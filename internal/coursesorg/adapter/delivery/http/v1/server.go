package v1

import "github.com/authena-ru/courses-organization/internal/coursesorg/app"

type HTTPServer struct {
	app app.Application
}

func NewHTTPServer(app app.Application) HTTPServer {
	return HTTPServer{app: app}
}
