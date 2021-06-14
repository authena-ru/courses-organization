package mongodb

import "github.com/authena-ru/courses-organization/internal/domain/course"

type courseModel struct {
	id            string      `bson:"_id,omitempty"`
	title         string      `bson:"title"`
	period        periodModel `bson:"period"`
	started       bool        `bson:"started"`
	creatorID     string      `bson:"creatorId"`
	collaborators []string    `bson:"collaborators"`
	students      []string    `bson:"students"`
	// tasks
	nextTaskNumber int `bson:"nextTaskNumber"`
}

type periodModel struct {
	academicStartYear int             `bson:"academicStartYear"`
	academicEndYear   int             `bson:"academicEndYear"`
	semester          course.Semester `bson:"semester"`
}

func newCourseModel(crs *course.Course) courseModel {
	return courseModel{
		id:    crs.ID(),
		title: crs.Title(),
		period: periodModel{
			academicStartYear: crs.Period().AcademicStartYear(),
			academicEndYear:   crs.Period().AcademicEndYear(),
			semester:          crs.Period().Semester(),
		},
		started:       crs.Started(),
		creatorID:     crs.CreatorID(),
		collaborators: crs.Collaborators(),
		students:      crs.Students(),
		// tasks
		// nextTaskNumber
	}
}

func newCourse(courseModel courseModel) (*course.Course, error) {
	// TODO: unmarshall other parameters
	return course.UnmarshallFromDatabase(course.UnmarshallingParams{
		ID:    courseModel.id,
		Title: courseModel.title,
	}), nil
}
