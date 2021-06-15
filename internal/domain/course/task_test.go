package course_test

import (
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/authena-ru/courses-organization/internal/domain/course"
)

func TestCourse_AddManualCheckingTask(t *testing.T) {
	t.Parallel()
	creator := course.MustNewAcademic("creator-id", course.TeacherType)
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
		Academic course.Academic
		Params   course.ManualCheckingTaskCreationParams
		IsErr    func(err error) bool
	}{
		{
			Name:     "add_task_to_course_and_obtain_number",
			Academic: creator,
			Params:   correctTaskCreationParams,
		},
		{
			Name:     "academic_cant_add_task",
			Academic: course.MustNewAcademic("not-course-teacher-id", course.TeacherType),
			Params:   correctTaskCreationParams,
			IsErr:    course.IsAcademicCantEditCourseError,
		},
		{
			Name:     "task_title_too_long",
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

			crs := course.MustNewCourse(course.CreationParams{
				ID:      "course-id",
				Creator: creator,
				Title:   "Docker",
				Period:  course.MustNewPeriod(2021, 2022, course.SecondSemester),
				Started: true,
			})

			number, err := crs.AddManualCheckingTask(c.Academic, c.Params)
			if c.IsErr != nil {
				require.Error(t, err)
				require.True(t, c.IsErr(err))
				require.Equal(t, 0, number)
				return
			}
			require.Equal(t, 1, crs.TasksNumber())
			task, err := crs.Task(number)
			require.NoError(t, err)
			require.Equal(t, number, task.Number())
			require.Equal(t, course.ManualCheckingType, task.Type())
			require.Equal(t, c.Params.Title, task.Title())
			require.Equal(t, c.Params.Description, task.Description())
			require.Equal(t, c.Params.Deadline, task.ManualCheckingOptional())
		})
	}
}

func TestCourse_AddAutoCodeCheckingTask(t *testing.T) {
	t.Parallel()
	creator := course.MustNewAcademic("creator-id", course.TeacherType)
	collaborator := course.MustNewAcademic("collaborator-id", course.TeacherType)
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
		Academic course.Academic
		Params   course.AutoCodeCheckingTaskCreationParams
		IsErr    func(err error) bool
	}{
		{
			Name:     "add_task_to_course_and_obtain_number",
			Academic: collaborator,
			Params:   correctTaskCreationParams,
		},
		{
			Name:     "academic_cant_add_task",
			Academic: course.MustNewAcademic("student-id", course.StudentType),
			Params:   correctTaskCreationParams,
			IsErr:    course.IsAcademicCantEditCourseError,
		},
		{
			Name:     "task_title_too_long",
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

			crs := course.MustNewCourse(course.CreationParams{
				ID:            "course-id",
				Creator:       creator,
				Title:         "Python Django",
				Period:        course.MustNewPeriod(2023, 2024, course.FirstSemester),
				Started:       true,
				Collaborators: []string{collaborator.ID()},
			})

			number, err := crs.AddAutoCodeCheckingTask(c.Academic, c.Params)
			if c.IsErr != nil {
				require.Error(t, err)
				require.True(t, c.IsErr(err))
				require.Equal(t, 0, number)
				return
			}
			require.Equal(t, 1, crs.TasksNumber())
			task, err := crs.Task(number)
			require.NoError(t, err)
			require.Equal(t, number, task.Number())
			require.Equal(t, course.AutoCodeCheckingType, task.Type())
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
	creator := course.MustNewAcademic("creator-id", course.TeacherType)
	correctTaskCreationParams := course.TestingTaskCreationParams{
		Title:       "Golang syntax",
		Description: "Choose right syntactic constructions",
		TestPoints:  []course.TestPoint{course.MustNewTestPoint("How to create channel in Go", []string{"make(chan int)", "chan int {}"}, []int{0})},
	}
	testCases := []struct {
		Name     string
		Academic course.Academic
		Params   course.TestingTaskCreationParams
		IsErr    func(err error) bool
	}{
		{
			Name:     "add_task_to_course_and_obtain_number",
			Academic: creator,
			Params:   correctTaskCreationParams,
		},
		{
			Name:     "academic_cant_add_task",
			Academic: course.MustNewAcademic("other-teacher-id", course.TeacherType),
			Params:   correctTaskCreationParams,
			IsErr:    course.IsAcademicCantEditCourseError,
		},
		{
			Name:     "task_title_too_long",
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

			crs := course.MustNewCourse(course.CreationParams{
				ID:      "course-id",
				Creator: creator,
				Title:   "Golang channels",
				Period:  course.MustNewPeriod(2021, 2022, course.FirstSemester),
				Started: true,
			})

			number, err := crs.AddTestingTask(c.Academic, c.Params)
			if c.IsErr != nil {
				require.Error(t, err)
				require.True(t, c.IsErr(err))
				require.Equal(t, 0, number)
				return
			}
			require.Equal(t, 1, crs.TasksNumber())
			task, err := crs.Task(number)
			require.NoError(t, err)
			require.Equal(t, number, task.Number())
			require.Equal(t, course.TestingType, task.Type())
			require.Equal(t, c.Params.Title, task.Title())
			require.Equal(t, c.Params.Description, task.Description())
			require.Equal(t, c.Params.TestPoints, task.TestingOptional())
		})
	}
}

func TestCourse_RenameTask(t *testing.T) {
	t.Parallel()
	creator := course.MustNewAcademic("creator-id", course.TeacherType)
	addTask := func(crs *course.Course) int {
		taskNumber, err := crs.AddManualCheckingTask(creator, course.ManualCheckingTaskCreationParams{
			Title: "Classes in TypeScript",
		})
		require.NoError(t, err)
		return taskNumber
	}
	testCases := []struct {
		Name        string
		Academic    course.Academic
		NewTitle    string
		PrepareTask func(crs *course.Course) int
		IsErr       func(err error) bool
	}{
		{
			Name:        "rename_task_to_new_valid_title",
			Academic:    creator,
			NewTitle:    "Classez in Typescript",
			PrepareTask: addTask,
		},
		{
			Name:        "academic_cant_rename_task",
			Academic:    course.MustNewAcademic("student-id", course.StudentType),
			NewTitle:    "TypeScript classes",
			PrepareTask: addTask,
			IsErr:       course.IsAcademicCantEditCourseError,
		},
		{
			Name:        "task_title_too_loong",
			Academic:    creator,
			NewTitle:    strings.Repeat("x", 201),
			PrepareTask: addTask,
			IsErr: func(err error) bool {
				return errors.Is(err, course.ErrTaskTitleTooLong)
			},
		},
		{
			Name:     "no_task_with_number",
			Academic: creator,
			NewTitle: "Classes",
			PrepareTask: func(crs *course.Course) int {
				return crs.TasksNumber() + 1
			},
			IsErr: func(err error) bool {
				return errors.Is(err, course.ErrCourseHasNoSuchTask)
			},
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			crs := course.MustNewCourse(course.CreationParams{
				ID:      "course-id",
				Creator: creator,
				Title:   "Learn TypeScript",
				Period:  course.MustNewPeriod(2021, 2022, course.FirstSemester),
			})
			taskNumber := c.PrepareTask(crs)

			err := crs.RenameTask(c.Academic, taskNumber, c.NewTitle)

			if c.IsErr != nil {
				require.Error(t, err)
				require.True(t, c.IsErr(err))
				return
			}
			require.NoError(t, err)
			task, err := crs.Task(taskNumber)
			require.NoError(t, err)
			require.Equal(t, c.NewTitle, task.Title())
		})
	}
}

func TestCourse_ReplaceTaskDescription(t *testing.T) {
	t.Parallel()
	creator := course.MustNewAcademic("creator-id", course.TeacherType)
	addTask := func(crs *course.Course) int {
		taskNumber, err := crs.AddManualCheckingTask(creator, course.ManualCheckingTaskCreationParams{
			Description: "Write your binary search",
		})
		require.NoError(t, err)
		return taskNumber
	}
	testCases := []struct {
		Name           string
		Academic       course.Academic
		NewDescription string
		PrepareTask    func(crs *course.Course) int
		IsErr          func(err error) bool
	}{
		{
			Name:           "replace_task_description_with_new_valid_description",
			Academic:       creator,
			NewDescription: "Write your search",
			PrepareTask:    addTask,
		},
		{
			Name:           "academic_cant_replace_description",
			Academic:       course.MustNewAcademic("student-id", course.StudentType),
			NewDescription: "Rewrite search",
			PrepareTask:    addTask,
			IsErr:          course.IsAcademicCantEditCourseError,
		},
		{
			Name:           "task_description_too_long",
			Academic:       creator,
			NewDescription: strings.Repeat("x", 1001),
			PrepareTask:    addTask,
			IsErr: func(err error) bool {
				return errors.Is(err, course.ErrTaskDescriptionTooLong)
			},
		},
		{
			Name:           "no_task_with_number",
			Academic:       creator,
			NewDescription: "Write search algorithm",
			PrepareTask: func(crs *course.Course) int {
				return crs.TasksNumber() + 1
			},
			IsErr: func(err error) bool {
				return errors.Is(err, course.ErrCourseHasNoSuchTask)
			},
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			crs := course.MustNewCourse(course.CreationParams{
				ID:      "course-id",
				Creator: creator,
				Title:   "C#",
				Period:  course.MustNewPeriod(2021, 2022, course.SecondSemester),
			})
			taskNumber := c.PrepareTask(crs)

			err := crs.ReplaceTaskDescription(c.Academic, taskNumber, c.NewDescription)

			if c.IsErr != nil {
				require.Error(t, err)
				require.True(t, c.IsErr(err))
				return
			}
			require.NoError(t, err)
			task, err := crs.Task(taskNumber)
			require.NoError(t, err)
			require.Equal(t, c.NewDescription, task.Description())
		})
	}
}

func TestCourse_ReplaceTaskDeadline(t *testing.T) {
	t.Parallel()
	creator := course.MustNewAcademic("creator-id", course.TeacherType)
	newDeadline := course.MustNewDeadline(
		time.Date(2023, time.March, 10, 0, 0, 0, 0, time.Local),
		time.Date(2023, time.March, 22, 0, 0, 0, 0, time.Local),
	)
	addTasks := func(crs *course.Course) (int, int) {
		manualTaskNumber, err := crs.AddManualCheckingTask(creator, course.ManualCheckingTaskCreationParams{
			Deadline: course.MustNewDeadline(
				time.Date(2023, time.March, 1, 0, 0, 0, 9, time.Local),
				time.Date(2023, time.March, 12, 0, 0, 0, 0, time.Local),
			),
		})
		require.NoError(t, err)
		testingTaskNumber, err := crs.AddTestingTask(creator, course.TestingTaskCreationParams{})
		require.NoError(t, err)
		return manualTaskNumber, testingTaskNumber
	}
	testCases := []struct {
		Name        string
		Academic    course.Academic
		NewDeadline course.Deadline
		PrepareTask func(crs *course.Course) int
		IsErr       func(err error) bool
	}{
		{
			Name:        "replace_task_deadline_with_new_valid_deadline",
			Academic:    creator,
			NewDeadline: newDeadline,
			PrepareTask: func(crs *course.Course) int {
				manualTaskNumber, _ := addTasks(crs)
				return manualTaskNumber
			},
		},
		{
			Name:        "task_has_no_deadline",
			Academic:    creator,
			NewDeadline: newDeadline,
			PrepareTask: func(crs *course.Course) int {
				_, testingTaskNumber := addTasks(crs)
				return testingTaskNumber
			},
			IsErr: func(err error) bool {
				return errors.Is(err, course.ErrTaskHasNoDeadline)
			},
		},
		{
			Name:        "academic_cant_replace_deadline",
			Academic:    course.MustNewAcademic("other-teacher-id", course.TeacherType),
			NewDeadline: newDeadline,
			PrepareTask: func(crs *course.Course) int {
				manualTaskNumber, _ := addTasks(crs)
				return manualTaskNumber
			},
			IsErr: course.IsAcademicCantEditCourseError,
		},
		{
			Name:        "no_task_with_number",
			Academic:    creator,
			NewDeadline: newDeadline,
			PrepareTask: func(crs *course.Course) int {
				return crs.TasksNumber() + 1
			},
			IsErr: func(err error) bool {
				return errors.Is(err, course.ErrCourseHasNoSuchTask)
			},
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			crs := course.MustNewCourse(course.CreationParams{
				ID:      "course-id",
				Creator: creator,
				Title:   "Python",
				Period:  course.MustNewPeriod(2022, 2023, course.SecondSemester),
			})
			taskNumber := c.PrepareTask(crs)

			err := crs.ReplaceTaskDeadline(c.Academic, taskNumber, c.NewDeadline)

			if c.IsErr != nil {
				require.Error(t, err)
				require.True(t, c.IsErr(err))
				return
			}
			require.NoError(t, err)
			task, err := crs.Task(taskNumber)
			require.NoError(t, err)
			switch task.Type() {
			case course.ManualCheckingType:
				require.Equal(t, c.NewDeadline, task.ManualCheckingOptional())
			case course.AutoCodeCheckingType:
				deadline, _ := task.AutoCodeCheckingOptional()
				require.Equal(t, c.NewDeadline, deadline)
			default:
				panic("unreachable")
			}
		})
	}
}

func TestCourse_ReplaceTaskTestData(t *testing.T) {
	t.Parallel()
	creator := course.MustNewAcademic("creator-id", course.TeacherType)
	addTasks := func(crs *course.Course) (int, int) {
		manualTaskNumber, err := crs.AddManualCheckingTask(creator, course.ManualCheckingTaskCreationParams{})
		require.NoError(t, err)
		autoCodeTaskNumber, err := crs.AddAutoCodeCheckingTask(creator, course.AutoCodeCheckingTaskCreationParams{
			TestData: []course.TestData{course.MustNewTestData("1 1 2 3", "7")},
		})
		require.NoError(t, err)
		return autoCodeTaskNumber, manualTaskNumber
	}
	newTestData := []course.TestData{course.MustNewTestData("1 1 2 3", "7"), course.MustNewTestData("1 1", "2")}
	testCases := []struct {
		Name        string
		Academic    course.Academic
		NewTestData []course.TestData
		PrepareTask func(crs *course.Course) int
		IsErr       func(err error) bool
	}{
		{
			Name:        "replace_task_test_data_with_new_test_data",
			Academic:    creator,
			NewTestData: newTestData,
			PrepareTask: func(crs *course.Course) int {
				autoCodeTaskNumber, _ := addTasks(crs)
				return autoCodeTaskNumber
			},
		},
		{
			Name:        "task_has_no_test_data",
			Academic:    creator,
			NewTestData: newTestData,
			PrepareTask: func(crs *course.Course) int {
				_, manualTaskNumber := addTasks(crs)
				return manualTaskNumber
			},
			IsErr: func(err error) bool {
				return errors.Is(err, course.ErrTaskHasNoTestData)
			},
		},
		{
			Name:        "academic_cant_replace_test_data",
			Academic:    course.MustNewAcademic("other-teacher-id", course.TeacherType),
			NewTestData: newTestData,
			PrepareTask: func(crs *course.Course) int {
				autoCodeTaskNumber, _ := addTasks(crs)
				return autoCodeTaskNumber
			},
			IsErr: course.IsAcademicCantEditCourseError,
		},
		{
			Name:        "no_task_with_number",
			Academic:    creator,
			NewTestData: newTestData,
			PrepareTask: func(crs *course.Course) int {
				return crs.TasksNumber() + 1
			},
			IsErr: func(err error) bool {
				return errors.Is(err, course.ErrCourseHasNoSuchTask)
			},
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			crs := course.MustNewCourse(course.CreationParams{
				ID:      "course-id",
				Creator: creator,
				Title:   "Golang",
				Period:  course.MustNewPeriod(2020, 2021, course.SecondSemester),
			})
			taskNumber := c.PrepareTask(crs)

			err := crs.ReplaceTaskTestData(c.Academic, taskNumber, c.NewTestData)

			if c.IsErr != nil {
				require.Error(t, err)
				require.True(t, c.IsErr(err))
				return
			}
			require.NoError(t, err)
			task, err := crs.Task(taskNumber)
			require.NoError(t, err)
			if task.Type() != course.AutoCodeCheckingType {
				panic("unreachable")
			}
			_, testData := task.AutoCodeCheckingOptional()
			require.Equal(t, c.NewTestData, testData)
		})
	}
}

func TestCourse_ReplaceTaskTestPoints(t *testing.T) {
	t.Parallel()
	creator := course.MustNewAcademic("creator-id", course.TeacherType)
	addTasks := func(crs *course.Course) (int, int) {
		autoCodeTaskNumber, err := crs.AddAutoCodeCheckingTask(creator, course.AutoCodeCheckingTaskCreationParams{})
		require.NoError(t, err)
		testingTaskNumber, err := crs.AddTestingTask(creator, course.TestingTaskCreationParams{
			TestPoints: []course.TestPoint{course.MustNewTestPoint("Spring is DI container", []string{"Yes", "No"}, []int{0})},
		})
		require.NoError(t, err)
		return testingTaskNumber, autoCodeTaskNumber
	}
	newTestPoints := []course.TestPoint{course.MustNewTestPoint("Spring is just DI container", []string{"Yes", "No"}, []int{1})}
	testCases := []struct {
		Name          string
		Academic      course.Academic
		NewTestPoints []course.TestPoint
		PrepareTask   func(crs *course.Course) int
		IsErr         func(err error) bool
	}{
		{
			Name:          "replace_task_test_points_with_new_test_points",
			Academic:      creator,
			NewTestPoints: newTestPoints,
			PrepareTask: func(crs *course.Course) int {
				testingTaskNumber, _ := addTasks(crs)
				return testingTaskNumber
			},
		},
		{
			Name:          "task_has_no_test_points",
			Academic:      creator,
			NewTestPoints: newTestPoints,
			PrepareTask: func(crs *course.Course) int {
				_, autoCodeTaskNumber := addTasks(crs)
				return autoCodeTaskNumber
			},
			IsErr: func(err error) bool {
				return errors.Is(err, course.ErrTaskHasNoTestPoints)
			},
		},
		{
			Name:          "academic_cant_replace_test_points",
			Academic:      course.MustNewAcademic("other-teacher-id", course.TeacherType),
			NewTestPoints: newTestPoints,
			PrepareTask: func(crs *course.Course) int {
				testingTaskNumber, _ := addTasks(crs)
				return testingTaskNumber
			},
			IsErr: course.IsAcademicCantEditCourseError,
		},
		{
			Name:          "no_task_with_number",
			Academic:      creator,
			NewTestPoints: newTestPoints,
			PrepareTask: func(crs *course.Course) int {
				return crs.TasksNumber() + 1
			},
			IsErr: func(err error) bool {
				return errors.Is(err, course.ErrCourseHasNoSuchTask)
			},
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			crs := course.MustNewCourse(course.CreationParams{
				ID:      "course-id",
				Creator: creator,
				Title:   "Spring framework",
				Period:  course.MustNewPeriod(2021, 2022, course.FirstSemester),
			})
			taskNumber := c.PrepareTask(crs)

			err := crs.ReplaceTaskTestPoints(c.Academic, taskNumber, c.NewTestPoints)

			if c.IsErr != nil {
				require.Error(t, err)
				require.True(t, c.IsErr(err))
				return
			}
			require.NoError(t, err)
			task, err := crs.Task(taskNumber)
			require.NoError(t, err)
			if task.Type() != course.TestingType {
				panic("unreachable")
			}
			require.Equal(t, c.NewTestPoints, task.TestingOptional())
		})
	}
}
