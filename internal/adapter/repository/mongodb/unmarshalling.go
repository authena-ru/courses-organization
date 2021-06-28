package mongodb

import (
	"github.com/authena-ru/courses-organization/internal/app"
	"github.com/authena-ru/courses-organization/internal/domain/course"
)

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
			TaskType:    td.Type,
			Deadline:    unmarshalDeadline(td.Deadline),
			TestData:    unmarshallTestData(td.TestData),
			TestPoints:  unmarshallTestPoints(td.TestPoints),
		})
	}
	return taskParams
}

func unmarshalDeadline(document *deadlineDocument) course.Deadline {
	if document == nil {
		return course.Deadline{}
	}
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

func unmarshallSpecificTask(academic course.Academic, document taskDocument) app.SpecificTask {
	forTeacher := academic.Type() == course.TeacherType
	return app.SpecificTask{
		Number:      document.Number,
		Title:       document.Title,
		Description: document.Description,
		Type:        document.Type,
		Deadline:    unmarshallQueryDeadline(document.Deadline),
		TestData:    unmarshallQueryTestData(forTeacher, document.TestData),
		Points:      unmarshallQueryTestPoints(forTeacher, document.TestPoints),
	}
}

func unmarshallQueryDeadline(document *deadlineDocument) *app.Deadline {
	if document == nil {
		return nil
	}
	return &app.Deadline{
		ExcellentGradeTime: document.ExcellentGradeTime,
		GoodGradeTime:      document.GoodGradeTime,
	}
}

func unmarshallQueryTestData(forTeacher bool, documents []testDataDocument) []app.TestData {
	if !forTeacher {
		return nil
	}
	queryTestData := make([]app.TestData, 0, len(documents))
	for _, d := range documents {
		queryTestData = append(queryTestData, app.TestData{
			InputData:  d.InputData,
			OutputData: d.OutputData,
		})
	}
	return queryTestData
}

func unmarshallQueryTestPoints(forTeacher bool, documents []testPointDocument) []app.TestPoint {
	queryTestPoints := make([]app.TestPoint, 0, len(documents))
	for _, d := range documents {
		var correctVariantNumbers []int
		if forTeacher {
			correctVariantNumbers = d.CorrectVariantNumbers
		}
		queryTestPoints = append(queryTestPoints, app.TestPoint{
			Description:           d.Description,
			Variants:              d.Variants,
			CorrectVariantNumbers: correctVariantNumbers,
			SingleCorrectVariant:  len(d.CorrectVariantNumbers) > 1,
		})
	}
	return queryTestPoints
}

func unmarshallGeneralTasks(documents []taskDocument) []app.GeneralTask {
	tasks := make([]app.GeneralTask, 0, len(documents))
	for _, d := range documents {
		tasks = append(tasks, app.GeneralTask{
			Number:      d.Number,
			Title:       d.Title,
			Description: d.Description,
			Type:        d.Type,
		})
	}
	return tasks
}
