package app

import (
	"time"

	"github.com/authena-ru/courses-organization/internal/domain/course"
)

type (
	CommonCourse struct {
		ID          string
		Title       string
		Period      Period
		CreatorID   string
		Started     bool
		TasksNumber int
	}

	SpecificTask struct {
		Number      int
		Title       string
		Description string
		Type        course.TaskType
		Deadline    *Deadline
		TestData    []TestData
		Points      []TestPoint
	}

	GeneralTask struct {
		Number      int
		Title       string
		Description string
		Type        course.TaskType
	}

	Period struct {
		AcademicStartYear int
		AcademicEndYear   int
		Semester          course.Semester
	}

	Deadline struct {
		ExcellentGradeTime time.Time
		GoodGradeTime      time.Time
	}

	TestData struct {
		InputData  string
		OutputData string
	}

	TestPoint struct {
		Description           string
		Variants              []string
		CorrectVariantNumbers []int
		SingleCorrectVariant  bool
	}
)
