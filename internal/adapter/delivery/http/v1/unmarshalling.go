package v1

import (
	"github.com/go-chi/render"
	"net/http"

	"github.com/authena-ru/courses-organization/internal/adapter/delivery/http/auth"
	"github.com/authena-ru/courses-organization/internal/adapter/delivery/http/httperr"
	"github.com/authena-ru/courses-organization/internal/domain/course"
)

func decode(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if err := render.Decode(r, v); err != nil {
		httperr.BadRequest("invalid-request-body", err, w, r)
		return false
	}
	return true
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
