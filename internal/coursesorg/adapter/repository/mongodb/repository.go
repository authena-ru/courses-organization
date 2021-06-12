package mongodb

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/authena-ru/courses-organization/internal/coursesorg/app"
	"github.com/authena-ru/courses-organization/internal/coursesorg/app/command"
	"github.com/authena-ru/courses-organization/internal/coursesorg/domain/course"
)

type CoursesRepository struct {
	courses *mongo.Collection
}

var coursesCollection = "courses"

func NewCoursesRepository(db *mongo.Database) *CoursesRepository {
	return &CoursesRepository{courses: db.Collection(coursesCollection)}
}

func (r *CoursesRepository) AddCourse(ctx context.Context, crs *course.Course) error {
	_, err := r.courses.InsertOne(ctx, newCourseModel(crs))
	return err
}

func (r *CoursesRepository) GetCourse(ctx context.Context, courseID string) (*course.Course, error) {
	var courseModel courseModel
	if err := r.courses.FindOne(ctx, bson.M{"_id": courseID}).Decode(&courseModel); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, app.Wrap(app.ErrCourseDoesntExist, err)
		}
		return nil, err
	}

	return newCourse(courseModel)
}

func (r *CoursesRepository) UpdateCourse(ctx context.Context, courseID string, updateFn command.UpdateFunction) error {
	session, err := r.courses.Database().Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		var courseModel courseModel
		if err := r.courses.FindOne(ctx, bson.M{"_id": courseID}).Decode(&courseModel); err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, app.Wrap(app.ErrCourseDoesntExist, err)
			}
			return nil, err
		}

		crs, err := newCourse(courseModel)
		if err != nil {
			return nil, err
		}
		updatedCourse, err := updateFn(ctx, crs)
		if err != nil {
			return nil, err
		}

		replaceOpts := options.Replace().SetUpsert(true)
		filter := bson.M{"_id": courseID}
		if _, err := r.courses.ReplaceOne(ctx, filter, newCourseModel(updatedCourse), replaceOpts); err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}
