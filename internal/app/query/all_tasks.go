package query

import (
	"context"
	"github.com/authena-ru/courses-organization/internal/app"

	"github.com/pkg/errors"

	"github.com/authena-ru/courses-organization/internal/domain/course"
)

type TasksFilterParams struct {
	Type course.TaskType
	Text string
}

type allTasksReadModel interface {
	FindAllTasks(
		ctx context.Context,
		academic course.Academic, courseID string,
		filterParams TasksFilterParams,
	) ([]app.GeneralTask, error)
}

type AllTasksHandler struct {
	readModel allTasksReadModel
}

func NewAllTasksHandler(readModel allTasksReadModel) AllTasksHandler {
	if readModel == nil {
		panic("readModel is nil")
	}
	return AllTasksHandler{readModel: readModel}
}

func (h AllTasksHandler) Handle(ctx context.Context, qry app.AllTasksQuery) ([]app.GeneralTask, error) {
	tasks, err := h.readModel.FindAllTasks(ctx, qry.Academic, qry.CourseID, TasksFilterParams{
		Type: qry.Type,
		Text: qry.Text,
	})
	return tasks, errors.Wrapf(err, "getting all tasks of course #%s", qry.CourseID)
}
