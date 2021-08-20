package v1

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/authena-ru/courses-organization/internal/app"
	"github.com/authena-ru/courses-organization/internal/domain/course"
	"github.com/authena-ru/courses-organization/pkg/httperr"
)

func (h handler) AddTaskToCourse(w http.ResponseWriter, r *http.Request, courseID string) {
	cmd, ok := unmarshalAddTaskCommand(w, r, courseID)
	if !ok {
		return
	}

	taskNumber, err := h.app.Commands.AddTask.Handle(r.Context(), cmd)
	if err == nil {
		w.Header().Set("Content-Location", fmt.Sprintf("/courses/%s/tasks/%d", courseID, taskNumber))
		w.WriteHeader(http.StatusCreated)

		return
	}

	if errors.Is(err, app.ErrCourseDoesntExist) {
		httperr.NotFound("course-not-found", err, w, r)

		return
	}

	if course.IsInvalidTaskParametersError(err) {
		httperr.UnprocessableEntity("invalid-task-parameters", err, w, r)

		return
	}

	if course.IsInvalidDeadlineError(err) {
		httperr.UnprocessableEntity("invalid-deadline", err, w, r)

		return
	}

	if course.IsInvalidTestDataError(err) {
		httperr.UnprocessableEntity("invalid-test-data", err, w, r)

		return
	}

	if course.IsInvalidTestPointError(err) {
		httperr.UnprocessableEntity("invalid-test-points", err, w, r)

		return
	}

	if course.IsAcademicCantEditCourseError(err) {
		httperr.Forbidden("academic-cant-edit-course", err, w, r)

		return
	}

	httperr.InternalServerError("unexpected-error", err, w, r)
}

func (h handler) GetCourseTasks(w http.ResponseWriter, r *http.Request, courseID string, params GetCourseTasksParams) {
	qry, ok := unmarshalAllTasksQuery(w, r, courseID, params)
	if !ok {
		return
	}

	tasks, err := h.app.Queries.AllTasks.Handle(r.Context(), qry)
	if err == nil {
		marshalGeneralTasks(w, r, tasks)

		return
	}

	if errors.Is(err, app.ErrCourseDoesntExist) {
		httperr.NotFound("course-not-found", err, w, r)

		return
	}

	httperr.InternalServerError("unexpected-error", err, w, r)
}

func (h handler) GetCourseTask(w http.ResponseWriter, r *http.Request, courseID string, taskNumber int) {
	qry, ok := unmarshalSpecificTaskQuery(w, r, courseID, taskNumber)
	if !ok {
		return
	}

	task, err := h.app.Queries.SpecificTask.Handle(r.Context(), qry)
	if err == nil {
		marshalSpecificTask(w, r, task)

		return
	}

	if errors.Is(err, app.ErrCourseDoesntExist) {
		httperr.NotFound("course-not-found", err, w, r)

		return
	}

	if errors.Is(err, app.ErrTaskDoesntExist) {
		httperr.NotFound("course-task-not-found", err, w, r)

		return
	}

	httperr.InternalServerError("unexpected-error", err, w, r)
}
