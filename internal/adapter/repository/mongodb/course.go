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
	Deadline    deadlineDocument    `bson:"deadline,omitempty"`
	TestPoints  []testPointDocument `bson:"testPoints,omitempty"`
	TestData    []testDataDocument  `bson:"testData,omitempty"`
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
		Tasks:         newTaskDocuments(crs.Tasks()),
	}
}

func newTaskDocuments(tasks []course.Task) []taskDocument {
	taskDocuments := make([]taskDocument, 0, len(tasks))
	for _, t := range tasks {
		deadline, _ := t.Deadline()
		testData, _ := t.TestData()
		testPoints, _ := t.TestPoints()
		taskDocuments = append(taskDocuments, taskDocument{
			Number:      t.Number(),
			Title:       t.Title(),
			Description: t.Description(),
			TaskType:    t.Type(),
			Deadline: deadlineDocument{
				GoodGradeTime:      deadline.GoodGradeTime(),
				ExcellentGradeTime: deadline.ExcellentGradeTime(),
			},
			TestData:   newTestDataDocuments(testData),
			TestPoints: newTestPointDocuments(testPoints),
		})
	}
	return taskDocuments
}

func newTestDataDocuments(testData []course.TestData) []testDataDocument {
	testDataDocuments := make([]testDataDocument, 0, len(testData))
	for _, td := range testData {
		testDataDocuments = append(testDataDocuments, testDataDocument{
			InputData:  td.InputData(),
			OutputData: td.OutputData(),
		})
	}
	return testDataDocuments
}

func newTestPointDocuments(testPoints []course.TestPoint) []testPointDocument {
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

func unmarshallCourse(document courseDocument) *course.Course {
	return course.UnmarshallFromDatabase(course.UnmarshallingParams{
		ID:            document.ID,
		Title:         document.Title,
		Period:        unmarshallPeriod(document.Period),
		Started:       document.Started,
		CreatorID:     document.CreatorID,
		Collaborators: document.Collaborators,
		Students:      document.Students,
		Tasks:         unmarshallTasks(document.Tasks),
	})
}

func unmarshallPeriod(document periodDocument) course.Period {
	return course.MustNewPeriod(document.AcademicStartYear, document.AcademicEndYear, document.Semester)
}

func unmarshallTasks(taskDocuments []taskDocument) []course.UnmarshallingTaskParams {
	taskParams := make([]course.UnmarshallingTaskParams, 0, len(taskDocuments))
	for _, td := range taskDocuments {
		taskParams = append(taskParams, course.UnmarshallingTaskParams{
			Number:      td.Number,
			Title:       td.Title,
			Description: td.Description,
			TaskType:    td.TaskType,
			Deadline:    unmarshalDeadline(td.Deadline),
			TestData:    unmarshallTestData(td.TestData),
			TestPoints:  unmarshallTestPoints(td.TestPoints),
		})
	}
	return taskParams
}

func unmarshalDeadline(document deadlineDocument) course.Deadline {
	return course.MustNewDeadline(document.ExcellentGradeTime, document.GoodGradeTime)
}

func unmarshallTestData(documents []testDataDocument) []course.TestData {
	testData := make([]course.TestData, 0, len(documents))
	for _, d := range documents {
		testData = append(testData, course.MustNewTestData(d.OutputData, d.OutputData))
	}
	return testData
}

func unmarshallTestPoints(documents []testPointDocument) []course.TestPoint {
	testPoints := make([]course.TestPoint, 0, len(documents))
	for _, d := range documents {
		testPoints = append(testPoints, course.MustNewTestPoint(d.Description, d.Variants, d.CorrectVariantNumbers))
	}
	return testPoints
}
