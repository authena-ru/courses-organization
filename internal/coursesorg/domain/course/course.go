package course

import "github.com/pkg/errors"

type Course struct {
	id             string
	title          string
	period         Period
	started        bool
	creatorID      string
	collaborators  map[string]bool
	students       map[string]bool
	nextTaskNumber uint
}

type CreationCourseParams struct {
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

func NewCourse(params CreationCourseParams) (*Course, error) {
	const (
		defaultCollaboratorsNumber = 5
		defaultStudentsNumber      = 10
	)

	if params.ID == "" {
		return nil, ErrEmptyCourseID
	}
	if params.Creator.IsZero() {
		return nil, ErrZeroCreator
	}
	if err := params.Creator.CanCreateCourse(); err != nil {
		return nil, err
	}
	if params.Title == "" {
		return nil, ErrEmptyCourseTitle
	}
	if params.Period.IsZero() {
		return nil, ErrZeroCoursePeriod
	}
	c := &Course{
		id:            params.ID,
		creatorID:     params.Creator.ID(),
		title:         params.Title,
		period:        params.Period,
		started:       params.Started,
		collaborators: make(map[string]bool, defaultCollaboratorsNumber),
		students:      make(map[string]bool, defaultStudentsNumber),
	}
	c.AddStudents(params.Creator, params.Students...)
	c.AddCollaborators(params.Creator, params.Collaborators...)
	return c, nil
}

func MustNewCourse(params CreationCourseParams) *Course {
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
