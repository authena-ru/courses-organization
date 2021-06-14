package v1

import "net/http"

func (h coursesOrganizationHandler) GetAllCourses(w http.ResponseWriter, r *http.Request, params GetAllCoursesParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h coursesOrganizationHandler) GetCourse(w http.ResponseWriter, r *http.Request, courseId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h coursesOrganizationHandler) CreateCourse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h coursesOrganizationHandler) EditCourse(w http.ResponseWriter, r *http.Request, courseId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h coursesOrganizationHandler) ExtendCourse(w http.ResponseWriter, r *http.Request, courseId string) {
	w.WriteHeader(http.StatusNotImplemented)
}
