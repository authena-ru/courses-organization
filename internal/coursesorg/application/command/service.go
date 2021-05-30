package command

import "github.com/authena-ru/courses-organization/internal/coursesorg/domain/course"

type academicsService interface {
	AcademicExists(academic course.Academic) error
}
