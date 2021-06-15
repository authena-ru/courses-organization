package v1

import "net/http"

func (h handler) GetCourseTasks(w http.ResponseWriter, r *http.Request, courseId string, params GetCourseTasksParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h handler) CreateTaskInCourse(w http.ResponseWriter, r *http.Request, courseId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h handler) GetCourseTask(w http.ResponseWriter, r *http.Request, courseId string, taskNumber int) {
	w.WriteHeader(http.StatusNotImplemented)
}
