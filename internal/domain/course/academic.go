package course

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

type AcademicType uint8

const (
	TeacherType AcademicType = iota + 1
	StudentType
)

func (at AcademicType) String() string {
	switch at {
	case TeacherType:
		return "teacher"
	case StudentType:
		return "student"
	}

	return "%!AcademicType(" + strconv.Itoa(int(at)) + ")"
}

func NewAcademicTypeFromString(value string) AcademicType {
	switch value {
	case "teacher":
		return TeacherType
	case "student":
		return StudentType
	}

	return AcademicType(0)
}

type Academic struct {
	t  AcademicType
	id string
}

func (at AcademicType) IsValid() bool {
	switch at {
	case TeacherType, StudentType:
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
		return "`teacher` access"
	case CreatorAccess:
		return "`creator` access"
	}

	return "%!Access(" + strconv.Itoa(int(a)) + ")"
}

type academicCantEditCourseError struct {
	academicType AcademicType
	access       Access
}

func (e academicCantEditCourseError) Error() string {
	if e.academicType == StudentType {
		return "student can't edit course"
	}

	return fmt.Sprintf("teacher can't edit course with %s", e.access)
}

func IsAcademicCantEditCourseError(err error) bool {
	var e academicCantEditCourseError

	return errors.As(err, &e)
}

func (c *Course) canAcademicEditWithAccess(academic Academic, access Access) error {
	if academic.Type() == TeacherType {
		if access == TeacherAccess && c.hasTeacher(academic.ID()) {
			return nil
		}

		if access == CreatorAccess && c.hasCreator(academic.ID()) {
			return nil
		}

		return academicCantEditCourseError{academicType: TeacherType, access: access}
	}

	return academicCantEditCourseError{academicType: StudentType}
}

func (a Academic) canCreateCourse() error {
	if a.Type() == TeacherType {
		return nil
	}

	return ErrNotTeacherCantCreateCourse
}
