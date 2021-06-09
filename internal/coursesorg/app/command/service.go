package command

type academicsService interface {
	// TeacherExists should return app.ErrTeacherDoesntExist
	// when academics service can't find teacher with such id.
	TeacherExists(teacherID string) error

	// StudentExists should return app.ErrStudentDoesntExist
	// when academics service can't find student with such id.
	StudentExists(studentID string) error

	// GroupExists should return app.ErrGroupDoesntExist
	// when academics service can't find group with such id.
	GroupExists(groupID string) error
}
