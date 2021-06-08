package course

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

type AcademicType uint8

const (
	Teacher AcademicType = iota + 1
	Student
)

func (at AcademicType) String() string {
	switch at {
	case Teacher:
		return "Teacher"
	case Student:
		return "Student"
	}
	return "%!AcademicType(" + strconv.Itoa(int(at)) + ")"
}

type Academic struct {
	t  AcademicType
	id string
}

func (at AcademicType) IsValid() bool {
	switch at {
	case Teacher, Student:
		return true
	}
	return false
}

var (
	ErrEmptyAcademicID            = errors.New("empty academic id")
	ErrInvalidAcademicType        = errors.New("invalid academic type")
	ErrNotTeacherCantCreateCourse = errors.New("not teacher can't create course")
)

func NewAcademic(id string, t AcademicType) (Academic, error) {
	if id == "" {
		return Academic{}, ErrEmptyAcademicID
	}
	if !t.IsValid() {
		return Academic{}, ErrInvalidAcademicType
	}
	return Academic{t: t, id: id}, nil
}

func MustNewAcademic(id string, t AcademicType) Academic {
	academic, err := NewAcademic(id, t)
	if err != nil {
		panic(err)
	}
	return academic
}

func (a Academic) Type() AcademicType {
	return a.t
}

func (a Academic) ID() string {
	return a.id
}

func (a Academic) IsZero() bool {
	return a == Academic{}
}

type Access uint8

const (
	TeacherAccess Access = iota + 1
	CreatorAccess
)

func (a Access) String() string {
	switch a {
	case TeacherAccess:
		return "teacher access"
	case CreatorAccess:
		return "creator access"
	}
	return "%!Access(" + strconv.Itoa(int(a)) + ")"
}

type academicCantEditCourseError struct {
	message string
}

func (e academicCantEditCourseError) Error() string {
	return e.message
}

func IsAcademicCantEditCourseError(err error) bool {
	var e academicCantEditCourseError
	return errors.As(err, &e)
}

func (c *Course) canAcademicEditWithAccess(academic Academic, access Access) error {
	if academic.Type() == Teacher {
		if access == TeacherAccess && c.hasTeacher(academic.ID()) {
			return nil
		}
		if access == CreatorAccess && c.hasCreator(academic.ID()) {
			return nil
		}
		return academicCantEditCourseError{message: fmt.Sprintf("teacher can't edit course with %s", access)}
	}
	return academicCantEditCourseError{message: "student can't edit course"}
}

func (a Academic) canCreateCourse() error {
	if a.Type() == Teacher {
		return nil
	}
	return ErrNotTeacherCantCreateCourse
}
