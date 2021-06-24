package query

import (
	"context"

	"github.com/pkg/errors"

	"github.com/authena-ru/courses-organization/internal/domain/course"
)

type AllTasksQuery struct {
	Academic course.Academic
	CourseID string
	Type     course.TaskType
	Text     string
}

type TasksFilterParams struct {
	Type course.TaskType
	Text string
}

type allTasksReadModel interface {
	FindAllTasks(
		ctx context.Context,
		academic course.Academic, courseID string,
		filterParams TasksFilterParams,
	) ([]GeneralTask, error)
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

// Handle is AllTasksQuery handler.
// Returns list of course tasks with general task parameters.
// Tasks filtered by type, title and description.
// If course doesn't exist, error equal app.ErrCourseDoesntExist
func (h AllTasksHandler) Handle(ctx context.Context, qry AllTasksQuery) ([]GeneralTask, error) {
	tasks, err := h.readModel.FindAllTasks(ctx, qry.Academic, qry.CourseID, TasksFilterParams{
		Type: qry.Type,
		Text: qry.Text,
	})
	return tasks, errors.Wrapf(err, "getting all tasks of course #%s", qry.CourseID)
}
