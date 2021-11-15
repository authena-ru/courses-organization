package mongodb_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"

	"github.com/authena-ru/courses-organization/internal/adapter/repository/mongodb"
	"github.com/authena-ru/courses-organization/internal/app"
	"github.com/authena-ru/courses-organization/internal/app/command"
	"github.com/authena-ru/courses-organization/internal/app/query"
	"github.com/authena-ru/courses-organization/internal/domain/course"
)

type CoursesRepositoryTestSuite struct {
	suite.Suite
	MongoTestFixtures

	repository *mongodb.CoursesRepository
}

func (s *CoursesRepositoryTestSuite) SetupTest() {
	s.repository = s.newCoursesRepository()
}

func (s *CoursesRepositoryTestSuite) TearDownTest() {
	err := s.repository.RemoveAllCourses(context.Background())
	s.Require().NoError(err)
}

func TestCoursesRepository(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration tests are skipped")
	}

	suite.Run(t, &CoursesRepositoryTestSuite{
		MongoTestFixtures: MongoTestFixtures{t: t},
	})
}

func (s *CoursesRepositoryTestSuite) TestCoursesRepository_FindCourse() {
	courseTeacher := course.MustNewAcademic("623bac67-d4b7-4474-8236-d4e5e3d2b374", course.TeacherType)
	courseStudent := course.MustNewAcademic("ad2ee15f-0805-42dc-8b7c-003c440e88c7", course.StudentType)
	anotherAcademic := course.MustNewAcademic("1f5f202a-8935-44eb-9544-ae78c4e90d46", course.StudentType)

	existingCourse := course.MustNewCourse(course.CreationParams{
		ID:            "a551847e-c529-48a3-bf4f-f5e91c888fa0",
		Creator:       course.MustNewAcademic("f7f02ce0-9d3f-4843-a379-a7773c29e7c2", course.TeacherType),
		Title:         "New brand existing course",
		Period:        course.MustNewPeriod(2041, 2042, course.SecondSemester),
		Started:       true,
		Collaborators: []string{courseTeacher.ID()},
		Students:      []string{courseStudent.ID()},
	})

	s.addCourses(existingCourse)

	testCases := []struct {
		Name           string
		Academic       course.Academic
		CourseIDToFind string
		ExpectedErr    error
	}{
		{
			Name:           "found_course_for_teacher",
			Academic:       courseTeacher,
			CourseIDToFind: existingCourse.ID(),
		},
		{
			Name:           "found_course_for_student",
			Academic:       courseStudent,
			CourseIDToFind: existingCourse.ID(),
		},
		{
			Name:           "not_found_course_for_another_academic",
			Academic:       anotherAcademic,
			CourseIDToFind: existingCourse.ID(),
			ExpectedErr:    app.ErrCourseDoesntExist,
		},
		{
			Name:           "course_doesnt_exist",
			Academic:       courseTeacher,
			CourseIDToFind: "27d5cd20-9da5-4636-85d5-12cd07a0c0ff",
			ExpectedErr:    app.ErrCourseDoesntExist,
		},
	}

	for _, c := range testCases {
		s.Run(c.Name, func() {
			foundCourse, err := s.repository.FindCourse(context.Background(), c.Academic, c.CourseIDToFind)

			if c.ExpectedErr != nil {
				s.Require().True(errors.Is(err, c.ExpectedErr))

				return
			}

			s.Require().NoError(err)
			s.Require().Equal(existingCourse.ID(), foundCourse.ID)
		})
	}
}

func (s *CoursesRepositoryTestSuite) TestCoursesRepository_FindAllCourses() {
	courseTeacher := course.MustNewAcademic("fdd8bd0a-a6da-4ba9-b60a-cd1093a480bf", course.TeacherType)
	courseStudent := course.MustNewAcademic("244a26f9-9c33-4585-843b-822b04d4cb76", course.StudentType)
	anotherAcademic := course.MustNewAcademic("e6377c3b-90ae-4028-896a-c867d7dc0738", course.StudentType)

	physicsCourse := course.MustNewCourse(course.CreationParams{
		ID:            "270d886e-097c-4334-8b0e-496464645dec",
		Creator:       course.MustNewAcademic("b639b230-80c1-456d-9f8e-e84459abc7a7", course.TeacherType),
		Title:         "Physics course",
		Period:        course.MustNewPeriod(2021, 2022, course.FirstSemester),
		Started:       true,
		Collaborators: []string{courseTeacher.ID()},
		Students:      []string{courseStudent.ID()},
	})
	mathCourse := course.MustNewCourse(course.CreationParams{
		ID:            "47b36fbb-61e3-42fe-bb72-3908b909542f",
		Creator:       course.MustNewAcademic("b3a30eb3-ea5c-441b-adba-a36a834494c4", course.TeacherType),
		Title:         "Math course",
		Period:        course.MustNewPeriod(2021, 2022, course.FirstSemester),
		Started:       true,
		Collaborators: []string{courseTeacher.ID()},
		Students:      []string{courseStudent.ID()},
	})

	s.addCourses(physicsCourse, mathCourse)

	testCases := []struct {
		Name         string
		Academic     course.Academic
		FilterParams query.CoursesFilterParams
		ExpectedIDs  []string
	}{
		{
			Name:         "found_courses_for_teacher",
			Academic:     courseTeacher,
			FilterParams: query.CoursesFilterParams{},
			ExpectedIDs:  []string{"270d886e-097c-4334-8b0e-496464645dec", "47b36fbb-61e3-42fe-bb72-3908b909542f"},
		},
		{
			Name:         "found_courses_for_student",
			Academic:     courseStudent,
			FilterParams: query.CoursesFilterParams{},
			ExpectedIDs:  []string{"270d886e-097c-4334-8b0e-496464645dec", "47b36fbb-61e3-42fe-bb72-3908b909542f"},
		},
		{
			Name:         "not_found_courses_for_another_academic",
			Academic:     anotherAcademic,
			FilterParams: query.CoursesFilterParams{},
		},
		{
			Name:         "found_course_by_title",
			Academic:     courseStudent,
			FilterParams: query.CoursesFilterParams{Title: "phys"},
			ExpectedIDs:  []string{"270d886e-097c-4334-8b0e-496464645dec"},
		},
		{
			Name:         "not_found_course_by_title",
			Academic:     courseTeacher,
			FilterParams: query.CoursesFilterParams{Title: "kqu"},
		},
	}

	for _, c := range testCases {
		s.Run(c.Name, func() {
			courses, err := s.repository.FindAllCourses(context.Background(), c.Academic, c.FilterParams)
			s.Require().NoError(err)

			s.Require().ElementsMatch(c.ExpectedIDs, mapCommonCoursesToIDS(courses))
		})
	}
}

func mapCommonCoursesToIDS(courses []app.CommonCourse) []string {
	ids := make([]string, 0, len(courses))
	for _, c := range courses {
		ids = append(ids, c.ID)
	}

	return ids
}

func (s *CoursesRepositoryTestSuite) TestCoursesRepository_AddCourse() {
	testCases := []struct {
		Name           string
		CoursesFactory func() *course.Course
	}{
		{
			Name: "course_without_collaborators_and_students",
			CoursesFactory: func() *course.Course {
				return course.MustNewCourse(course.CreationParams{
					ID:      "8f0d3461-175e-4e34-aebe-f292d83c3957",
					Creator: course.MustNewAcademic("aec5b656-1028-472c-a2be-f4bd2b0424b1", course.TeacherType),
					Title:   "New brand course",
					Period:  course.MustNewPeriod(2023, 2024, course.FirstSemester),
					Started: false,
				})
			},
		},
		{
			Name: "course_with_collaborators",
			CoursesFactory: func() *course.Course {
				return course.MustNewCourse(course.CreationParams{
					ID:      "1e2339f9-7d31-42c9-bbfc-302c1471f48b",
					Creator: course.MustNewAcademic("93fedf0c-ae0d-4301-81b5-1f2291bdfa32", course.TeacherType),
					Title:   "New brand course with collaborators",
					Period:  course.MustNewPeriod(2024, 2025, course.SecondSemester),
					Started: true,
					Collaborators: []string{
						"8452d33b-a739-42d5-8090-627e080176b0",
						"1793f030-fedf-48b1-a996-7f8a2ff37a6f",
						"86094876-1d58-4a10-a7c7-4c0cb4e6584d",
					},
				})
			},
		},
		{
			Name: "course_with_students",
			CoursesFactory: func() *course.Course {
				return course.MustNewCourse(course.CreationParams{
					ID:      "7a7d3715-a54d-40f7-a568-95fa823b597e",
					Creator: course.MustNewAcademic("32a1be5d-62d5-4554-8784-507e95679966", course.TeacherType),
					Title:   "New brand course with students",
					Period:  course.MustNewPeriod(2025, 2026, course.FirstSemester),
					Started: true,
					Students: []string{
						"c5206f35-32d7-4ae8-903f-c4a2cad4c325",
						"57219c45-f7dc-4ac1-8695-d8aa2bdeb151",
						"e4841ac4-2bb6-4dd4-aabe-0e2a735ee613",
						"e81347e9-55f1-44a2-8adc-e8d65dd8b5ae",
					},
				})
			},
		},
	}

	for _, c := range testCases {
		s.Run(c.Name, func() {
			expectedCourse := c.CoursesFactory()

			err := s.repository.AddCourse(context.Background(), expectedCourse)
			s.Require().NoError(err)

			s.requirePersistedCourseEquals(expectedCourse)
		})
	}
}

func (s *CoursesRepositoryTestSuite) TestCoursesRepository_GetCourse() {
	existingCourse := course.MustNewCourse(course.CreationParams{
		ID:      "dfe52190-ed85-4130-aba7-409c0bc21a49",
		Creator: course.MustNewAcademic("7b07ce84-d8aa-454b-9d01-305b211fe730", course.TeacherType),
		Title:   "New brand existing course",
		Period:  course.MustNewPeriod(2027, 2028, course.SecondSemester),
		Started: false,
	})

	s.addCourses(existingCourse)

	testCases := []struct {
		Name          string
		CourseIDToGet string
		ExpectedErr   error
	}{
		{
			Name:          "existing_course",
			CourseIDToGet: existingCourse.ID(),
		},
		{
			Name:          "non_existing_course",
			CourseIDToGet: "23c645b3-9c94-4d5b-86e5-c8a3f9a47db3",
			ExpectedErr:   app.ErrCourseDoesntExist,
		},
	}

	for _, c := range testCases {
		s.Run(c.Name, func() {
			givenCourse, err := s.repository.GetCourse(context.Background(), c.CourseIDToGet)

			if c.ExpectedErr != nil {
				s.Require().True(errors.Is(err, c.ExpectedErr))

				return
			}

			s.Require().NoError(err)
			s.Require().Equal(existingCourse, givenCourse)
		})
	}
}

func (s *CoursesRepositoryTestSuite) TestCoursesRepository_UpdateCourse() {
	courseToUpdate := course.MustNewCourse(course.CreationParams{
		ID:            "a9ddc049-38c7-474f-b4b7-593b686482ea",
		Creator:       course.MustNewAcademic("243886c7-5215-4574-9cb5-50c84cf7e2fe", course.TeacherType),
		Title:         "New course that will change",
		Period:        course.MustNewPeriod(2030, 2031, course.FirstSemester),
		Started:       false,
		Collaborators: []string{"5a475f3f-c8cf-477e-9ba5-0e428923b089"},
		Students:      []string{"d4a28d9b-7d05-4639-a067-4a5e5cb4e057"},
	})

	s.addCourses(courseToUpdate)

	testCases := []struct {
		Name            string
		UpdatedCourseID string
		Update          command.UpdateFunction
		CourseUpdated   func(crs *course.Course) bool
		ExpectedErr     error
	}{
		{
			Name:            "non_existing_course",
			UpdatedCourseID: "b59aa2a2-2d2a-4cd8-916b-234fdb2827a7",
			Update: func(_ context.Context, crs *course.Course) (*course.Course, error) {
				return crs, nil
			},
			ExpectedErr: app.ErrCourseDoesntExist,
		},
	}

	for _, c := range testCases {
		s.Run(c.Name, func() {
			ctx := context.Background()

			err := s.repository.UpdateCourse(ctx, c.UpdatedCourseID, c.Update)

			if c.ExpectedErr != nil {
				s.Require().True(errors.Is(err, c.ExpectedErr))

				return
			}

			s.Require().NoError(err)

			courseFromDB, err := s.repository.GetCourse(ctx, c.UpdatedCourseID)
			s.Require().NoError(err)

			s.Require().True(c.CourseUpdated(courseFromDB))
		})
	}
}

func (s *CoursesRepositoryTestSuite) newCoursesRepository() *mongodb.CoursesRepository {
	return mongodb.NewCoursesRepository(s.db)
}

func (s *CoursesRepositoryTestSuite) addCourses(courses ...*course.Course) {
	s.T().Helper()

	ctx := context.Background()

	for _, c := range courses {
		err := s.repository.AddCourse(ctx, c)
		s.Require().NoError(err)
	}
}

func (s *CoursesRepositoryTestSuite) requirePersistedCourseEquals(expectedCourse *course.Course) {
	s.T().Helper()

	persistedCourse, err := s.repository.GetCourse(context.Background(), expectedCourse.ID())
	s.Require().NoError(err)

	s.Require().Equal(expectedCourse, persistedCourse)
}
