package query

import (
	"time"

	"github.com/authena-ru/courses-organization/internal/domain/course"
)

type SpecificTask struct {
	Number      int
	Title       string
	Description string
	Type        course.TaskType
	Deadline    *Deadline
	TestData    []TestData
	Points      []TestPoint
}

type Deadline struct {
	ExcellentGradeTime time.Time
	GoodGradeTime      time.Time
}

type TestData struct {
	InputData  string
	OutputData string
}

type TestPoint struct {
	Description           string
	Variants              []string
	CorrectVariantNumbers []int
	SingleCorrectVariant  bool
}
