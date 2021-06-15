package v1

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/authena-ru/courses-organization/internal/adapter/delivery/http/httperr"
	"github.com/authena-ru/courses-organization/internal/app"
	"github.com/authena-ru/courses-organization/internal/app/command"
	"github.com/authena-ru/courses-organization/internal/domain/course"
)

func (h handler) GetAllCourses(w http.ResponseWriter, r *http.Request, params GetAllCoursesParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h handler) GetCourse(w http.ResponseWriter, r *http.Request, courseId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h handler) CreateCourse(w http.ResponseWriter, r *http.Request) {
	academic, ok := unmarshallAcademic(w, r)
	if !ok {
		return
	}
	rb := &CreateCourseRequest{}
	if ok := decode(w, r, rb); !ok {
		return
	}
	period, ok := unmarshallPeriod(w, r, &rb.Period)
	if !ok {
		return
	}
	cmd := command.CreateCourseCommand{
		Academic:      academic,
		CourseStarted: rb.Started,
		CourseTitle:   rb.Title,
		CoursePeriod:  period,
	}

	createdCourseID, err := h.app.Commands.CreateCourse.Handle(r.Context(), cmd)
	if err == nil {
		w.Header().Set("Content-Location", fmt.Sprintf("/courses/%s", createdCourseID))
		w.WriteHeader(http.StatusCreated)
		return
	}
	if errors.Is(err, course.ErrZeroCreator) {
		httperr.BadRequest("zero-creator", err, w, r)
		return
	}
	if errors.Is(err, course.ErrNotTeacherCantCreateCourse) {
		httperr.Forbidden("not-teacher-cant-create-course", err, w, r)
		return
	}
	if errors.Is(err, course.ErrEmptyCourseTitle) {
		httperr.BadRequest("empty-course-title", err, w, r)
		return
	}
	if errors.Is(err, course.ErrZeroCoursePeriod) {
		httperr.BadRequest("zero-course-period", err, w, r)
		return
	}
	httperr.InternalServerError("unexpected-error", err, w, r)
}

func (h handler) EditCourse(w http.ResponseWriter, r *http.Request, courseId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h handler) ExtendCourse(w http.ResponseWriter, r *http.Request, courseId string) {
	academic, ok := unmarshallAcademic(w, r)
	if !ok {
		return
	}
	rb := &ExtendCourseRequest{}
	if ok := decode(w, r, rb); !ok {
		return
	}
	period, ok := unmarshallPeriod(w, r, rb.Period)
	if !ok {
		return
	}
	cmd := command.ExtendCourseCommand{
		Academic:       academic,
		OriginCourseID: courseId,
		CourseStarted:  rb.Started,
		CourseTitle:    rb.Title,
		CoursePeriod:   period,
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
	if errors.Is(err, course.ErrZeroCreator) {
		httperr.BadRequest("zero-creator", err, w, r)
		return
	}
	if errors.Is(err, course.ErrNotTeacherCantCreateCourse) {
		httperr.Forbidden("not-teacher-cant-create-course", err, w, r)
		return
	}
	httperr.InternalServerError("unexpected-error", err, w, r)
}
