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

type academicCantSeeCourseError struct {
	requestingAcademicID string
	courseID             string
}

func (e academicCantSeeCourseError) Error() string {
	return fmt.Sprintf("academic #%s can't see course #%s", e.requestingAcademicID, e.courseID)
}

func IsAcademicCantSeeCourseError(err error) bool {
	var e academicCantSeeCourseError
	return errors.As(err, &e)
}

func (c *Course) CanAcademicSee(academic Academic) error {
	if academic.Type() == Teacher && c.hasTeacher(academic.ID()) {
		return nil
	}
	if academic.Type() == Student && c.hasStudent(academic.ID()) {
		return nil
	}
	return academicCantSeeCourseError{
		requestingAcademicID: academic.ID(),
		courseID:             c.ID(),
	}
}

type Access uint8

const (
	TeacherAccess Access = iota + 1
	CreatorAccess
)

func (a Access) String() string {
	switch a {
	case TeacherAccess:
		return "Teacher access"
	case CreatorAccess:
		return "Creator access"
	}
	return "%!Access(" + strconv.Itoa(int(a)) + ")"
}

type academicCantEditCourseError struct {
	requestingAcademicID string
	access               Access
	courseID             string
}

func (e academicCantEditCourseError) Error() string {
	return fmt.Sprintf(
		"academic #%s can't edit course #%s with %s",
		e.requestingAcademicID, e.courseID, e.access,
	)
}

func IsAcademicCantEditCourseError(err error) bool {
	var e academicCantEditCourseError
	return errors.As(err, &e)
}

func (c *Course) CanAcademicEditWithAccess(academic Academic, access Access) error {
	if academic.Type() == Teacher {
		if access == TeacherAccess && c.hasTeacher(academic.ID()) {
			return nil
		}
		if access == CreatorAccess && c.hasCreator(academic.ID()) {
			return nil
		}
	}
	return academicCantEditCourseError{
		requestingAcademicID: academic.ID(),
		access:               access,
		courseID:             c.ID(),
	}
}

func (a Academic) CanCreateCourse() error {
	if a.Type() == Teacher {
		return nil
	}
	return ErrNotTeacherCantCreateCourse
}
