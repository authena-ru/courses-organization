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
