package v1

import "net/http"

func (h handler) AttachAuxiliaryMaterialToCourse(w http.ResponseWriter, _ *http.Request, _ string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h handler) GetAllCourseAuxiliaryMaterials(w http.ResponseWriter, _ *http.Request, _ string, _ GetAllCourseAuxiliaryMaterialsParams) {
	w.WriteHeader(http.StatusNotImplemented)
}
