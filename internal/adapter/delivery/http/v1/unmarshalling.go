package v1

import (
	"github.com/authena-ru/courses-organization/internal/app/query"
	"net/http"

	"github.com/go-chi/render"

	"github.com/authena-ru/courses-organization/internal/adapter/delivery/http/auth"
	"github.com/authena-ru/courses-organization/internal/adapter/delivery/http/httperr"
	"github.com/authena-ru/courses-organization/internal/app/command"
	"github.com/authena-ru/courses-organization/internal/domain/course"
)

func decode(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if err := render.Decode(r, v); err != nil {
		httperr.BadRequest("invalid-request-body", err, w, r)
		return false
	}
	return true
}

func unmarshallAllTasksQuery(
	w http.ResponseWriter, r *http.Request,
	courseID string, params GetCourseTasksParams,
) (qry query.AllTasksQuery, ok bool) {
	academic, ok := unmarshallAcademic(w, r)
	if !ok {
		return
	}

	var title, description string
	if params.Title != nil {
		title = *params.Title
	}
	if params.Description != nil {
		description = *params.Description
	}
	var taskType course.TaskType
	switch *params.Type {
	case TaskTypeMANUALCHECKING:
		taskType = course.ManualCheckingType
	case TaskTypeAUTOCODECHECKING:
		taskType = course.AutoCodeCheckingType
	case TaskTypeTESTING:
		taskType = course.TestingType
	default:
		taskType = course.TaskType(0)
	}

	return query.AllTasksQuery{
		Academic:    academic,
		CourseID:    courseID,
		Type:        taskType,
		Title:       title,
		Description: description,
	}, true
}

func unmarshallSpecificTaskQuery(
	w http.ResponseWriter, r *http.Request,
	courseID string, taskNumber int,
) (qry query.SpecificTaskQuery, ok bool) {
	academic, ok := unmarshallAcademic(w, r)
	if !ok {
		return
	}
	return query.SpecificTaskQuery{
		Academic:   academic,
		CourseID:   courseID,
		TaskNumber: taskNumber,
	}, true
}

func unmarshallAddStudentCommand(
	w http.ResponseWriter, r *http.Request,
	courseID string,
) (cmd command.AddStudentCommand, ok bool) {
	academic, ok := unmarshallAcademic(w, r)
	if !ok {
		return
	}
	var rb AddStudentToCourseRequest
	if ok = decode(w, r, &rb); !ok {
		return
	}
	return command.AddStudentCommand{
		Academic:  academic,
		CourseID:  courseID,
		StudentID: rb.Id,
	}, false
}

func unmarshallRemoveStudentCommand(
	w http.ResponseWriter, r *http.Request,
	courseID, studentID string,
) (cmd command.RemoveStudentCommand, ok bool) {
	academic, ok := unmarshallAcademic(w, r)
	if !ok {
		return
	}
	return command.RemoveStudentCommand{
		Academic:  academic,
		CourseID:  courseID,
		StudentID: studentID,
	}, true
}

func unmarshallAddCollaboratorCommand(
	w http.ResponseWriter, r *http.Request,
	courseID string,
) (cmd command.AddCollaboratorCommand, ok bool) {
	academic, ok := unmarshallAcademic(w, r)
	if !ok {
		return
	}
	var rb AddCollaboratorToCourseRequest
	if ok = decode(w, r, &rb); !ok {
		return
	}
	return command.AddCollaboratorCommand{
		Academic:       academic,
		CourseID:       courseID,
		CollaboratorID: rb.Id,
	}, true
}

func unmarshallRemoveCollaboratorCommand(
	w http.ResponseWriter, r *http.Request,
	courseID, collaboratorID string,
) (cmd command.RemoveCollaboratorCommand, ok bool) {
	academic, ok := unmarshallAcademic(w, r)
	if !ok {
		return
	}
	return command.RemoveCollaboratorCommand{
		Academic:       academic,
		CourseID:       courseID,
		CollaboratorID: collaboratorID,
	}, true
}

func unmarshallCreateCourseCommand(w http.ResponseWriter, r *http.Request) (cmd command.CreateCourseCommand, ok bool) {
	academic, ok := unmarshallAcademic(w, r)
	if !ok {
		return
	}
	var rb CreateCourseRequest
	if ok = decode(w, r, &rb); !ok {
		return
	}
	period, ok := unmarshallPeriod(w, r, &rb.Period)
	if !ok {
		return
	}
	return command.CreateCourseCommand{
		Academic:      academic,
		CourseStarted: rb.Started,
		CourseTitle:   rb.Title,
		CoursePeriod:  period,
	}, true
}

func unmarshallExtendCourseCommand(
	w http.ResponseWriter, r *http.Request,
	courseID string,
) (cmd command.ExtendCourseCommand, ok bool) {
	academic, ok := unmarshallAcademic(w, r)
	if !ok {
		return
	}
	var rb ExtendCourseRequest
	if ok = decode(w, r, &rb); !ok {
		return
	}
	period, ok := unmarshallPeriod(w, r, rb.Period)
	if !ok {
		return
	}
	return command.ExtendCourseCommand{
		Academic:       academic,
		OriginCourseID: courseID,
		CourseStarted:  rb.Started,
		CourseTitle:    rb.Title,
		CoursePeriod:   period,
	}, true
}

func unmarshallAddTaskCommand(
	w http.ResponseWriter, r *http.Request,
	courseID string,
) (cmd command.AddTaskCommand, ok bool) {
	academic, ok := unmarshallAcademic(w, r)
	if !ok {
		return
	}
	rb := struct {
		Task
		Deadline Deadline
		TestData []TestData
		Points   []TestPoint
	}{}
	if ok = decode(w, r, &rb); !ok {
		return
	}
	taskType, ok := unmarshallTaskType(w, r, rb.Type)
	if !ok {
		return
	}
	deadline, ok := unmarshallDeadline(w, r, &rb.Deadline)
	if !ok {
		return
	}
	testData, ok := unmarshallTestData(w, r, &rb.TestData)
	if !ok {
		return
	}
	testPoints, ok := unmarshallTestPoints(w, r, &rb.Points)
	if !ok {
		return
	}
	return command.AddTaskCommand{
		Academic:        academic,
		CourseID:        courseID,
		TaskTitle:       rb.Title,
		TaskDescription: rb.Description,
		TaskType:        taskType,
		Deadline:        deadline,
		TestPoints:      testPoints,
		TestData:        testData,
	}, true
}

func unmarshallTaskType(w http.ResponseWriter, r *http.Request, apiTaskType TaskType) (course.TaskType, bool) {
	switch apiTaskType {
	case TaskTypeMANUALCHECKING:
		return course.ManualCheckingType, true
	case TaskTypeAUTOCODECHECKING:
		return course.AutoCodeCheckingType, true
	case TaskTypeTESTING:
		return course.TestingType, true
	}
	httperr.BadRequest("invalid-task-type", nil, w, r)
	return course.TaskType(0), false
}

func unmarshallPeriod(w http.ResponseWriter, r *http.Request, apiPeriod *CoursePeriod) (course.Period, bool) {
	if apiPeriod == nil {
		return course.Period{}, true
	}
	var semester course.Semester
	switch apiPeriod.Semester {
	case SemesterFIRST:
		semester = course.FirstSemester
	case SemesterSECOND:
		semester = course.SecondSemester
	}
	domainPeriod, err := course.NewPeriod(apiPeriod.AcademicStartYear, apiPeriod.AcademicEndYear, semester)
	if err != nil {
		httperr.BadRequest("invalid-course-period", err, w, r)
		return course.Period{}, false
	}
	return domainPeriod, true
}

func unmarshallAcademic(w http.ResponseWriter, r *http.Request) (course.Academic, bool) {
	academic, err := auth.AcademicFromCtx(r.Context())
	if err != nil {
		httperr.Unauthorized("no-user-in-context", err, w, r)
		return course.Academic{}, false
	}
	return academic, true
}

func unmarshallDeadline(w http.ResponseWriter, r *http.Request, apiDeadline *Deadline) (course.Deadline, bool) {
	if apiDeadline == nil {
		return course.Deadline{}, true
	}
	deadline, err := course.NewDeadline(apiDeadline.ExcellentGradeTime.Time, apiDeadline.GoodGradeTime.Time)
	if err != nil {
		httperr.BadRequest("invalid-task", err, w, r)
		return course.Deadline{}, false
	}
	return deadline, true
}

func unmarshallTestData(w http.ResponseWriter, r *http.Request, apiTestData *[]TestData) ([]course.TestData, bool) {
	if apiTestData == nil {
		return nil, true
	}
	dereferencedAPITestData := *apiTestData
	testData := make([]course.TestData, 0, len(dereferencedAPITestData))
	for _, atd := range dereferencedAPITestData {
		var inputData, outputData string
		if atd.InputData != nil {
			inputData = *atd.InputData
		}
		if atd.OutputData != nil {
			outputData = *atd.OutputData
		}
		td, err := course.NewTestData(inputData, outputData)
		if err != nil {
			httperr.BadRequest("invalid-task", err, w, r)
			return nil, false
		}
		testData = append(testData, td)
	}
	return testData, true
}

func unmarshallTestPoints(w http.ResponseWriter, r *http.Request, apiTestPoints *[]TestPoint) ([]course.TestPoint, bool) {
	if apiTestPoints == nil {
		return nil, true
	}
	dereferencedAPITestPoints := *apiTestPoints
	testPoints := make([]course.TestPoint, 0, len(dereferencedAPITestPoints))
	for _, atp := range dereferencedAPITestPoints {
		var correctVariantNumbers []int
		if atp.CorrectVariantNumbers != nil {
			correctVariantNumbers = *atp.CorrectVariantNumbers
		}
		tp, err := course.NewTestPoint(atp.Description, atp.Variants, correctVariantNumbers)
		if err != nil {
			httperr.BadRequest("invalid-task", err, w, r)
			return nil, false
		}
		testPoints = append(testPoints, tp)
	}
	return testPoints, true
}
