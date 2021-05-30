package application

import "github.com/authena-ru/courses-organization/internal/coursesorg/application/command"

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateCourse    command.CreateCourseHandler
	AddCollaborator command.AddCollaboratorHandler
	AddStudent      command.AddStudentHandler
}

type Queries struct {
}
