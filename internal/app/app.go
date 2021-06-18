package app

import (
	"github.com/authena-ru/courses-organization/internal/app/command"
	"github.com/authena-ru/courses-organization/internal/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateCourse       command.CreateCourseHandler
	ExtendCourse       command.ExtendCourseHandler
	AddCollaborator    command.AddCollaboratorHandler
	RemoveCollaborator command.RemoveCollaboratorHandler
	AddStudent         command.AddStudentHandler
	RemoveStudent      command.RemoveStudentHandler
	AddTask            command.AddTaskHandler
}

type Queries struct {
	SpecificTask query.SpecificTaskHandler
}
