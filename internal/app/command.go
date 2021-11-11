package app

import "github.com/authena-ru/courses-organization/internal/domain/course"

type (
	AddCollaboratorCommand struct {
		Academic       course.Academic
		CourseID       string
		CollaboratorID string
	}

	AddStudentCommand struct {
		Academic  course.Academic
		CourseID  string
		StudentID string
	}

	AddTaskCommand struct {
		Academic        course.Academic
		CourseID        string
		TaskTitle       string
		TaskDescription string
		TaskType        course.TaskType
		Deadline        course.Deadline
		TestPoints      []course.TestPoint
		TestData        []course.TestData
	}

	CreateCourseCommand struct {
		Academic      course.Academic
		CourseStarted bool
		CourseTitle   string
		CoursePeriod  course.Period
	}

	ExtendCourseCommand struct {
		Academic       course.Academic
		OriginCourseID string
		CourseStarted  bool
		CourseTitle    string
		CoursePeriod   course.Period
	}

	RemoveCollaboratorCommand struct {
		Academic       course.Academic
		CourseID       string
		CollaboratorID string
	}

	RemoveStudentCommand struct {
		Academic  course.Academic
		CourseID  string
		StudentID string
	}
)
