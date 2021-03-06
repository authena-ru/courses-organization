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

func (r *CoursesRepository) FindCourse(
	ctx context.Context,
	academic course.Academic,
	courseID string,
) (app.CommonCourse, error) {
	filter := makeCourseForAcademicFilter(academic, courseID)

	var document courseDocument
	if err := r.courses.FindOne(ctx, filter).Decode(&document); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return app.CommonCourse{}, app.Wrap(app.ErrCourseDoesntExist, err)
		}
	}

	return unmarshalCommonCourse(document), nil
}

func (r *CoursesRepository) FindAllCourses(
	ctx context.Context,
	academic course.Academic,
	params query.CoursesFilterParams,
) ([]app.CommonCourse, error) {
	filter := makeFindAllCoursesFilter(academic, params)

	cursor, err := r.courses.Find(ctx, filter)
	if err != nil {
		return nil, app.Wrap(app.ErrDatabaseProblems, err)
	}

	var documents []courseDocument
	if err := cursor.All(ctx, &documents); err != nil {
		return nil, app.Wrap(app.ErrDatabaseProblems, err)
	}

	return unmarshalCommonCourses(documents), nil
}

func makeFindAllCoursesFilter(academic course.Academic, filterParams query.CoursesFilterParams) bson.D {
	return bson.D{makeCoursesForAcademicFilter(academic), makeFindCoursesTitleFilter(filterParams)}
}

func makeFindCoursesTitleFilter(filterParams query.CoursesFilterParams) bson.E {
	if filterParams.Title == "" {
		return bson.E{}
	}

	regex := primitive.Regex{Pattern: filterParams.Title, Options: "i"}

	return bson.E{
		Key:   "title",
		Value: regex,
	}
}

func (r *CoursesRepository) AddCourse(ctx context.Context, crs *course.Course) error {
	_, err := r.courses.InsertOne(ctx, marshalCourseDocument(crs))

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

	return unmarshalCourse(document), nil
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

		crs := unmarshalCourse(document)
		updatedCourse, err := updateFn(ctx, crs)
		if err != nil {
			return nil, err
		}
		updatedCourseDocument := marshalCourseDocument(updatedCourse)

		replaceOpt := options.Replace().SetUpsert(true)
		filter := bson.M{"_id": updatedCourseDocument.ID}
		if _, err := r.courses.ReplaceOne(ctx, filter, updatedCourseDocument, replaceOpt); err != nil {
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
	findOpt := options.FindOne().SetProjection(projection)

	var document courseDocument
	if err := r.courses.FindOne(ctx, filter, findOpt).Decode(&document); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return app.SpecificTask{}, app.Wrap(app.ErrCourseDoesntExist, err)
		}

		return app.SpecificTask{}, app.Wrap(app.ErrDatabaseProblems, err)
	}

	if len(document.Tasks) == 0 {
		return app.SpecificTask{}, app.ErrTaskDoesntExist
	}

	return unmarshalSpecificTask(academic, document.Tasks[0]), nil
}

func makeFindTaskProjection(taskNumber int) bson.D {
	return bson.D{{
		Key: "tasks", Value: bson.D{{
			Key: "$elemMatch", Value: bson.D{{
				Key: "number", Value: bson.D{{
					Key: "$eq", Value: taskNumber,
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

	return unmarshalGeneralTasks(document.Tasks), nil
}

func makeFindAllTasksPipeline(
	academic course.Academic,
	courseID string,
	filterParams query.TasksFilterParams,
) mongo.Pipeline {
	matchState := bson.D{{Key: "$match", Value: makeCourseForAcademicFilter(academic, courseID)}}
	projectStage := bson.D{{
		Key: "$project", Value: bson.D{{
			Key: "tasks", Value: bson.D{{
				Key: "$filter", Value: bson.D{
					{Key: "input", Value: "$tasks"},
					{
						Key: "cond",
						Value: bson.D{{
							Key: "$and", Value: bson.A{
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
		Key: "$or",
		Value: bson.A{
			bson.D{{
				Key: "$regexMatch", Value: bson.D{
					{Key: "input", Value: "$$this.title"},
					{Key: "regex", Value: regex},
				},
			}},
			bson.D{{
				Key: "$regexMatch", Value: bson.D{
					{Key: "input", Value: "$$this.description"},
					{Key: "regex", Value: regex},
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
		Key: "$eq", Value: bson.A{"$$this.type", filterParams.Type},
	}}
}

func makeCourseForAcademicFilter(academic course.Academic, courseID string) bson.D {
	return bson.D{{Key: "_id", Value: courseID}, makeCoursesForAcademicFilter(academic)}
}

func makeCoursesForAcademicFilter(academic course.Academic) bson.E {
	if academic.Type() == course.StudentType {
		return bson.E{
			Key: "students", Value: bson.D{{
				Key: "$elemMatch", Value: bson.D{{
					Key: "$eq", Value: academic.ID(),
				}},
			}},
		}
	}

	return bson.E{
		Key: "$or", Value: bson.A{
			bson.D{{
				Key: "creatorId", Value: academic.ID(),
			}},
			bson.D{{
				Key: "collaborators", Value: bson.D{{
					Key: "$elemMatch", Value: bson.D{{
						Key: "$eq", Value: academic.ID(),
					}},
				}},
			}},
		},
	}
}

func (r *CoursesRepository) RemoveAllCourses(ctx context.Context) error {
	_, err := r.courses.DeleteMany(ctx, bson.D{})

	return errors.Wrap(err, "unable to remove all courses")
}
