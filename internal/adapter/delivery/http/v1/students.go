package v1

import "net/http"

func (h handler) GetAllCourseStudents(w http.ResponseWriter, r *http.Request, courseId string, params GetAllCourseStudentsParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h handler) AddStudentToCourse(w http.ResponseWriter, r *http.Request, courseId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h handler) RemoveStudentFromCourse(w http.ResponseWriter, r *http.Request, courseId string, studentId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h handler) AddGroupToCourse(w http.ResponseWriter, r *http.Request, courseId string) {
	w.WriteHeader(http.StatusNotImplemented)
}
