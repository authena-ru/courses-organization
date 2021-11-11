package app

import "context"

type Application struct {
	Commands Commands
	Queries  Queries
}

type (
	Commands struct {
		CreateCourse       createCourseHandler
		ExtendCourse       extendCourseHandler
		AddCollaborator    addCollaboratorHandler
		RemoveCollaborator removeCollaboratorHandler
		AddStudent         addStudentHandler
		RemoveStudent      removeStudentHandler
		AddTask            addTaskHandler
	}

	createCourseHandler interface {
		// Handle is CreateCourseCommand handler.
		// Creates course, returns ID of new brand course and one of possible errors:
		// app.ErrDatabaseProblems, course.ErrNotTeacherCantCreateCourse, error that can
		// be detected using methods course.IsInvalidTaskParametersError and others without definition.
		Handle(ctx context.Context, cmd CreateCourseCommand) (string, error)
	}

	extendCourseHandler interface {
		// Handle is ExtendCourseCommand handler.
		// Extends origin course, returns extended course ID and one of possible errors:
		// app.ErrCourseDoesntExist, app.ErrDatabaseProblems, course.ErrNotTeacherCantCreateCourse,
		// errors that can be detected using method course.IsInvalidTaskParametersError,
		// course.IsAcademicCantEditCourseError, and others without definition.
		Handle(ctx context.Context, cmd ExtendCourseCommand) (string, error)
	}

	addCollaboratorHandler interface {
		// Handle is AddCollaboratorCommand handler.
		// Adds one collaborator to course, returns one of possible errors:
		// app.ErrTeacherDoesntExist, app.ErrCourseDoesntExist, app.ErrDatabaseProblems,
		// error that can be detected using course.IsAcademicCantEditCourseError and
		// others without definition.
		Handle(ctx context.Context, cmd AddCollaboratorCommand) error
	}

	removeCollaboratorHandler interface {
		// Handle is RemoveCollaboratorCommand handler.
		// Removes one collaborator from course, returns one of possible errors:
		// app.ErrCourseDoesntExist, app.ErrDatabaseProblems, course.ErrCourseHasNoSuchCollaborator
		// error that can be detected using method course.IsAcademicCantEditCourseError and others without definition.
		Handle(ctx context.Context, cmd RemoveCollaboratorCommand) error
	}

	addStudentHandler interface {
		// Handle is AddStudentCommand handler.
		// Adds one student to course, returns one of possible errors:
		// app.ErrStudentDoesntExist, app.ErrCourseDoesntExist, app.ErrDatabaseProblems,
		// error that can be detected using method course.IsAcademicCantEditCourseError and
		// others without definition.
		Handle(ctx context.Context, cmd AddStudentCommand) error
	}

	removeStudentHandler interface {
		// Handle is RemoveStudentCommand handler.
		// Removes one student from course, returns one of possible errors:
		// app.ErrCourseDoesntExist, app.ErrDatabaseProblems, course.ErrCourseHasNoSuchStudent
		// error that can be detected using method course.IsAcademicCantEditCourseError and others without definition.
		Handle(ctx context.Context, cmd RemoveStudentCommand) error
	}

	addTaskHandler interface {
		// Handle is AddTaskCommand handler.
		// Adds task with manual checking, auto code checking or testing type,
		// returns one of possible errors. app.ErrCourseDoesntExist, errors that can
		// be detected using methods course.IsInvalidTaskParametersError,
		// course.IsAcademicCantEditCourseError and others without definition.
		Handle(ctx context.Context, cmd AddTaskCommand) (int, error)
	}
)

type (
	Queries struct {
		SpecificCourse specificCourseHandler
		SpecificTask   specificTaskHandler
		AllTasks       allTasksHandler
	}

	specificCourseHandler interface {
		// Handle is SpecificCourseQuery handler.
		// Returns course with given ID.
		// If course doesn't exist, an error equal app.ErrCourseDoesntExist.
		Handle(ctx context.Context, qry SpecificCourseQuery) (CommonCourse, error)
	}

	specificTaskHandler interface {
		// Handle is SpecificTaskQuery handler.
		// Returns course task with given number.
		// If course doesn't exist, an error equal app.ErrCourseDoesntExist.
		// If task doesn't exist, an error equal app.ErrTaskDoesntExist.
		Handle(ctx context.Context, qry SpecificTaskQuery) (SpecificTask, error)
	}

	allTasksHandler interface {
		// Handle is AllTasksQuery handler.
		// Returns list of course tasks with general task parameters.
		// Tasks filtered by type, title and description.
		// If course doesn't exist, error equal app.ErrCourseDoesntExist.
		Handle(ctx context.Context, qry AllTasksQuery) ([]GeneralTask, error)
	}
)
