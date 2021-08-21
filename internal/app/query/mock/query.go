package mock

import (
	"context"

	"github.com/authena-ru/courses-organization/internal/app"
)

type AllTasksHandler func(ctx context.Context, qry app.AllTasksQuery) ([]app.GeneralTask, error)

func (m AllTasksHandler) Handle(ctx context.Context, qry app.AllTasksQuery) ([]app.GeneralTask, error) {
	return m(ctx, qry)
}
