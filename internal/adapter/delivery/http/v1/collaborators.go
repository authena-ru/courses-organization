package v1

import "net/http"

func (h handler) GetAllCourseCollaborators(w http.ResponseWriter, r *http.Request, courseId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h handler) AddCollaboratorToCourse(w http.ResponseWriter, r *http.Request, courseId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h handler) RemoveCollaboratorFromCourse(w http.ResponseWriter, r *http.Request, courseId string, teacherId string) {
	w.WriteHeader(http.StatusNotImplemented)
}
