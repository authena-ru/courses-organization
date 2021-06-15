package v1

import "net/http"

func (h handler) AttachAuxiliaryMaterialToCourse(w http.ResponseWriter, r *http.Request, courseId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h handler) GetAllCourseAuxiliaryMaterials(w http.ResponseWriter, r *http.Request, courseId string, params GetAllCourseAuxiliaryMaterialsParams) {
	w.WriteHeader(http.StatusNotImplemented)
}
