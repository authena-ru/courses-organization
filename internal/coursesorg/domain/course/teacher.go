package course

import "github.com/pkg/errors"

func (c *Course) Collaborators() []string {
	collaborators := make([]string, 0, len(c.collaborators))
	for c := range c.collaborators {
		collaborators = append(collaborators, c)
	}
	return collaborators
}

func (c *Course) AddCollaborators(academic Academic, teacherIDs ...string) error {
	if err := c.canAcademicEditWithAccess(academic, TeacherAccess); err != nil {
		return err
	}
	for _, tid := range teacherIDs {
		c.collaborators[tid] = true
	}
	return nil
}

var ErrCourseHasNoSuchCollaborator = errors.New("course has no such collaborator")

func (c *Course) RemoveCollaborator(academic Academic, teacherID string) error {
	if err := c.canAcademicEditWithAccess(academic, CreatorAccess); err != nil {
		return err
	}
	if !c.hasTeacher(teacherID) {
		return ErrCourseHasNoSuchCollaborator
	}
	delete(c.collaborators, teacherID)
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
