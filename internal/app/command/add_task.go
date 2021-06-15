package command

import (
	"context"

	"github.com/pkg/errors"

	"github.com/authena-ru/courses-organization/internal/domain/course"
)

type AddTaskCommand struct {
	Academic        course.Academic
	CourseID        string
	TaskTitle       string
	TaskDescription string
	TaskType        course.TaskType
	Deadline        course.Deadline
	TestPoints      []course.TestPoint
	TestData        []course.TestData
}

type AddTaskHandler struct {
	coursesRepository coursesRepository
}

func NewAddTaskHandler(repository coursesRepository) AddTaskHandler {
	if repository == nil {
		panic("coursesRepository is nil")
	}
	return AddTaskHandler{coursesRepository: repository}
}

// Handle is AddTaskCommand handler.
// Adds task with manual checking, auto code checking or testing type,
// returns one of possible errors: app.ErrCourseDoesntExist, course.ErrTaskTitleTooLong,
// app.ErrDatabaseProblems, course.ErrTaskDescriptionTooLong, error that can be detected
// using method course.IsAcademicCantEditCourseError and others without definition.
func (h AddTaskHandler) Handle(ctx context.Context, cmd AddTaskCommand) (taskNumber int, err error) {
	defer func() {
		err = errors.Wrapf(
			err,
			"adding %s task to course #%s by academic #%s",
			cmd.TaskType, cmd.CourseID, cmd.Academic.ID(),
		)
	}()

	givenTaskNumber := new(int)
	if err := h.coursesRepository.UpdateCourse(ctx, cmd.CourseID, addTask(cmd, givenTaskNumber)); err != nil {
		return 0, err
	}
	return *givenTaskNumber, nil
}

var errInvalidTaskType = errors.New("invalid task type")

func addTask(cmd AddTaskCommand, givenTaskNumber *int) UpdateFunction {
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
