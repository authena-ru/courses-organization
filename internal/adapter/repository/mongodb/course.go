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

func unmarshallCourse(document courseDocument) (*course.Course, error) {
	period, err := unmarshallPeriod(document.Period)
	if err != nil {
		return nil, err
	}
	tasks, err := unmarshallTasks(document.Tasks)
	if err != nil {
		return nil, err
	}
	return course.UnmarshallFromDatabase(course.UnmarshallingParams{
		ID:            document.ID,
		Title:         document.Title,
		Period:        period,
		Started:       document.Started,
		CreatorID:     document.CreatorID,
		Collaborators: document.Collaborators,
		Students:      document.Students,
		Tasks:         tasks,
	}), nil
}

func unmarshallPeriod(document periodDocument) (course.Period, error) {
	return course.NewPeriod(document.AcademicStartYear, document.AcademicEndYear, document.Semester)
}

func unmarshallTasks(taskDocuments []taskDocument) ([]course.UnmarshallingTaskParams, error) {
	taskParams := make([]course.UnmarshallingTaskParams, 0, len(taskDocuments))
	for _, td := range taskDocuments {
		deadline, err := unmarshalDeadline(td.Deadline)
		if err != nil {
			return nil, err
		}
		testData, err := unmarshallTestData(td.TestData)
		if err != nil {
			return nil, err
		}
		testPoints, err := unmarshallTestPoints(td.TestPoints)
		if err != nil {
			return nil, err
		}
		taskParams = append(taskParams, course.UnmarshallingTaskParams{
			Number:      td.Number,
			Title:       td.Title,
			Description: td.Description,
			TaskType:    td.TaskType,
			Deadline:    deadline,
			TestData:    testData,
			TestPoints:  testPoints,
		})
	}
	return taskParams, nil
}

func unmarshalDeadline(document deadlineDocument) (course.Deadline, error) {
	deadline, err := course.NewDeadline(document.ExcellentGradeTime, document.GoodGradeTime)
	if err != nil {
		return course.Deadline{}, err
	}
	return deadline, nil
}

func unmarshallTestData(documents []testDataDocument) ([]course.TestData, error) {
	testData := make([]course.TestData, 0, len(documents))
	for _, d := range documents {
		td, err := course.NewTestData(d.InputData, d.OutputData)
		if err != nil {
			return nil, err
		}
		testData = append(testData, td)
	}
	return testData, nil
}

func unmarshallTestPoints(documents []testPointDocument) ([]course.TestPoint, error) {
	testPoints := make([]course.TestPoint, 0, len(documents))
	for _, d := range documents {
		tp, err := course.NewTestPoint(d.Description, d.Variants, d.CorrectVariantNumbers)
		if err != nil {
			return nil, err
		}
		testPoints = append(testPoints, tp)
	}
	return testPoints, nil
}
