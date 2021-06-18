package query

import (
	"context"
	"github.com/authena-ru/courses-organization/internal/domain/course"
	"github.com/pkg/errors"
)

type SpecificTaskQuery struct {
	Academic   course.Academic
	CourseID   string
	TaskNumber int
}

type specificTaskReadModel interface {
	FindTask(ctx context.Context, academic course.Academic, courseID string, taskNumber int) (SpecificTask, error)
}

type SpecificTaskHandler struct {
	readModel specificTaskReadModel
}

func NewSpecificTaskHandler(readModel specificTaskReadModel) SpecificTaskHandler {
	if readModel == nil {
		panic("readModel is nil")
	}
	return SpecificTaskHandler{readModel: readModel}
}

func (h SpecificTaskHandler) Handle(ctx context.Context, qry SpecificTaskQuery) (SpecificTask, error) {
	task, err := h.readModel.FindTask(ctx, qry.Academic, qry.CourseID, qry.TaskNumber)
	return task, errors.Wrapf(err, "getting task %d of course %s", qry.TaskNumber, qry.CourseID)
}
