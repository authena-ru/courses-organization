package command

import (
	"context"

	"github.com/pkg/errors"

	"github.com/authena-ru/courses-organization/internal/app"
	"github.com/authena-ru/courses-organization/internal/domain/course"
)

type AddTaskHandler struct {
	coursesRepository coursesRepository
}

func NewAddTaskHandler(repository coursesRepository) AddTaskHandler {
	if repository == nil {
		panic("coursesRepository is nil")
	}

	return AddTaskHandler{coursesRepository: repository}
}

func (h AddTaskHandler) Handle(ctx context.Context, cmd app.AddTaskCommand) (taskNumber int, err error) {
	err = h.coursesRepository.UpdateCourse(ctx, cmd.CourseID, addTask(cmd, &taskNumber))

	return taskNumber, errors.Wrapf(
		err,
		"adding %s task to course #%s by academic #%s",
		cmd.TaskType, cmd.CourseID, cmd.Academic.ID(),
	)
}

var errInvalidTaskType = errors.New("invalid task type")

func addTask(cmd app.AddTaskCommand, givenTaskNumber *int) UpdateFunction {
	return func(_ context.Context, crs *course.Course) (*course.Course, error) {
		var (
			number int
			err    error
		)

		switch cmd.TaskType {
		case course.ManualCheckingType:
			number, err = crs.AddManualCheckingTask(cmd.Academic, course.ManualCheckingTaskCreationParams{
				Title:       cmd.TaskTitle,
				Description: cmd.TaskDescription,
				Deadline:    cmd.Deadline,
			})
		case course.AutoCodeCheckingType:
			number, err = crs.AddAutoCodeCheckingTask(cmd.Academic, course.AutoCodeCheckingTaskCreationParams{
				Title:       cmd.TaskTitle,
				Description: cmd.TaskDescription,
				Deadline:    cmd.Deadline,
				TestData:    cmd.TestData,
			})
		case course.TestingType:
			number, err = crs.AddTestingTask(cmd.Academic, course.TestingTaskCreationParams{
				Title:       cmd.TaskTitle,
				Description: cmd.TaskDescription,
				TestPoints:  cmd.TestPoints,
			})
		default:
			number, err = 0, errInvalidTaskType
		}

		*givenTaskNumber = number

		if err != nil {
			return nil, err
		}

		return crs, nil
	}
}
