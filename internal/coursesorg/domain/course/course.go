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

var (
	ErrEmptyCourseID    = errors.New("empty course id")
	ErrEmptyCreatorID   = errors.New("empty creator id")
	ErrEmptyCourseTitle = errors.New("empty course title")
	ErrZeroCoursePeriod = errors.New("zero course period")
)

func NewCourse(id string, creatorID string, title string, period Period) (*Course, error) {
	if id == "" {
		return nil, ErrEmptyCourseID
	}
	if creatorID == "" {
		return nil, ErrEmptyCreatorID
	}
	if title == "" {
		return nil, ErrEmptyCourseTitle
	}
	if period.IsZero() {
		return nil, ErrZeroCoursePeriod
	}
	return &Course{
		id:        id,
		creatorID: creatorID,
		title:     title,
		period:    period,
	}, nil
}
