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
	if err := c.canAcademicSee(params.Creator); err != nil {
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
		collaborators:  make(map[string]bool, len(c.collaborators)+len(params.Collaborators)),
		students:       make(map[string]bool, len(c.students)+len(params.Students)),
		tasks:          make(map[int]*Task, len(c.tasks)),
		nextTaskNumber: len(c.tasks) + 1,
	}
	crs.putCollaborators(append(c.Collaborators(), params.Collaborators...))
	crs.putStudents(append(c.Students(), params.Students...))
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

// TODO: umarshall tasks, write doc
func UnmarshallFromDatabase(params UnmarshallingParams) *Course {
	crs := &Course{
		id:             params.ID,
		title:          params.Title,
		period:         params.Period,
		started:        params.Started,
		creatorID:      params.CreatorID,
		collaborators:  make(map[string]bool, len(params.Collaborators)),
		students:       make(map[string]bool, len(params.Students)),
		tasks:          make(map[int]*Task, len(params.Tasks)),
		nextTaskNumber: 1,
	}
	crs.putCollaborators(params.Collaborators)
	crs.putStudents(params.Students)
	return crs
}
