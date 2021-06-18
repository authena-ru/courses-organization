package app

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	ErrCourseDoesntExist     = errors.New("course doesn't exist")
	ErrTeacherDoesntExist    = errors.New("teacher doesn't exist")
	ErrStudentDoesntExist    = errors.New("student doesn't exist")
	ErrGroupDoesntExist      = errors.New("group doesn't exist")
	ErrCourseTaskDoesntExist = errors.New("course task doesn't exist")
	ErrDatabaseProblems      = errors.New("database problems")
)

type errorWrapper struct {
	appErr    error
	originErr error
}

func Wrap(applicationError error, originError error) error {
	if originError == nil {
		return nil
	}
	if applicationError == nil {
		return originError
	}
	return errors.WithStack(&errorWrapper{
		appErr:    applicationError,
		originErr: originError,
	})
}

func (e errorWrapper) Error() string {
	return fmt.Sprintf("%s: %s", e.appErr, e.originErr)
}

func (e errorWrapper) Cause() error {
	return e.appErr
}

func (e errorWrapper) Unwrap() error {
	return e.appErr
}
