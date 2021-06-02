package course

import "github.com/pkg/errors"

type TaskType uint8

const (
	ManualChecking TaskType = iota + 1
	AutoCodeChecking
	Testing
)

type taskOptional struct {
	deadline   Deadline
	testPoints []TestPoint
	testData   []TestData
}

type Task struct {
	number      int
	title       string
	description string
	taskType    TaskType
	optional    taskOptional
}

func (t *Task) Number() int {
	return t.number
}

func (t *Task) Title() string {
	return t.title
}

func (t *Task) Description() string {
	return t.description
}

func (t *Task) Type() TaskType {
	return t.taskType
}

func (t *Task) ManualCheckingOptional() Deadline {
	return t.optional.deadline
}

func (t *Task) AutoCodeCheckingOptional() (Deadline, []TestData) {
	testDataCopy := make([]TestData, len(t.optional.testData))
	copy(testDataCopy, t.optional.testData)
	return t.optional.deadline, testDataCopy
}

func (t *Task) TestingOptional() []TestPoint {
	testPointsCopy := make([]TestPoint, len(t.optional.testPoints))
	copy(testPointsCopy, t.optional.testPoints)
	return testPointsCopy
}

const (
	taskTitleMaxLen       = 200
	taskDescriptionMaxLen = 1000
)

var (
	ErrTaskHasNoDeadline      = errors.New("task has no deadline")
	ErrTaskHasNoTestPoints    = errors.New("task has no Task points")
	ErrTaskHasNoTestData      = errors.New("task has no test data")
	ErrTaskTitleTooLong       = errors.New("task title too long")
	ErrTaskDescriptionTooLong = errors.New("task description too long")
	ErrTaskNumberOutOfBounds  = errors.New("task number out of bounds")
)

func (t *Task) rename(title string) error {
	if len(title) > taskTitleMaxLen {
		return ErrTaskTitleTooLong
	}
	t.title = title
	return nil
}

func (t *Task) replaceDescription(description string) error {
	if len(description) > taskDescriptionMaxLen {
		return ErrTaskDescriptionTooLong
	}
	t.description = description
	return nil
}

func (t *Task) replaceDeadline(deadline Deadline) error {
	if t.taskType == Testing {
		return ErrTaskHasNoDeadline
	}
	t.optional.deadline = deadline
	return nil
}

func (t *Task) replaceTestPoints(testPoints []TestPoint) error {
	if t.taskType != Testing {
		return ErrTaskHasNoTestPoints
	}
	t.optional.testPoints = testPoints
	return nil
}

func (t *Task) replaceTestData(testData []TestData) error {
	if t.taskType != AutoCodeChecking {
		return ErrTaskHasNoTestData
	}
	t.optional.testData = testData
	return nil
}

func (c *Course) Task(taskNumber int) (Task, error) {
	if taskNumber >= len(c.tasks) {
		return Task{}, ErrTaskNumberOutOfBounds
	}
	return c.tasks[taskNumber], nil
}

type ManualCheckingTaskCreationParams struct {
	Title       string
	Description string
	Deadline    Deadline
}

func (c *Course) AddManualCheckingTask(academic Academic, params ManualCheckingTaskCreationParams) (int, error) {
	if err := c.CanAcademicEditWithAccess(academic, TeacherAccess); err != nil {
		return 0, err
	}
	task, err := c.newTask(params.Title, params.Description, ManualChecking)
	if err != nil {
		return 0, err
	}
	task.optional = taskOptional{deadline: params.Deadline}
	return task.number, nil
}

type AutoCodeCheckingTaskCreationParams struct {
	Title       string
	Description string
	Deadline    Deadline
	TestData    []TestData
}

func (c *Course) AddAutoCodeCheckingTask(academic Academic, params AutoCodeCheckingTaskCreationParams) (int, error) {
	if err := c.CanAcademicEditWithAccess(academic, TeacherAccess); err != nil {
		return 0, err
	}
	task, err := c.newTask(params.Title, params.Description, AutoCodeChecking)
	if err != nil {
		return 0, err
	}
	testDataCopy := make([]TestData, len(params.TestData))
	copy(testDataCopy, params.TestData)
	task.optional = taskOptional{
		deadline: params.Deadline,
		testData: testDataCopy,
	}
	return task.number, nil
}

type TestingTaskCreationParams struct {
	Title       string
	Description string
	TestPoints  []TestPoint
}

func (c *Course) AddTestingTask(academic Academic, params TestingTaskCreationParams) (int, error) {
	if err := c.CanAcademicEditWithAccess(academic, TeacherAccess); err != nil {
		return 0, err
	}
	task, err := c.newTask(params.Title, params.Description, Testing)
	if err != nil {
		return 0, err
	}
	testPointsCopy := make([]TestPoint, len(params.TestPoints))
	task.optional = taskOptional{testPoints: testPointsCopy}
	return task.number, nil
}

func (c *Course) RenameTask(academic Academic, taskNumber int, title string) error {
	if err := c.CanAcademicEditWithAccess(academic, TeacherAccess); err != nil {
		return err
	}
	task, err := c.Task(taskNumber)
	if err != nil {
		return err
	}
	return task.rename(title)
}

func (c *Course) ReplaceTaskDescription(academic Academic, taskNumber int, description string) error {
	if err := c.CanAcademicEditWithAccess(academic, TeacherAccess); err != nil {
		return err
	}
	task, err := c.Task(taskNumber)
	if err != nil {
		return err
	}
	return task.replaceDescription(description)
}

func (c *Course) ReplaceTaskDeadline(academic Academic, taskNumber int, deadline Deadline) error {
	if err := c.CanAcademicEditWithAccess(academic, TeacherAccess); err != nil {
		return err
	}
	task, err := c.Task(taskNumber)
	if err != nil {
		return err
	}
	return task.replaceDeadline(deadline)
}

func (c *Course) ReplaceTaskTestPoints(academic Academic, taskNumber int, testPoints []TestPoint) error {
	if err := c.CanAcademicEditWithAccess(academic, TeacherAccess); err != nil {
		return err
	}
	task, err := c.Task(taskNumber)
	if err != nil {
		return err
	}
	return task.replaceTestPoints(testPoints)
}

func (c *Course) ReplaceTaskTestData(academic Academic, taskNumber int, testData []TestData) error {
	if err := c.CanAcademicEditWithAccess(academic, TeacherAccess); err != nil {
		return err
	}
	task, err := c.Task(taskNumber)
	if err != nil {
		return err
	}
	return task.replaceTestData(testData)
}

func (c *Course) TasksNumber() int {
	return len(c.tasks)
}

func (c *Course) newTask(title string, description string, taskType TaskType) (Task, error) {
	task := Task{
		number:   c.nextTaskNumber,
		taskType: taskType,
	}
	if err := task.rename(title); err != nil {
		return Task{}, err
	}
	if err := task.replaceDescription(description); err != nil {
		return Task{}, err
	}
	c.tasks[c.nextTaskNumber] = task
	c.nextTaskNumber++
	return task, nil
}
