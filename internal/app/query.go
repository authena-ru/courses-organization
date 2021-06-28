package app

import "github.com/authena-ru/courses-organization/internal/domain/course"

type AllTasksQuery struct {
	Academic course.Academic
	CourseID string
	Type     course.TaskType
	Text     string
}

type SpecificTaskQuery struct {
	Academic   course.Academic
	CourseID   string
	TaskNumber int
}
