package mongodb

import (
	"github.com/authena-ru/courses-organization/internal/app"
	"github.com/authena-ru/courses-organization/internal/domain/course"
)

func unmarshalCourse(document courseDocument) *course.Course {
	return course.UnmarshalFromDatabase(course.UnmarshallingParams{
		ID:            document.ID,
		Title:         document.Title,
		Period:        unmarshalPeriod(document.Period),
		Started:       document.Started,
		CreatorID:     document.CreatorID,
		Collaborators: document.Collaborators,
		Students:      document.Students,
		Tasks:         unmarshalTasks(document.Tasks),
	})
}

func unmarshalPeriod(document periodDocument) course.Period {
	return course.MustNewPeriod(document.AcademicStartYear, document.AcademicEndYear, document.Semester)
}

func unmarshalTasks(taskDocuments []taskDocument) []course.UnmarshallingTaskParams {
	taskParams := make([]course.UnmarshallingTaskParams, 0, len(taskDocuments))
	for _, td := range taskDocuments {
		taskParams = append(taskParams, course.UnmarshallingTaskParams{
			Number:      td.Number,
			Title:       td.Title,
			Description: td.Description,
			TaskType:    td.Type,
			Deadline:    unmarshalDeadline(td.Deadline),
			TestData:    unmarshalTestData(td.TestData),
			TestPoints:  unmarshalTestPoints(td.TestPoints),
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

func unmarshalTestData(documents []testDataDocument) []course.TestData {
	testData := make([]course.TestData, 0, len(documents))
	for _, d := range documents {
		testData = append(testData, course.MustNewTestData(d.OutputData, d.OutputData))
	}

	return testData
}

func unmarshalTestPoints(documents []testPointDocument) []course.TestPoint {
	testPoints := make([]course.TestPoint, 0, len(documents))
	for _, d := range documents {
		testPoints = append(testPoints, course.MustNewTestPoint(d.Description, d.Variants, d.CorrectVariantNumbers))
	}

	return testPoints
}

func unmarshalCommonCourses(documents []courseDocument) []app.CommonCourse {
	courses := make([]app.CommonCourse, 0, len(documents))
	for _, d := range documents {
		courses = append(courses, unmarshalCommonCourse(d))
	}

	return courses
}

func unmarshalCommonCourse(document courseDocument) app.CommonCourse {
	return app.CommonCourse{
		ID:          document.ID,
		Title:       document.Title,
		Period:      unmarshalQueryPeriod(document.Period),
		CreatorID:   document.CreatorID,
		Started:     document.Started,
		TasksNumber: len(document.Tasks),
	}
}

func unmarshalSpecificTask(academic course.Academic, document taskDocument) app.SpecificTask {
	forTeacher := academic.Type() == course.TeacherType

	return app.SpecificTask{
		Number:      document.Number,
		Title:       document.Title,
		Description: document.Description,
		Type:        document.Type,
		Deadline:    unmarshalQueryDeadline(document.Deadline),
		TestData:    unmarshalQueryTestData(forTeacher, document.TestData),
		Points:      unmarshalQueryTestPoints(forTeacher, document.TestPoints),
	}
}

func unmarshalQueryPeriod(document periodDocument) app.Period {
	return app.Period{
		AcademicStartYear: document.AcademicStartYear,
		AcademicEndYear:   document.AcademicEndYear,
		Semester:          document.Semester,
	}
}

func unmarshalQueryDeadline(document *deadlineDocument) *app.Deadline {
	if document == nil {
		return nil
	}

	return &app.Deadline{
		ExcellentGradeTime: document.ExcellentGradeTime,
		GoodGradeTime:      document.GoodGradeTime,
	}
}

func unmarshalQueryTestData(forTeacher bool, documents []testDataDocument) []app.TestData {
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

func unmarshalQueryTestPoints(forTeacher bool, documents []testPointDocument) []app.TestPoint {
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

func unmarshalGeneralTasks(documents []taskDocument) []app.GeneralTask {
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
