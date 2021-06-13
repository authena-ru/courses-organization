package v1

import "net/http"

func (s HTTPServer) AttachAuxiliaryMaterialToCourse(w http.ResponseWriter, r *http.Request, courseId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s HTTPServer) GetAllCourseAuxiliaryMaterials(w http.ResponseWriter, r *http.Request, courseId string, params GetAllCourseAuxiliaryMaterialsParams) {
	w.WriteHeader(http.StatusNotImplemented)
}
