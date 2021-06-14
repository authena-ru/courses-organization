package v1

import "net/http"

func (h coursesOrganizationHandler) GetCourseTasks(w http.ResponseWriter, r *http.Request, courseId string, params GetCourseTasksParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h coursesOrganizationHandler) CreateTaskInCourse(w http.ResponseWriter, r *http.Request, courseId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h coursesOrganizationHandler) GetCourseTask(w http.ResponseWriter, r *http.Request, courseId string, taskNumber int) {
	w.WriteHeader(http.StatusNotImplemented)
}
