package mongodb

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/authena-ru/courses-organization/internal/app"
	"github.com/authena-ru/courses-organization/internal/app/command"
	"github.com/authena-ru/courses-organization/internal/app/query"
	"github.com/authena-ru/courses-organization/internal/domain/course"
)

type CoursesRepository struct {
	courses *mongo.Collection
}

const coursesCollection = "courses"

func NewCoursesRepository(db *mongo.Database) *CoursesRepository {
	return &CoursesRepository{courses: db.Collection(coursesCollection)}
}

func (r *CoursesRepository) AddCourse(ctx context.Context, crs *course.Course) error {
	_, err := r.courses.InsertOne(ctx, marshallCourseDocument(crs))
	return app.Wrap(app.ErrDatabaseProblems, err)
}

func (r *CoursesRepository) GetCourse(ctx context.Context, courseID string) (*course.Course, error) {
	var document courseDocument
	if err := r.courses.FindOne(ctx, bson.M{"_id": courseID}).Decode(&document); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, app.Wrap(app.ErrCourseDoesntExist, err)
		}
		return nil, app.Wrap(app.ErrDatabaseProblems, err)
	}

	return unmarshallCourse(document), nil
}

func (r *CoursesRepository) UpdateCourse(ctx context.Context, courseID string, updateFn command.UpdateFunction) error {
	session, err := r.courses.Database().Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		var document courseDocument
		if err := r.courses.FindOne(ctx, bson.M{"_id": courseID}).Decode(&document); err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, app.Wrap(app.ErrCourseDoesntExist, err)
			}
			return nil, app.Wrap(app.ErrDatabaseProblems, err)
		}

		crs := unmarshallCourse(document)
		updatedCourse, err := updateFn(ctx, crs)
		if err != nil {
			return nil, err
		}
		updatedCourseDocument := marshallCourseDocument(updatedCourse)

		replaceOpts := options.Replace().SetUpsert(true)
		filter := bson.M{"_id": updatedCourseDocument.ID}
		if _, err := r.courses.ReplaceOne(ctx, filter, updatedCourseDocument, replaceOpts); err != nil {
			return nil, app.Wrap(app.ErrDatabaseProblems, err)
		}
		return nil, nil
	})
	return err
}

func (r *CoursesRepository) FindTask(
	ctx context.Context,
	academic course.Academic, courseID string, taskNumber int,
) (app.SpecificTask, error) {
	filter := makeCourseForAcademicFilter(academic, courseID)
	projection := makeFindTaskProjection(taskNumber)
	opt := options.FindOne().SetProjection(projection)

	var document courseDocument
	if err := r.courses.FindOne(ctx, filter, opt).Decode(&document); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return app.SpecificTask{}, app.Wrap(app.ErrCourseDoesntExist, err)
		}
		return app.SpecificTask{}, app.Wrap(app.ErrDatabaseProblems, err)
	}

	if len(document.Tasks) == 0 {
		return app.SpecificTask{}, app.ErrTaskDoesntExist
	}

	return unmarshallSpecificTask(academic, document.Tasks[0]), nil
}

func makeFindTaskProjection(taskNumber int) bson.D {
	return bson.D{{
		"tasks", bson.D{{
			"$elemMatch", bson.D{{
				"number", bson.D{{
					"$eq", taskNumber,
				}},
			}},
		}},
	}}
}

func (r *CoursesRepository) FindAllTasks(
	ctx context.Context,
	academic course.Academic, courseID string,
	filterParams query.TasksFilterParams,
) ([]app.GeneralTask, error) {
	pipeline := makeFindAllTasksPipeline(academic, courseID, filterParams)
	cursor, err := r.courses.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, app.Wrap(app.ErrDatabaseProblems, err)
	}
	if !cursor.Next(ctx) {
		return nil, app.ErrCourseDoesntExist
	}
	var document courseDocument
	if err := cursor.Decode(&document); err != nil {
		return nil, app.Wrap(app.ErrDatabaseProblems, err)
	}
	return unmarshallGeneralTasks(document.Tasks), nil
}

func makeFindAllTasksPipeline(
	academic course.Academic,
	courseID string,
	filterParams query.TasksFilterParams,
) mongo.Pipeline {
	matchState := bson.D{{"$match", makeCourseForAcademicFilter(academic, courseID)}}
	projectStage := bson.D{{
		"$project", bson.D{{
			"tasks", bson.D{{
				"$filter", bson.D{
					{"input", "$tasks"},
					{
						"cond",
						bson.D{{
							"$and", bson.A{
								makeFindAllTasksTextFilter(filterParams),
								makeFindAllTasksTypeFilter(filterParams),
							},
						}},
					},
				},
			}},
		}},
	}}
	return mongo.Pipeline{matchState, projectStage}
}

func makeFindAllTasksTextFilter(filterParams query.TasksFilterParams) bson.D {
	if filterParams.Text == "" {
		return bson.D{}
	}
	regex := primitive.Regex{Pattern: filterParams.Text, Options: "i"}
	return bson.D{{
		"$or",
		bson.A{
			bson.D{{
				"$regexMatch", bson.D{
					{"input", "$$this.title"},
					{"regex", regex},
				},
			}},
			bson.D{{
				"$regexMatch", bson.D{
					{"input", "$$this.description"},
					{"regex", regex},
				},
			}},
		},
	}}
}

func makeFindAllTasksTypeFilter(filterParams query.TasksFilterParams) bson.D {
	if !filterParams.Type.IsValid() {
		return bson.D{}
	}
	return bson.D{{
		"$eq", bson.A{"$$this.type", filterParams.Type},
	}}
}

func makeCourseForAcademicFilter(academic course.Academic, courseID string) bson.D {
	var academicSubFilter bson.E
	if academic.Type() == course.StudentType {
		academicSubFilter = bson.E{
			Key: "students", Value: bson.D{{
				"$elemMatch", bson.D{{"$eq", academic.ID()}}}},
		}
	} else {
		academicSubFilter = bson.E{
			Key: "$or", Value: bson.A{
				bson.D{{
					"creatorId", academic.ID(),
				}},
				bson.D{{
					"collaborators", bson.D{{
						"$elemMatch", bson.D{{"$eq", academic.ID()}}}},
				}},
			},
		}
	}
	return bson.D{{"_id", courseID}, academicSubFilter}
}
