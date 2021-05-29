package course

func (c *Course) Students() []string {
	students := make([]string, 0, len(c.students))
	for s := range c.students {
		students = append(students, s)
	}
	return students
}

func (c *Course) AddStudents(academic Academic, studentIDs ...string) error {
	if err := c.CanAcademicEditWithAccess(academic, TeacherAccess); err != nil {
		return err
	}
	for _, sid := range studentIDs {
		c.students[sid] = true
	}
	return nil
}

func (c *Course) RemoveStudents(academic Academic, studentIDs ...string) error {
	if err := c.CanAcademicEditWithAccess(academic, TeacherAccess); err != nil {
		return err
	}
	for _, sid := range studentIDs {
		delete(c.students, sid)
	}
	return nil
}
