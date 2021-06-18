package v1

import (
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/go-chi/render"

	"github.com/authena-ru/courses-organization/internal/app/query"
	"github.com/authena-ru/courses-organization/internal/domain/course"
)

func marshallSpecificTask(w http.ResponseWriter, r *http.Request, task query.SpecificTask) {
	type taskResponse struct {
		TaskResponse
		Deadline *Deadline   `json:"deadline,omitempty"`
		TestData []TestData  `json:"testData,omitempty"`
		Points   []TestPoint `json:"points,omitempty"`
	}
	response := taskResponse{
		TaskResponse: TaskResponse{
			Number: task.Number,
			Task: Task{
				Title:       task.Title,
				Description: task.Description,
				Type:        marshallTaskType(task.Type),
			},
		},
		Deadline: marshallDeadline(task.Deadline),
		TestData: marshallTestData(task.TestData),
		Points:   marshallTestPoints(task.Points),
	}
	render.Respond(w, r, response)
}

func marshallTaskType(taskType course.TaskType) TaskType {
	switch taskType {
	case course.ManualCheckingType:
		return TaskTypeMANUALCHECKING
	case course.AutoCodeCheckingType:
		return TaskTypeAUTOCODECHECKING
	case course.TestingType:
		return TaskTypeTESTING
	}
	return "UNKNOWN"
}

func marshallDeadline(deadline *query.Deadline) *Deadline {
	if deadline == nil {
		return nil
	}
	return &Deadline{
		ExcellentGradeTime: types.Date{Time: deadline.ExcellentGradeTime},
		GoodGradeTime:      types.Date{Time: deadline.GoodGradeTime},
	}
}

func marshallTestData(testData []query.TestData) []TestData {
	marshalled := make([]TestData, 0, len(testData))
	for _, td := range testData {
		marshalled = append(marshalled, TestData{
			InputData:  &td.InputData,
			OutputData: &td.OutputData,
		})
	}
	return marshalled
}

func marshallTestPoints(testPoints []query.TestPoint) []TestPoint {
	marshalled := make([]TestPoint, 0, len(testPoints))
	for _, tp := range testPoints {
		var correctVariantNumbers *[]int
		if tp.CorrectVariantNumbers != nil {
			correctVariantNumbers = &tp.CorrectVariantNumbers
		}
		marshalled = append(marshalled, TestPoint{
			Description:           tp.Description,
			Variants:              tp.Variants,
			CorrectVariantNumbers: correctVariantNumbers,
			SingleCorrectVariant:  &tp.SingleCorrectVariant,
		})
	}
	return marshalled
}
