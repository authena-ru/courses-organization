package v1

import "net/http"

func (h coursesOrganizationHandler) GetAllCourseStudents(w http.ResponseWriter, r *http.Request, courseId string, params GetAllCourseStudentsParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h coursesOrganizationHandler) AddStudentToCourse(w http.ResponseWriter, r *http.Request, courseId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h coursesOrganizationHandler) RemoveStudentFromCourse(w http.ResponseWriter, r *http.Request, courseId string, studentId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h coursesOrganizationHandler) AddGroupToCourse(w http.ResponseWriter, r *http.Request, courseId string) {
	w.WriteHeader(http.StatusNotImplemented)
}
