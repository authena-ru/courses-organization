package mongodb

import (
	"github.com/authena-ru/courses-organization/internal/domain/course"
	"time"
)

type courseDocument struct {
	ID            string         `bson:"_id,omitempty"`
	Title         string         `bson:"title"`
	Period        periodDocument `bson:"period"`
	Started       bool           `bson:"started"`
	CreatorID     string         `bson:"creatorId"`
	Collaborators []string       `bson:"collaborators,omitempty"`
	Students      []string       `bson:"students,omitempty"`
	Tasks         []taskDocument `bson:"tasks,omitempty"`
}

type periodDocument struct {
	AcademicStartYear int             `bson:"academicStartYear"`
	AcademicEndYear   int             `bson:"academicEndYear"`
	Semester          course.Semester `bson:"semester"`
}

type taskDocument struct {
	Number      int                 `bson:"number"`
	Title       string              `bson:"title"`
	Description string              `bson:"description"`
	TaskType    course.TaskType     `bson:"taskType"`
	Deadline    deadlineDocument    `bson:"deadline"`
	TestPoints  []testPointDocument `bson:"testPoints"`
	TestData    []testDataDocument  `bson:"testData"`
}

type deadlineDocument struct {
	GoodGradeTime      time.Time `bson:"goodGradeTime"`
	ExcellentGradeTime time.Time `bson:"ExcellentGradeTime"`
}

type testPointDocument struct {
	Description           string   `bson:"description"`
	Variants              []string `bson:"variants"`
	CorrectVariantNumbers []int    `bson:"correctVariantNumbers"`
}

type testDataDocument struct {
	InputData  string `bson:"inputData"`
	OutputData string `bson:"outputData"`
}

func newCourseDocument(crs *course.Course) courseDocument {
	return courseDocument{
		ID:    crs.ID(),
		Title: crs.Title(),
		Period: periodDocument{
			AcademicStartYear: crs.Period().AcademicStartYear(),
			AcademicEndYear:   crs.Period().AcademicEndYear(),
			Semester:          crs.Period().Semester(),
		},
		Started:       crs.Started(),
		CreatorID:     crs.CreatorID(),
		Collaborators: crs.Collaborators(),
		Students:      crs.Students(),
		// tasks
	}
}

func unmarshallCourse(document courseDocument) (*course.Course, error) {
	// TODO: unmarshall other parameters
	return course.UnmarshallFromDatabase(course.UnmarshallingParams{
		ID:    document.ID,
		Title: document.Title,
	}), nil
}
