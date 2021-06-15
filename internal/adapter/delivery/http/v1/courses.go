package v1

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/pkg/errors"

	"github.com/authena-ru/courses-organization/internal/adapter/delivery/http/auth"
	"github.com/authena-ru/courses-organization/internal/adapter/delivery/http/httperr"
	"github.com/authena-ru/courses-organization/internal/app/command"
	"github.com/authena-ru/courses-organization/internal/domain/course"
)

func (h coursesOrganizationHandler) GetAllCourses(w http.ResponseWriter, r *http.Request, params GetAllCoursesParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h coursesOrganizationHandler) GetCourse(w http.ResponseWriter, r *http.Request, courseId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h coursesOrganizationHandler) CreateCourse(w http.ResponseWriter, r *http.Request) {
	academic, err := auth.AcademicFromCtx(r.Context())
	if err != nil {
		httperr.Unauthorized("no-user-in-context", err, w, r)
		return
	}

	rb := &CreateCourseRequest{}
	if err := render.Decode(r, rb); err != nil {
		httperr.BadRequest("invalid-request-body", err, w, r)
		return
	}

	var semester course.Semester
	switch rb.Period.Semester {
	case SemesterFIRST:
		semester = course.FirstSemester
	case SemesterSECOND:
		semester = course.SecondSemester
	default:
		httperr.BadRequest("invalid-course-period-semester", nil, w, r)
		return
	}
	period, err := course.NewPeriod(rb.Period.AcademicStartYear, rb.Period.AcademicEndYear, semester)
	if err != nil {
		httperr.BadRequest("invalid-course-period", err, w, r)
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

func (h coursesOrganizationHandler) EditCourse(w http.ResponseWriter, r *http.Request, courseId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h coursesOrganizationHandler) ExtendCourse(w http.ResponseWriter, r *http.Request, courseId string) {
	w.WriteHeader(http.StatusNotImplemented)
}
