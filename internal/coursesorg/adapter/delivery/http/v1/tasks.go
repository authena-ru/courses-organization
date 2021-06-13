package v1

import "net/http"

func (s HTTPServer) GetCourseTasks(w http.ResponseWriter, r *http.Request, courseId string, params GetCourseTasksParams) {

}

func (s HTTPServer) CreateTaskInCourse(w http.ResponseWriter, r *http.Request, courseId string) {

}

func (s HTTPServer) GetCourseTask(w http.ResponseWriter, r *http.Request, courseId string, taskNumber int) {

}
