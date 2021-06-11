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
		collaborators:  make(map[string]bool, len(params.Collaborators)),
		students:       make(map[string]bool, len(params.Students)),
		tasks:          make(map[int]*Task, len(c.tasks)),
		nextTaskNumber: len(c.tasks) + 1,
	}
	for _, c := range append(c.Collaborators(), params.Collaborators...) {
		crs.collaborators[c] = true
	}
	for _, s := range append(c.Students(), params.Students...) {
		crs.students[s] = true
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
