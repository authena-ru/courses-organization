package app

import "github.com/authena-ru/courses-organization/internal/domain/course"

type AddCollaboratorCommand struct {
	Academic       course.Academic
	CourseID       string
	CollaboratorID string
}

type AddStudentCommand struct {
	Academic  course.Academic
	CourseID  string
	StudentID string
}

type AddTaskCommand struct {
	Academic        course.Academic
	CourseID        string
	TaskTitle       string
	TaskDescription string
	TaskType        course.TaskType
	Deadline        course.Deadline
	TestPoints      []course.TestPoint
	TestData        []course.TestData
}

type CreateCourseCommand struct {
	Academic      course.Academic
	CourseStarted bool
	CourseTitle   string
	CoursePeriod  course.Period
}

type ExtendCourseCommand struct {
	Academic       course.Academic
	OriginCourseID string
	CourseStarted  bool
	CourseTitle    string
	CoursePeriod   course.Period
}

type RemoveCollaboratorCommand struct {
	Academic       course.Academic
	CourseID       string
	CollaboratorID string
}

type RemoveStudentCommand struct {
	Academic  course.Academic
	CourseID  string
	StudentID string
}
