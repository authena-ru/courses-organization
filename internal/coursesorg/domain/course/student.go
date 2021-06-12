package course

import "github.com/pkg/errors"

func (c *Course) Students() []string {
	students := make([]string, 0, len(c.students))
	for s := range c.students {
		students = append(students, s)
	}
	return students
}

func (c *Course) AddStudents(academic Academic, studentIDs ...string) error {
	if err := c.canAcademicEditWithAccess(academic, TeacherAccess); err != nil {
		return err
	}
	c.putStudents(studentIDs)
	return nil
}

func (c *Course) putStudents(studentIDs []string) {
	for _, sid := range studentIDs {
		c.students[sid] = true
	}
}

var ErrCourseHasNoSuchStudent = errors.New("course has no such student")

func (c *Course) RemoveStudent(academic Academic, studentID string) error {
	if err := c.canAcademicEditWithAccess(academic, TeacherAccess); err != nil {
		return err
	}
	if !c.hasStudent(studentID) {
		return ErrCourseHasNoSuchStudent
	}
	delete(c.students, studentID)
	return nil
}
