package mongodb

import (
	"time"

	"github.com/authena-ru/courses-organization/internal/domain/course"
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
	Deadline    *deadlineDocument   `bson:"deadline,omitempty"`
	TestPoints  []testPointDocument `bson:"testPoints,omitempty"`
	TestData    []testDataDocument  `bson:"testData,omitempty"`
}

type deadlineDocument struct {
	GoodGradeTime      time.Time `bson:"goodGradeTime"`
	ExcellentGradeTime time.Time `bson:"excellentGradeTime"`
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

func marshallCourseDocument(crs *course.Course) courseDocument {
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
		Tasks:         marshallTaskDocuments(crs.Tasks()),
	}
}

func marshallTaskDocuments(tasks []course.Task) []taskDocument {
	taskDocuments := make([]taskDocument, 0, len(tasks))
	for _, t := range tasks {
		deadline, _ := t.Deadline()
		var deadlineDoc *deadlineDocument
		if !deadline.IsZero() {
			deadlineDoc = &deadlineDocument{
				GoodGradeTime:      deadline.GoodGradeTime(),
				ExcellentGradeTime: deadline.ExcellentGradeTime(),
			}
		}
		testData, _ := t.TestData()
		testPoints, _ := t.TestPoints()
		taskDocuments = append(taskDocuments, taskDocument{
			Number:      t.Number(),
			Title:       t.Title(),
			Description: t.Description(),
			TaskType:    t.Type(),
			Deadline:    deadlineDoc,
			TestData:    marshallTestDataDocuments(testData),
			TestPoints:  marshallTestPointDocuments(testPoints),
		})
	}
	return taskDocuments
}

func marshallTestDataDocuments(testData []course.TestData) []testDataDocument {
	testDataDocuments := make([]testDataDocument, 0, len(testData))
	for _, td := range testData {
		testDataDocuments = append(testDataDocuments, testDataDocument{
			InputData:  td.InputData(),
			OutputData: td.OutputData(),
		})
	}
	return testDataDocuments
}

func marshallTestPointDocuments(testPoints []course.TestPoint) []testPointDocument {
	testPointDocuments := make([]testPointDocument, 0, len(testPoints))
	for _, tp := range testPoints {
		testPointDocuments = append(testPointDocuments, testPointDocument{
			Description:           tp.Description(),
			Variants:              tp.Variants(),
			CorrectVariantNumbers: tp.CorrectVariantNumbers(),
		})
	}
	return testPointDocuments
}
