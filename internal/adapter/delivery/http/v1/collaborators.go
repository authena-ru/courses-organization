package v1

import "net/http"

func (h coursesOrganizationHandler) GetAllCourseCollaborators(w http.ResponseWriter, r *http.Request, courseId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h coursesOrganizationHandler) AddCollaboratorToCourse(w http.ResponseWriter, r *http.Request, courseId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h coursesOrganizationHandler) RemoveCollaboratorFromCourse(w http.ResponseWriter, r *http.Request, courseId string, teacherId string) {
	w.WriteHeader(http.StatusNotImplemented)
}
