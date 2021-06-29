package v1

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/authena-ru/courses-organization/internal/adapter/delivery/http/httperr"
	"github.com/authena-ru/courses-organization/internal/app"
	"github.com/authena-ru/courses-organization/internal/domain/course"
)

func (h handler) GetAllCourses(w http.ResponseWriter, _ *http.Request, _ GetAllCoursesParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h handler) GetCourse(w http.ResponseWriter, _ *http.Request, _ string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h handler) CreateCourse(w http.ResponseWriter, r *http.Request) {
	cmd, ok := unmarshallCreateCourseCommand(w, r)
	if !ok {
		return
	}

	createdCourseID, err := h.app.Commands.CreateCourse.Handle(r.Context(), cmd)
	if err == nil {
		w.Header().Set("Content-Location", fmt.Sprintf("/courses/%s", createdCourseID))
		w.WriteHeader(http.StatusCreated)

		return
	}

	if course.IsInvalidCourseParametersError(err) {
		httperr.BadRequest("invalid-course-parameters", err, w, r)

		return
	}

	if errors.Is(err, course.ErrNotTeacherCantCreateCourse) {
		httperr.Forbidden("not-teacher-cant-create-course", err, w, r)

		return
	}

	httperr.InternalServerError("unexpected-error", err, w, r)
}

func (h handler) EditCourse(w http.ResponseWriter, _ *http.Request, _ string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h handler) ExtendCourse(w http.ResponseWriter, r *http.Request, courseID string) {
	cmd, ok := unmarshallExtendCourseCommand(w, r, courseID)
	if !ok {
		return
	}

	extendedCourseID, err := h.app.Commands.ExtendCourse.Handle(r.Context(), cmd)
	if err == nil {
		w.Header().Set("Content-Location", fmt.Sprintf("/courses/%s", extendedCourseID))
		w.WriteHeader(http.StatusCreated)

		return
	}

	if errors.Is(err, app.ErrCourseDoesntExist) {
		httperr.NotFound("course-not-found", err, w, r)

		return
	}

	if course.IsInvalidCourseParametersError(err) {
		httperr.BadRequest("invalid-course-parameters", err, w, r)

		return
	}

	if course.IsAcademicCantEditCourseError(err) {
		httperr.BadRequest("academic-cant-edit-course", err, w, r)

		return
	}

	if errors.Is(err, course.ErrNotTeacherCantCreateCourse) {
		httperr.Forbidden("not-teacher-cant-create-course", err, w, r)

		return
	}

	httperr.InternalServerError("unexpected-error", err, w, r)
}
