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
	ID        string
	CreatorID string
	Title     string
	Period    Period
	Started   bool
}

var (
	ErrEmptyCourseID    = errors.New("empty course id")
	ErrEmptyCreatorID   = errors.New("empty course creator id")
	ErrEmptyCourseTitle = errors.New("empty course title")
	ErrZeroCoursePeriod = errors.New("zero course period")
)

func NewCourse(params CreationCourseParams) (*Course, error) {
	if params.ID == "" {
		return nil, ErrEmptyCourseID
	}
	if params.CreatorID == "" {
		return nil, ErrEmptyCreatorID
	}
	if params.Title == "" {
		return nil, ErrEmptyCourseTitle
	}
	if params.Period.IsZero() {
		return nil, ErrZeroCoursePeriod
	}
	return &Course{
		id:        params.ID,
		creatorID: params.CreatorID,
		title:     params.Title,
		period:    params.Period,
		started:   params.Started,
	}, nil
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

func (c *Course) Collaborators() []string {
	collaborators := make([]string, 0, len(c.collaborators))
	for c := range c.collaborators {
		collaborators = append(collaborators, c)
	}
	return collaborators
}

func (c *Course) Students() []string {
	students := make([]string, 0, len(c.students))
	for s := range c.students {
		students = append(students, s)
	}
	return students
}
