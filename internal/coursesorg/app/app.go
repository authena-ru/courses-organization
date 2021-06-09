package app

import "github.com/authena-ru/courses-organization/internal/coursesorg/app/command"

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateCourse       command.CreateCourseHandler
	AddCollaborator    command.AddCollaboratorHandler
	RemoveCollaborator command.RemoveCollaboratorHandler
	AddStudent         command.AddStudentHandler
	RemoveStudent      command.RemoveStudentHandler
}

type Queries struct {
}
