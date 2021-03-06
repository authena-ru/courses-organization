package v1

import (
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/go-chi/render"

	"github.com/authena-ru/courses-organization/internal/app"
	"github.com/authena-ru/courses-organization/internal/domain/course"
)

func marshalCommonCourses(w http.ResponseWriter, r *http.Request, courses []app.CommonCourse) {
	response := make(GetAllCoursesResponse, 0, len(courses))
	for _, c := range courses {
		response = append(response, marshalCommonCourseToCourseResponse(c))
	}

	render.Respond(w, r, response)
}

func marshalCommonCourse(w http.ResponseWriter, r *http.Request, crs app.CommonCourse) {
	response := marshalCommonCourseToCourseResponse(crs)

	render.Respond(w, r, response)
}

func marshalCommonCourseToCourseResponse(crs app.CommonCourse) Course {
	return Course{
		Id:          crs.ID,
		Title:       crs.Title,
		Period:      marshalPeriod(crs.Period),
		CreatorId:   crs.CreatorID,
		Started:     crs.Started,
		TasksNumber: crs.TasksNumber,
	}
}

func marshalSpecificTask(w http.ResponseWriter, r *http.Request, task app.SpecificTask) {
	type taskResponse struct {
		TaskResponse
		Deadline *Deadline   `json:"deadline,omitempty"`
		TestData []TestData  `json:"testData,omitempty"`
		Points   []TestPoint `json:"points,omitempty"`
	}

	response := taskResponse{
		TaskResponse: TaskResponse{
			Number: task.Number,
			Task: Task{
				Title:       task.Title,
				Description: task.Description,
				Type:        marshalTaskType(task.Type),
			},
		},
		Deadline: marshalDeadline(task.Deadline),
		TestData: marshalTestData(task.TestData),
		Points:   marshalTestPoints(task.Points),
	}

	render.Respond(w, r, response)
}

func marshalGeneralTasks(w http.ResponseWriter, r *http.Request, tasks []app.GeneralTask) {
	response := make([]TaskResponse, 0, len(tasks))
	for _, t := range tasks {
		response = append(response, TaskResponse{
			Number: t.Number,
			Task: Task{
				Title:       t.Title,
				Description: t.Description,
				Type:        marshalTaskType(t.Type),
			},
		})
	}

	render.Respond(w, r, response)
}

func marshalTaskType(taskType course.TaskType) TaskType {
	switch taskType {
	case course.ManualCheckingType:
		return TaskTypeMANUALCHECKING
	case course.AutoCodeCheckingType:
		return TaskTypeAUTOCODECHECKING
	case course.TestingType:
		return TaskTypeTESTING
	}

	return "UNKNOWN"
}

func marshalSemester(semester course.Semester) Semester {
	switch semester {
	case course.FirstSemester:
		return SemesterFIRST
	case course.SecondSemester:
		return SemesterSECOND
	}

	return "UNKNOWN"
}

func marshalPeriod(period app.Period) CoursePeriod {
	return CoursePeriod{
		AcademicStartYear: period.AcademicStartYear,
		AcademicEndYear:   period.AcademicEndYear,
		Semester:          marshalSemester(period.Semester),
	}
}

func marshalDeadline(deadline *app.Deadline) *Deadline {
	if deadline == nil {
		return nil
	}

	return &Deadline{
		ExcellentGradeTime: types.Date{Time: deadline.ExcellentGradeTime},
		GoodGradeTime:      types.Date{Time: deadline.GoodGradeTime},
	}
}

func marshalTestData(testData []app.TestData) []TestData {
	marshalled := make([]TestData, 0, len(testData))

	for _, td := range testData {
		inputData, outputData := td.InputData, td.OutputData

		marshalled = append(marshalled, TestData{
			InputData:  &inputData,
			OutputData: &outputData,
		})
	}

	return marshalled
}

func marshalTestPoints(testPoints []app.TestPoint) []TestPoint {
	marshalled := make([]TestPoint, 0, len(testPoints))

	for _, tp := range testPoints {
		var correctVariantNumbers []int

		if tp.CorrectVariantNumbers != nil {
			correctVariantNumbers = make([]int, len(tp.CorrectVariantNumbers))
			copy(correctVariantNumbers, tp.CorrectVariantNumbers)
		}

		singleCorrectVariant := tp.SingleCorrectVariant

		marshalled = append(marshalled, TestPoint{
			Description:           tp.Description,
			Variants:              tp.Variants,
			CorrectVariantNumbers: &correctVariantNumbers,
			SingleCorrectVariant:  &singleCorrectVariant,
		})
	}

	return marshalled
}
