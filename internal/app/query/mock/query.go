package mock

import (
	"context"

	"github.com/authena-ru/courses-organization/internal/app"
)

type AllTasksHandler func(ctx context.Context, qry app.AllTasksQuery) ([]app.GeneralTask, error)

func (m AllTasksHandler) Handle(ctx context.Context, qry app.AllTasksQuery) ([]app.GeneralTask, error) {
	return m(ctx, qry)
}

type SpecificTaskHandler func(ctx context.Context, qry app.SpecificTaskQuery) (app.SpecificTask, error)

func (m SpecificTaskHandler) Handle(ctx context.Context, qry app.SpecificTaskQuery) (app.SpecificTask, error) {
	return m(ctx, qry)
}
