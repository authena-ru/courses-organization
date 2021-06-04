package course_test

import (
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/authena-ru/courses-organization/internal/coursesorg/domain/course"
)

func TestCourse_AddManualCheckingTask(t *testing.T) {
	t.Parallel()
	creator := course.MustNewAcademic("creator-id", course.Teacher)
	crs := course.MustNewCourse(course.CreationParams{
		ID:      "course-id",
		Creator: creator,
		Title:   "Docker",
		Period:  course.MustNewPeriod(2021, 2022, course.SecondSemester),
		Started: true,
	})
	correctTaskCreationParams := course.ManualCheckingTaskCreationParams{
		Title:       "Make container",
		Description: "Containerization",
		Deadline: course.MustNewDeadline(
			time.Date(2021, time.September, 9, 0, 0, 0, 0, time.Local),
			time.Date(2021, time.September, 21, 0, 0, 0, 0, time.Local),
		),
	}
	testCases := []struct {
		Name     string
		Course   course.Course
		Academic course.Academic
		Params   course.ManualCheckingTaskCreationParams
		IsErr    func(err error) bool
	}{
		{
			Name:     "add_task_to_course_and_obtain_number",
			Course:   *crs,
			Academic: creator,
			Params:   correctTaskCreationParams,
		},
		{
			Name:     "academic_cant_add_task",
			Course:   *crs,
			Academic: course.MustNewAcademic("not-course-teacher-id", course.Teacher),
			Params:   correctTaskCreationParams,
			IsErr:    course.IsAcademicCantEditCourseError,
		},
		{
			Name:     "task_title_too_long",
			Course:   *crs,
			Academic: creator,
			Params: course.ManualCheckingTaskCreationParams{
				Title: strings.Repeat("x", 201),
			},
			IsErr: func(err error) bool {
				return errors.Is(err, course.ErrTaskTitleTooLong)
			},
		},
		{
			Name:     "task_description_too_long",
			Course:   *crs,
			Academic: creator,
			Params: course.ManualCheckingTaskCreationParams{
				Description: strings.Repeat("x", 1001),
			},
			IsErr: func(err error) bool {
				return errors.Is(err, course.ErrTaskDescriptionTooLong)
			},
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			number, err := c.Course.AddManualCheckingTask(c.Academic, c.Params)
			if c.IsErr != nil {
				require.Error(t, err)
				require.True(t, c.IsErr(err))
				require.Equal(t, 0, number)
				return
			}
			require.Equal(t, 1, crs.TasksNumber())
			task, err := c.Course.Task(number)
			require.NoError(t, err)
			require.Equal(t, number, task.Number())
			require.Equal(t, course.ManualChecking, task.Type())
			require.Equal(t, c.Params.Title, task.Title())
			require.Equal(t, c.Params.Description, task.Description())
			require.Equal(t, c.Params.Deadline, task.ManualCheckingOptional())
		})
	}
}

func TestCourse_AddAutoCodeCheckingTask(t *testing.T) {
	t.Parallel()
	creator := course.MustNewAcademic("creator-id", course.Teacher)
	collaborator := course.MustNewAcademic("collaborator-id", course.Teacher)
	crs := course.MustNewCourse(course.CreationParams{
		ID:            "course-id",
		Creator:       creator,
		Title:         "Python Django",
		Period:        course.MustNewPeriod(2023, 2024, course.FirstSemester),
		Started:       true,
		Collaborators: []string{collaborator.ID()},
	})
	correctTaskCreationParams := course.AutoCodeCheckingTaskCreationParams{
		Title:       "Print sum of two integers",
		Description: "You should read two integers from console and print sum",
		Deadline: course.MustNewDeadline(
			time.Date(2023, time.November, 10, 0, 0, 0, 0, time.Local),
			time.Date(2023, time.November, 20, 0, 0, 0, 0, time.Local),
		),
		TestData: []course.TestData{course.MustNewTestData("1 3", "4"), course.MustNewTestData("2 2", "4")},
	}
	testCases := []struct {
		Name     string
		Course   course.Course
		Academic course.Academic
		Params   course.AutoCodeCheckingTaskCreationParams
		IsErr    func(err error) bool
	}{
		{
			Name:     "add_task_to_course_and_obtain_number",
			Course:   *crs,
			Academic: collaborator,
			Params:   correctTaskCreationParams,
		},
		{
			Name:     "academic_cant_add_task",
			Course:   *crs,
			Academic: course.MustNewAcademic("student-id", course.Student),
			Params:   correctTaskCreationParams,
			IsErr:    course.IsAcademicCantEditCourseError,
		},
		{
			Name:     "task_title_too_long",
			Course:   *crs,
			Academic: creator,
			Params: course.AutoCodeCheckingTaskCreationParams{
				Title: strings.Repeat("x", 201),
			},
			IsErr: func(err error) bool {
				return errors.Is(err, course.ErrTaskTitleTooLong)
			},
		},
		{
			Name:     "task_description_too_long",
			Course:   *crs,
			Academic: collaborator,
			Params: course.AutoCodeCheckingTaskCreationParams{
				Description: strings.Repeat("x", 1001),
			},
			IsErr: func(err error) bool {
				return errors.Is(err, course.ErrTaskDescriptionTooLong)
			},
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			number, err := c.Course.AddAutoCodeCheckingTask(c.Academic, c.Params)
			if c.IsErr != nil {
				require.Error(t, err)
				require.True(t, c.IsErr(err))
				require.Equal(t, 0, number)
				return
			}
			require.Equal(t, 1, crs.TasksNumber())
			task, err := c.Course.Task(number)
			require.NoError(t, err)
			require.Equal(t, number, task.Number())
			require.Equal(t, course.AutoCodeChecking, task.Type())
			require.Equal(t, c.Params.Title, task.Title())
			require.Equal(t, c.Params.Description, task.Description())
			deadline, testData := task.AutoCodeCheckingOptional()
			require.Equal(t, c.Params.Deadline, deadline)
			require.Equal(t, c.Params.TestData, testData)
		})
	}
}

func TestCourse_AddTestingTask(t *testing.T) {
	t.Parallel()
	creator := course.MustNewAcademic("creator-id", course.Teacher)
	crs := course.MustNewCourse(course.CreationParams{
		ID:      "course-id",
		Creator: creator,
		Title:   "Golang channels",
		Period:  course.MustNewPeriod(2021, 2022, course.FirstSemester),
		Started: true,
	})
	correctTaskCreationParams := course.TestingTaskCreationParams{
		Title:       "Golang syntax",
		Description: "Choose right syntactic constructions",
		TestPoints:  []course.TestPoint{course.MustNewTestPoint("How to create channel in Go", []string{"make(chan int)", "chan int {}"}, []int{0})},
	}
	testCases := []struct {
		Name     string
		Course   course.Course
		Academic course.Academic
		Params   course.TestingTaskCreationParams
		IsErr    func(err error) bool
	}{
		{
			Name:     "add_task_to_course_and_obtain_number",
			Course:   *crs,
			Academic: creator,
			Params:   correctTaskCreationParams,
		},
		{
			Name:     "academic_cant_add_task",
			Course:   *crs,
			Academic: course.MustNewAcademic("other-teacher-id", course.Teacher),
			Params:   correctTaskCreationParams,
			IsErr:    course.IsAcademicCantEditCourseError,
		},
		{
			Name:     "task_title_too_long",
			Course:   *crs,
			Academic: creator,
			Params: course.TestingTaskCreationParams{
				Title: strings.Repeat("x", 201),
			},
			IsErr: func(err error) bool {
				return errors.Is(err, course.ErrTaskTitleTooLong)
			},
		},
		{
			Name:     "task_description_too_long",
			Course:   *crs,
			Academic: creator,
			Params: course.TestingTaskCreationParams{
				Description: strings.Repeat("x", 1001),
			},
			IsErr: func(err error) bool {
				return errors.Is(err, course.ErrTaskDescriptionTooLong)
			},
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			number, err := c.Course.AddTestingTask(c.Academic, c.Params)
			if c.IsErr != nil {
				require.Error(t, err)
				require.True(t, c.IsErr(err))
				require.Equal(t, 0, number)
				return
			}
			require.Equal(t, 1, c.Course.TasksNumber())
			task, err := c.Course.Task(number)
			require.NoError(t, err)
			require.Equal(t, number, task.Number())
			require.Equal(t, course.Testing, task.Type())
			require.Equal(t, c.Params.Title, task.Title())
			require.Equal(t, c.Params.Description, task.Description())
			require.Equal(t, c.Params.TestPoints, task.TestingOptional())
		})
	}
}

func TestCourse_RenameTask(t *testing.T) {
	t.Parallel()
	creator := course.MustNewAcademic("creator-id", course.Teacher)
	crs := course.MustNewCourse(course.CreationParams{
		ID:      "course-id",
		Creator: creator,
		Title:   "Learn TypeScript",
		Period:  course.MustNewPeriod(2021, 2022, course.FirstSemester),
	})
	taskNumber, err := crs.AddManualCheckingTask(creator, course.ManualCheckingTaskCreationParams{
		Title: "Classes in TypeScript",
	})
	require.NoError(t, err)
	testCases := []struct {
		Name       string
		Course     course.Course
		TaskNumber int
		Academic   course.Academic
		NewTitle   string
		IsErr      func(err error) bool
	}{
		{
			Name:       "rename_task_to_new_valid_title",
			Course:     *crs,
			TaskNumber: taskNumber,
			Academic:   creator,
			NewTitle:   "Classez in Typescript",
		},
		{
			Name:       "academic_cant_rename_task",
			Course:     *crs,
			TaskNumber: taskNumber,
			Academic:   course.MustNewAcademic("student-id", course.Student),
			NewTitle:   "TypeScript classes",
			IsErr:      course.IsAcademicCantEditCourseError,
		},
		{
			Name:       "task_title_too_loong",
			Course:     *crs,
			TaskNumber: taskNumber,
			Academic:   creator,
			NewTitle:   strings.Repeat("x", 201),
			IsErr: func(err error) bool {
				return errors.Is(err, course.ErrTaskTitleTooLong)
			},
		},
		{
			Name:       "no_task_with_number",
			Course:     *crs,
			TaskNumber: crs.TasksNumber() + 1,
			Academic:   creator,
			NewTitle:   "Classes",
			IsErr: func(err error) bool {
				return errors.Is(err, course.ErrNoTaskWithNumber)
			},
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			err := c.Course.RenameTask(c.Academic, c.TaskNumber, c.NewTitle)

			if c.IsErr != nil {
				require.Error(t, err)
				require.True(t, c.IsErr(err))
				return
			}
			require.NoError(t, err)
			task, err := c.Course.Task(c.TaskNumber)
			require.NoError(t, err)
			require.Equal(t, c.NewTitle, task.Title())
		})
	}
}

func TestCourse_ReplaceTaskDescription(t *testing.T) {
	t.Parallel()
	creator := course.MustNewAcademic("creator-id", course.Teacher)
	crs := course.MustNewCourse(course.CreationParams{
		ID:      "course-id",
		Creator: creator,
		Title:   "C#",
		Period:  course.MustNewPeriod(2021, 2022, course.SecondSemester),
	})
	taskNumber, err := crs.AddManualCheckingTask(creator, course.ManualCheckingTaskCreationParams{
		Description: "Write your binary search",
	})
	require.NoError(t, err)
	testCases := []struct {
		Name           string
		Course         course.Course
		TaskNumber     int
		Academic       course.Academic
		NewDescription string
		IsErr          func(err error) bool
	}{
		{
			Name:           "replace_task_description_with_new_valid_description",
			Course:         *crs,
			TaskNumber:     taskNumber,
			Academic:       creator,
			NewDescription: "Write your search",
		},
		{
			Name:           "academic_cant_replace_description",
			Course:         *crs,
			TaskNumber:     taskNumber,
			Academic:       course.MustNewAcademic("student-id", course.Student),
			NewDescription: "Rewrite search",
			IsErr:          course.IsAcademicCantEditCourseError,
		},
		{
			Name:           "task_description_too_long",
			Course:         *crs,
			TaskNumber:     taskNumber,
			Academic:       creator,
			NewDescription: strings.Repeat("x", 1001),
			IsErr: func(err error) bool {
				return errors.Is(err, course.ErrTaskDescriptionTooLong)
			},
		},
		{
			Name:           "no_task_with_number",
			Course:         *crs,
			TaskNumber:     crs.TasksNumber() + 1,
			Academic:       creator,
			NewDescription: "Write search algorithm",
			IsErr: func(err error) bool {
				return errors.Is(err, course.ErrNoTaskWithNumber)
			},
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			err := c.Course.ReplaceTaskDescription(c.Academic, c.TaskNumber, c.NewDescription)

			if c.IsErr != nil {
				require.Error(t, err)
				require.True(t, c.IsErr(err))
				return
			}
			require.NoError(t, err)
			task, err := c.Course.Task(c.TaskNumber)
			require.NoError(t, err)
			require.Equal(t, c.NewDescription, task.Description())
		})
	}
}
