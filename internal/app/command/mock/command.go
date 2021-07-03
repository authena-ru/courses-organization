package mock

import (
	"context"

	"github.com/authena-ru/courses-organization/internal/app"
)

type AddCollaboratorsHandler func(ctx context.Context, cmd app.AddCollaboratorCommand) error

func (m AddCollaboratorsHandler) Handle(ctx context.Context, cmd app.AddCollaboratorCommand) error {
	return m(ctx, cmd)
}

type RemoveCollaboratorHandler func(ctx context.Context, cmd app.RemoveCollaboratorCommand) error

func (m RemoveCollaboratorHandler) Handle(ctx context.Context, cmd app.RemoveCollaboratorCommand) error {
	return m(ctx, cmd)
}

type CreateCourseHandler func(ctx context.Context, cmd app.CreateCourseCommand) (string, error)

func (m CreateCourseHandler) Handle(ctx context.Context, cmd app.CreateCourseCommand) (string, error) {
	return m(ctx, cmd)
}

type ExtendCourseHandler func(ctx context.Context, cmd app.ExtendCourseCommand) (string, error)

func (m ExtendCourseHandler) Handle(ctx context.Context, cmd app.ExtendCourseCommand) (string, error) {
	return m(ctx, cmd)
}

type AddStudentHandler func(ctx context.Context, cmd app.AddStudentCommand) error

func (m AddStudentHandler) Handle(ctx context.Context, cmd app.AddStudentCommand) error {
	return m(ctx, cmd)
}
