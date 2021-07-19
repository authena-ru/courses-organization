package v1

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/authena-ru/courses-organization/internal/app"
	"github.com/authena-ru/courses-organization/internal/domain/course"
	"github.com/authena-ru/courses-organization/pkg/httperr"
)

func (h handler) GetAllCourseStudents(w http.ResponseWriter, _ *http.Request, _ string, _ GetAllCourseStudentsParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h handler) AddStudentToCourse(w http.ResponseWriter, r *http.Request, courseID string) {
	cmd, ok := unmarshalAddStudentCommand(w, r, courseID)
	if !ok {
		return
	}

	err := h.app.Commands.AddStudent.Handle(r.Context(), cmd)

	if err == nil {
		w.WriteHeader(http.StatusNoContent)

		return
	}

	if errors.Is(err, app.ErrCourseDoesntExist) {
		httperr.NotFound("course-not-found", err, w, r)

		return
	}

	if errors.Is(err, app.ErrStudentDoesntExist) {
		httperr.UnprocessableEntity("student-not-found", err, w, r)

		return
	}

	if course.IsAcademicCantEditCourseError(err) {
		httperr.Forbidden("academic-cant-edit-course", err, w, r)

		return
	}

	httperr.InternalServerError("unexpected-error", err, w, r)
}

func (h handler) RemoveStudentFromCourse(w http.ResponseWriter, r *http.Request, courseID string, studentID string) {
	cmd, ok := unmarshalRemoveStudentCommand(w, r, courseID, studentID)
	if !ok {
		return
	}

	err := h.app.Commands.RemoveStudent.Handle(r.Context(), cmd)
	if err == nil {
		w.WriteHeader(http.StatusNoContent)

		return
	}

	if errors.Is(err, app.ErrCourseDoesntExist) {
		httperr.NotFound("course-not-found", err, w, r)

		return
	}

	if errors.Is(err, course.ErrCourseHasNoSuchStudent) {
		httperr.NotFound("course-student-not-found", err, w, r)

		return
	}

	if course.IsAcademicCantEditCourseError(err) {
		httperr.Forbidden("academic-cant-edit-course", err, w, r)

		return
	}

	httperr.InternalServerError("unexpected-error", err, w, r)
}

func (h handler) AddGroupToCourse(w http.ResponseWriter, _ *http.Request, _ string) {
	w.WriteHeader(http.StatusNotImplemented)
}
