package command

type academicsService interface {
	TeacherExists(teacherID string) error

	StudentExists(studentID string) error

	GroupExists(groupID string) error
}
