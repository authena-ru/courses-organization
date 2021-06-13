package v1

import "net/http"

func (s HTTPServer) GetAllCourses(w http.ResponseWriter, r *http.Request, params GetAllCoursesParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s HTTPServer) GetCourse(w http.ResponseWriter, r *http.Request, courseId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s HTTPServer) CreateCourse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s HTTPServer) EditCourse(w http.ResponseWriter, r *http.Request, courseId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s HTTPServer) ExtendCourse(w http.ResponseWriter, r *http.Request, courseId string) {
	w.WriteHeader(http.StatusNotImplemented)
}
