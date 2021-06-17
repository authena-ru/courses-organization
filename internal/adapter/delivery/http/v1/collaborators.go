package v1

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/authena-ru/courses-organization/internal/adapter/delivery/http/httperr"
	"github.com/authena-ru/courses-organization/internal/app"
	"github.com/authena-ru/courses-organization/internal/domain/course"
)

func (h handler) GetAllCourseCollaborators(w http.ResponseWriter, r *http.Request, courseId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h handler) AddCollaboratorToCourse(w http.ResponseWriter, r *http.Request, courseID string) {
	cmd, ok := unmarshallAddCollaboratorCommand(w, r, courseID)
	if !ok {
		return
	}
	err := h.app.Commands.AddCollaborator.Handle(r.Context(), cmd)
	if err == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if errors.Is(err, app.ErrCourseDoesntExist) {
		httperr.NotFound("course-not-found", err, w, r)
		return
	}
	if errors.Is(err, app.ErrTeacherDoesntExist) {
		httperr.UnprocessableEntity("teacher-not-found", err, w, r)
		return
	}
	if course.IsAcademicCantEditCourseError(err) {
		httperr.Forbidden("academic-cant-edit-course", err, w, r)
		return
	}
	httperr.InternalServerError("unexpected-error", err, w, r)
}

func (h handler) RemoveCollaboratorFromCourse(
	w http.ResponseWriter, r *http.Request,
	courseID, teacherID string,
) {
	cmd, ok := unmarshallRemoveCollaboratorCommand(w, r, courseID, teacherID)
	if !ok {
		return
	}
	err := h.app.Commands.RemoveCollaborator.Handle(r.Context(), cmd)
	if err == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if errors.Is(err, app.ErrCourseDoesntExist) {
		httperr.NotFound("course-not-found", err, w, r)
		return
	}
	if errors.Is(err, course.ErrCourseHasNoSuchCollaborator) {
		httperr.NotFound("course-collaborator-not-found", err, w, r)
		return
	}
	if course.IsAcademicCantEditCourseError(err) {
		httperr.Forbidden("academic-cant-edit-course", err, w, r)
		return
	}
	httperr.InternalServerError("unexpected-error", err, w, r)
}
