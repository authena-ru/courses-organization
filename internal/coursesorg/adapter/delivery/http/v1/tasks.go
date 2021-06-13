package v1

import "net/http"

func (s HTTPServer) GetCourseTasks(w http.ResponseWriter, r *http.Request, courseId string, params GetCourseTasksParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s HTTPServer) CreateTaskInCourse(w http.ResponseWriter, r *http.Request, courseId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s HTTPServer) GetCourseTask(w http.ResponseWriter, r *http.Request, courseId string, taskNumber int) {
	w.WriteHeader(http.StatusNotImplemented)
}
