package course

import "github.com/pkg/errors"

type Course struct {
	id      string
	title   string
	period  Period
	started bool

	creatorID     string
	collaborators map[string]bool
	students      map[string]bool

	tasks          map[int]*Task
	nextTaskNumber int
}

type CreationParams struct {
	ID            string
	Creator       Academic
	Title         string
	Period        Period
	Started       bool
	Collaborators []string
	Students      []string
}

var (
	ErrEmptyCourseID    = errors.New("empty course id")
	ErrZeroCreator      = errors.New("empty course creator id")
	ErrEmptyCourseTitle = errors.New("empty course title")
	ErrZeroCoursePeriod = errors.New("zero course period")
)

func IsInvalidCourseParametersError(err error) bool {
	return errors.Is(err, ErrEmptyCourseID) ||
		errors.Is(err, ErrZeroCreator) ||
		errors.Is(err, ErrEmptyCourseTitle) ||
		errors.Is(err, ErrZeroCoursePeriod)
}

func NewCourse(params CreationParams) (*Course, error) {
	if params.ID == "" {
		return nil, ErrEmptyCourseID
	}
	if params.Creator.IsZero() {
		return nil, ErrZeroCreator
	}
	if err := params.Creator.canCreateCourse(); err != nil {
		return nil, err
	}
	if params.Title == "" {
		return nil, ErrEmptyCourseTitle
	}
	if params.Period.IsZero() {
		return nil, ErrZeroCoursePeriod
	}
	crs := &Course{
		id:             params.ID,
		creatorID:      params.Creator.ID(),
		title:          params.Title,
		period:         params.Period,
		started:        params.Started,
		collaborators:  make(map[string]bool, len(params.Collaborators)),
		students:       make(map[string]bool, len(params.Students)),
		tasks:          make(map[int]*Task),
		nextTaskNumber: 1,
	}
	for _, c := range params.Collaborators {
		crs.collaborators[c] = true
	}
	for _, s := range params.Students {
		crs.students[s] = true
	}
	return crs, nil
}

func (c *Course) Extend(params CreationParams) (*Course, error) {
	if params.ID == "" {
		return nil, ErrEmptyCourseID
	}
	if params.Creator.IsZero() {
		return nil, ErrZeroCreator
	}
	if err := params.Creator.canCreateCourse(); err != nil {
		return nil, err
	}
	if err := c.canAcademicEditWithAccess(params.Creator, TeacherAccess); err != nil {
		return nil, err
	}
	extendedCourseTitle := c.Title()
	if params.Title != "" {
		extendedCourseTitle = params.Title
	}
	extendedCoursePeriod := c.period.next()
	if !params.Period.IsZero() {
		extendedCoursePeriod = params.Period
	}
	crs := &Course{
		id:             params.ID,
		creatorID:      params.Creator.ID(),
		title:          extendedCourseTitle,
		period:         extendedCoursePeriod,
		started:        params.Started,
		collaborators:  unmarshallIDs(append(c.Collaborators(), params.Collaborators...)),
		students:       unmarshallIDs(append(c.Students(), params.Students...)),
		tasks:          make(map[int]*Task, len(c.tasks)),
		nextTaskNumber: len(c.tasks) + 1,
	}
	for i, t := range c.tasksCopy() {
		number := i + 1
		crs.tasks[number] = t
		crs.tasks[number].number = number
	}
	return crs, nil
}

func MustNewCourse(params CreationParams) *Course {
	crs, err := NewCourse(params)
	if err != nil {
		panic(err)
	}
	return crs
}

func (c *Course) ID() string {
	return c.id
}

func (c *Course) Title() string {
	return c.title
}

func (c *Course) Period() Period {
	return c.period
}

func (c *Course) Started() bool {
	return c.started
}

func (c *Course) CreatorID() string {
	return c.creatorID
}

type UnmarshallingParams struct {
	ID            string
	Title         string
	Period        Period
	Started       bool
	CreatorID     string
	Collaborators []string
	Students      []string
	Tasks         []UnmarshallingTaskParams
}

type UnmarshallingTaskParams struct {
	Number      int
	Title       string
	Description string
	TaskType    TaskType
	Deadline    Deadline
	TestPoints  []TestPoint
	TestData    []TestData
}

// UnmarshallFromDatabase unmarshalls Course from the database.
// It should be used only for unmarshalling from the database!
// Using UnmarshallFromDatabase may put domain into the invalid state!
func UnmarshallFromDatabase(params UnmarshallingParams) *Course {
	tasks, lastNumber := unmarshallTasks(params.Tasks)
	crs := &Course{
		id:             params.ID,
		title:          params.Title,
		period:         params.Period,
		started:        params.Started,
		creatorID:      params.CreatorID,
		collaborators:  unmarshallIDs(params.Collaborators),
		students:       unmarshallIDs(params.Students),
		tasks:          tasks,
		nextTaskNumber: lastNumber + 1,
	}
	return crs
}

func unmarshallIDs(ids []string) map[string]bool {
	unmarshalled := make(map[string]bool, len(ids))
	for _, id := range ids {
		unmarshalled[id] = true
	}
	return unmarshalled
}

func unmarshallTasks(taskParams []UnmarshallingTaskParams) (map[int]*Task, int) {
	tasks := make(map[int]*Task, len(taskParams))
	lastNumber := 0
	for _, tp := range taskParams {
		tasks[tp.Number] = &Task{
			number:      tp.Number,
			title:       tp.Title,
			description: tp.Description,
			taskType:    tp.TaskType,
			optional: taskOptional{
				deadline:   tp.Deadline,
				testData:   tp.TestData,
				testPoints: tp.TestPoints,
			},
		}
		if tp.Number > lastNumber {
			lastNumber = tp.Number
		}
	}
	return tasks, lastNumber
}
