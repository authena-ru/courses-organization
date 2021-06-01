package course

func (c *Course) Collaborators() []string {
	collaborators := make([]string, 0, len(c.collaborators))
	for c := range c.collaborators {
		collaborators = append(collaborators, c)
	}
	return collaborators
}

func (c *Course) AddCollaborators(academic Academic, teacherIDs ...string) error {
	if err := c.CanAcademicEditWithAccess(academic, TeacherAccess); err != nil {
		return err
	}
	for _, tid := range teacherIDs {
		c.collaborators[tid] = true
	}
	return nil
}

func (c *Course) RemoveCollaborators(academic Academic, teacherIDs ...string) error {
	if err := c.CanAcademicEditWithAccess(academic, CreatorAccess); err != nil {
		return err
	}
	for _, tid := range teacherIDs {
		delete(c.collaborators, tid)
	}
	return nil
}

func (c *Course) hasTeacher(teacherID string) bool {
	return c.hasCreator(teacherID) || c.collaborators[teacherID]
}

func (c *Course) hasCreator(teacherID string) bool {
	return c.creatorID == teacherID
}

func (c *Course) hasStudent(studentID string) bool {
	return c.students[studentID]
}
