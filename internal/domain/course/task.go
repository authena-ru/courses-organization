package course

import (
	"sort"
	"strconv"

	"github.com/pkg/errors"
)

type TaskType uint8

const (
	ManualCheckingType TaskType = iota + 1
	AutoCodeCheckingType
	TestingType
)

func (t TaskType) String() string {
	switch t {
	case ManualCheckingType:
		return "manual checking"
	case AutoCodeCheckingType:
		return "auto code checking"
	case TestingType:
		return "testing"
	}

	return "%!TaskType(" + strconv.Itoa(int(t)) + ")"
}

func (t TaskType) IsValid() bool {
	switch t {
	case ManualCheckingType, AutoCodeCheckingType, TestingType:
		return true
	}

	return false
}

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

func (t *Task) Deadline() (Deadline, bool) {
	if t.taskType == ManualCheckingType ||
		t.taskType == AutoCodeCheckingType {
		return t.deadline(), true
	}

	return Deadline{}, false
}

func (t *Task) TestData() ([]TestData, bool) {
	if t.taskType == AutoCodeCheckingType {
		return t.testData(), true
	}

	return nil, false
}

func (t *Task) TestPoints() ([]TestPoint, bool) {
	if t.taskType == TestingType {
		return t.testPoints(), true
	}

	return nil, false
}

func (t *Task) deadline() Deadline {
	return t.optional.deadline
}

func (t *Task) testData() []TestData {
	testDataCopy := make([]TestData, len(t.optional.testData))
	copy(testDataCopy, t.optional.testData)

	return testDataCopy
}

func (t *Task) testPoints() []TestPoint {
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
	ErrTaskHasNoTestPoints    = errors.New("task has no test points")
	ErrTaskHasNoTestData      = errors.New("task has no test data")
	ErrTaskTitleTooLong       = errors.New("task title too long")
	ErrTaskDescriptionTooLong = errors.New("task description too long")
	ErrCourseHasNoSuchTask    = errors.New("course has no such task")
)

func IsInvalidTaskParametersError(err error) bool {
	return errors.Is(err, ErrTaskTitleTooLong) ||
		errors.Is(err, ErrTaskDescriptionTooLong)
}

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
	if t.taskType == TestingType {
		return ErrTaskHasNoDeadline
	}

	t.optional.deadline = deadline

	return nil
}

func (t *Task) replaceTestPoints(testPoints []TestPoint) error {
	if t.taskType != TestingType {
		return ErrTaskHasNoTestPoints
	}

	t.optional.testPoints = testPoints

	return nil
}

func (t *Task) replaceTestData(testData []TestData) error {
	if t.taskType != AutoCodeCheckingType {
		return ErrTaskHasNoTestData
	}

	t.optional.testData = testData

	return nil
}

func (t *Task) copy() *Task {
	return &Task{
		number:      t.Number(),
		title:       t.Title(),
		description: t.Description(),
		taskType:    t.Type(),
		optional: taskOptional{
			deadline:   Deadline{},
			testPoints: t.testPoints(),
			testData:   t.testData(),
		},
	}
}

func (c *Course) Task(taskNumber int) (Task, error) {
	task, err := c.obtainTask(taskNumber)
	if err != nil {
		return Task{}, err
	}

	return *task, nil
}

func (c *Course) Tasks() []Task {
	tasks := make([]Task, 0, len(c.tasks))
	for _, t := range c.tasks {
		tasks = append(tasks, *t)
	}

	return tasks
}

type ManualCheckingTaskCreationParams struct {
	Title       string
	Description string
	Deadline    Deadline
}

func (c *Course) AddManualCheckingTask(academic Academic, params ManualCheckingTaskCreationParams) (int, error) {
	if err := c.canAcademicEditWithAccess(academic, TeacherAccess); err != nil {
		return 0, err
	}

	task, err := c.newTask(params.Title, params.Description, ManualCheckingType, taskOptional{deadline: params.Deadline})
	if err != nil {
		return 0, err
	}

	return task.number, nil
}

type AutoCodeCheckingTaskCreationParams struct {
	Title       string
	Description string
	Deadline    Deadline
	TestData    []TestData
}

func (c *Course) AddAutoCodeCheckingTask(academic Academic, params AutoCodeCheckingTaskCreationParams) (int, error) {
	if err := c.canAcademicEditWithAccess(academic, TeacherAccess); err != nil {
		return 0, err
	}

	testDataCopy := make([]TestData, len(params.TestData))
	copy(testDataCopy, params.TestData)

	task, err := c.newTask(params.Title, params.Description, AutoCodeCheckingType, taskOptional{
		deadline: params.Deadline,
		testData: testDataCopy,
	})
	if err != nil {
		return 0, err
	}

	return task.number, nil
}

type TestingTaskCreationParams struct {
	Title       string
	Description string
	TestPoints  []TestPoint
}

func (c *Course) AddTestingTask(academic Academic, params TestingTaskCreationParams) (int, error) {
	if err := c.canAcademicEditWithAccess(academic, TeacherAccess); err != nil {
		return 0, err
	}

	testPointsCopy := make([]TestPoint, len(params.TestPoints))
	copy(testPointsCopy, params.TestPoints)

	task, err := c.newTask(params.Title, params.Description, TestingType, taskOptional{testPoints: testPointsCopy})
	if err != nil {
		return 0, err
	}

	return task.number, nil
}

func (c *Course) RenameTask(academic Academic, taskNumber int, title string) error {
	if err := c.canAcademicEditWithAccess(academic, TeacherAccess); err != nil {
		return err
	}

	task, err := c.obtainTask(taskNumber)
	if err != nil {
		return err
	}

	return task.rename(title)
}

func (c *Course) ReplaceTaskDescription(academic Academic, taskNumber int, description string) error {
	if err := c.canAcademicEditWithAccess(academic, TeacherAccess); err != nil {
		return err
	}

	task, err := c.obtainTask(taskNumber)
	if err != nil {
		return err
	}

	return task.replaceDescription(description)
}

func (c *Course) ReplaceTaskDeadline(academic Academic, taskNumber int, deadline Deadline) error {
	if err := c.canAcademicEditWithAccess(academic, TeacherAccess); err != nil {
		return err
	}

	task, err := c.obtainTask(taskNumber)
	if err != nil {
		return err
	}

	return task.replaceDeadline(deadline)
}

func (c *Course) ReplaceTaskTestPoints(academic Academic, taskNumber int, testPoints []TestPoint) error {
	if err := c.canAcademicEditWithAccess(academic, TeacherAccess); err != nil {
		return err
	}

	task, err := c.obtainTask(taskNumber)
	if err != nil {
		return err
	}

	return task.replaceTestPoints(testPoints)
}

func (c *Course) ReplaceTaskTestData(academic Academic, taskNumber int, testData []TestData) error {
	if err := c.canAcademicEditWithAccess(academic, TeacherAccess); err != nil {
		return err
	}

	task, err := c.obtainTask(taskNumber)
	if err != nil {
		return err
	}

	return task.replaceTestData(testData)
}

func (c *Course) TasksNumber() int {
	return len(c.tasks)
}

func (c *Course) newTask(title string, description string, taskType TaskType, optional taskOptional) (*Task, error) {
	task := &Task{
		number:   c.nextTaskNumber,
		taskType: taskType,
		optional: optional,
	}
	if err := task.rename(title); err != nil {
		return nil, err
	}

	if err := task.replaceDescription(description); err != nil {
		return nil, err
	}

	c.tasks[c.nextTaskNumber] = task
	c.nextTaskNumber++

	return task, nil
}

func (c *Course) obtainTask(taskNumber int) (*Task, error) {
	task, ok := c.tasks[taskNumber]
	if !ok {
		return nil, ErrCourseHasNoSuchTask
	}

	return task, nil
}

func (c *Course) tasksCopy() []*Task {
	tasks := make([]*Task, 0, len(c.tasks))
	for _, t := range c.tasks {
		tasks = append(tasks, t.copy())
	}

	sort.SliceStable(tasks, func(i, j int) bool {
		return tasks[i].Number() < tasks[j].Number()
	})

	return tasks
}
