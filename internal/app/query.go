package app

import "github.com/authena-ru/courses-organization/internal/domain/course"

type (
	AllCoursesQuery struct {
		Academic course.Academic
		Title    string
	}

	SpecificCourseQuery struct {
		Academic course.Academic
		CourseID string
	}

	AllTasksQuery struct {
		Academic course.Academic
		CourseID string
		Type     course.TaskType
		Text     string
	}

	SpecificTaskQuery struct {
		Academic   course.Academic
		CourseID   string
		TaskNumber int
	}
)
