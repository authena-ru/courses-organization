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
	ErrEmptyAcademicID     = errors.New("empty course id")
	ErrInvalidAcademicType = errors.New("invalid academic type")
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

func CanAcademicSeeCourse(academic Academic, course Course) error {
	if academic.Type() == Teacher && course.hasTeacher(academic.ID()) {
		return nil
	}
	if academic.Type() == Student && course.hasStudent(academic.ID()) {
		return nil
	}
	return academicCantSeeCourseError{requestingAcademicID: academic.ID(), courseID: course.ID()}
}

type Access uint8

const (
	TeacherAccess Access = iota + 1
	CreatorAccess
)

func (a Access) String() string {
	return [...]string{"teacher access", "creator access"}[a-1]
}

type academicCantEditCourseError struct {
	requestingAcademicID     string
	requestingAcademicAccess Access
	courseID                 string
}

func (e academicCantEditCourseError) Error() string {
	return fmt.Sprintf(
		"academic $%s with %s can't edit course #%s",
		e.requestingAcademicID, e.requestingAcademicAccess, e.courseID,
	)
}

func IsAcademicCantEditCourseError(err error) bool {
	var e academicCantEditCourseError
	return errors.As(err, &e)
}

func CanAcademicEditCourseWithAccess(academic Academic, course Course, access Access) error {
	if academic.Type() == Teacher {
		if access == TeacherAccess && course.hasTeacher(academic.ID()) {
			return nil
		}
		if access == CreatorAccess && course.hasCreator(academic.ID()) {
			return nil
		}
	}
	return academicCantEditCourseError{
		requestingAcademicID:     academic.ID(),
		requestingAcademicAccess: access,
		courseID:                 course.ID(),
	}
}
