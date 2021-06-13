package v1

import "net/http"

func (s HTTPServer) GetAllCourseStudents(w http.ResponseWriter, r *http.Request, courseId string, params GetAllCourseStudentsParams) {

}

func (s HTTPServer) AddStudentToCourse(w http.ResponseWriter, r *http.Request, courseId string) {

}

func (s HTTPServer) RemoveStudentFromCourse(w http.ResponseWriter, r *http.Request, courseId string, studentId string) {

}

func (s HTTPServer) AddGroupToCourse(w http.ResponseWriter, r *http.Request, courseId string) {

}
