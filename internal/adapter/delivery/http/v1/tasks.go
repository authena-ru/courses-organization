package v1

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/authena-ru/courses-organization/internal/adapter/delivery/http/httperr"
	"github.com/authena-ru/courses-organization/internal/app"
	"github.com/authena-ru/courses-organization/internal/domain/course"
)

func (h handler) GetCourseTasks(w http.ResponseWriter, r *http.Request, courseId string, params GetCourseTasksParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h handler) AddTaskToCourse(w http.ResponseWriter, r *http.Request, courseID string) {
	cmd, ok := unmarshallAddTaskCommand(w, r, courseID)
	if !ok {
		return
	}
	taskNumber, err := h.app.Commands.AddTask.Handle(r.Context(), cmd)
	if err == nil {
		w.Header().Set("Content-Location", fmt.Sprintf("courses/%s/tasks/%d", courseID, taskNumber))
		w.WriteHeader(http.StatusCreated)
		return
	}
	if errors.Is(err, app.ErrCourseDoesntExist) {
		httperr.NotFound("course-not-found", err, w, r)
		return
	}
	if course.IsInvalidTaskParametersError(err) {
		httperr.BadRequest("invalid-task-parameters", err, w, r)
		return
	}
	if course.IsInvalidDeadlineError(err) {
		httperr.BadRequest("invalid-deadline", err, w, r)
		return
	}
	if course.IsInvalidTestDataError(err) {
		httperr.BadRequest("invalid-test-data", err, w, r)
		return
	}
	if course.IsInvalidTestPointError(err) {
		httperr.BadRequest("invalid-test-points", err, w, r)
		return
	}
	if course.IsAcademicCantEditCourseError(err) {
		httperr.Forbidden("academic-cant-edit-course", err, w, r)
		return
	}
	httperr.InternalServerError("unexpected-error", err, w, r)
}

func (h handler) GetCourseTask(w http.ResponseWriter, r *http.Request, courseId string, taskNumber int) {
	w.WriteHeader(http.StatusNotImplemented)
}
