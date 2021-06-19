package mongodb

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
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
) (query.SpecificTask, error) {
	filter := makeFindTaskFilter(academic, courseID, taskNumber)
	projection := bson.M{"tasks": bson.M{"$elemMatch": bson.M{"number": taskNumber}}}
	opt := options.FindOne().SetProjection(projection)

	var document courseDocument
	if err := r.courses.FindOne(ctx, filter, opt).Decode(&document); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return query.SpecificTask{}, app.Wrap(app.ErrCourseTaskDoesntExist, err)
		}
		return query.SpecificTask{}, app.Wrap(app.ErrDatabaseProblems, err)
	}

	return unmarshallSpecificTask(academic, document.Tasks[0]), nil
}

func makeFindTaskFilter(academic course.Academic, courseID string, taskNumber int) bson.M {
	filter := makeCourseForAcademicFilter(academic, courseID)
	filter["tasks.number"] = taskNumber
	return filter
}

func (r *CoursesRepository) FindAllTasks(
	ctx context.Context,
	academic course.Academic, courseID string,
	filterParams query.TasksFilterParams,
) ([]query.GeneralTask, error) {
	filter := makeCourseForAcademicFilter(academic, courseID)
	projection := makeFindAllTasksProjection(filterParams)
	opt := options.FindOne().SetProjection(projection)

	var document courseDocument
	if err := r.courses.FindOne(ctx, filter, opt).Decode(&document); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, app.Wrap(app.ErrCourseDoesntExist, err)
		}
		return nil, app.Wrap(app.ErrDatabaseProblems, err)
	}

	return unmarshallGeneralTasks(document.Tasks), nil
}

func makeFindAllTasksProjection(filterParams query.TasksFilterParams) bson.M {
	return bson.M{"tasks": true}
}

func makeCourseForAcademicFilter(academic course.Academic, courseID string) bson.M {
	filter := bson.M{"_id": courseID}
	if academic.Type() == course.StudentType {
		filter["students"] = bson.M{"$elemMatch": bson.M{"$eq": academic.ID()}}
	} else {
		filter["$or"] = bson.A{
			bson.M{"creatorId": academic.ID()},
			bson.M{"collaborators": bson.M{"$elemMatch": bson.M{"$eq": academic.ID()}}},
		}
	}
	return filter
}
