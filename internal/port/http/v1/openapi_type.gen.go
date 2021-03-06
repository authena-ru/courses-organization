// Package v1 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.8.2 DO NOT EDIT.
package v1

import (
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
)

const (
	BearerAuthScopes = "bearerAuth.Scopes"
)

// Defines values for ResourceType.
const (
	ResourceTypeOTHER ResourceType = "OTHER"

	ResourceTypePRESENTATION ResourceType = "PRESENTATION"

	ResourceTypeTRAININGMANUAL ResourceType = "TRAINING_MANUAL"

	ResourceTypeVIDEO ResourceType = "VIDEO"
)

// Defines values for Semester.
const (
	SemesterFIRST Semester = "FIRST"

	SemesterSECOND Semester = "SECOND"
)

// Defines values for TaskType.
const (
	TaskTypeAUTOCODECHECKING TaskType = "AUTO_CODE_CHECKING"

	TaskTypeMANUALCHECKING TaskType = "MANUAL_CHECKING"

	TaskTypeTESTING TaskType = "TESTING"
)

// AddAutoCodeCheckingTaskRequest defines model for AddAutoCodeCheckingTaskRequest.
type AddAutoCodeCheckingTaskRequest struct {
	// Embedded struct due to allOf(#/components/schemas/AddTaskRequest)
	AddTaskRequest `yaml:",inline"`
	// Embedded struct due to allOf(#/components/schemas/AutoCodeCheckingTaskPart)
	AutoCodeCheckingTaskPart `yaml:",inline"`
}

// AddAuxiliaryMaterialRequest defines model for AddAuxiliaryMaterialRequest.
type AddAuxiliaryMaterialRequest AuxiliaryMaterial

// AddCollaboratorToCourseRequest defines model for AddCollaboratorToCourseRequest.
type AddCollaboratorToCourseRequest struct {
	Id string `json:"id"`
}

// AddGroupToCourseRequest defines model for AddGroupToCourseRequest.
type AddGroupToCourseRequest struct {
	Id string `json:"id"`
}

// AddManualCheckingTaskRequest defines model for AddManualCheckingTaskRequest.
type AddManualCheckingTaskRequest struct {
	// Embedded struct due to allOf(#/components/schemas/AddTaskRequest)
	AddTaskRequest `yaml:",inline"`
	// Embedded struct due to allOf(#/components/schemas/ManualCheckingTaskPart)
	ManualCheckingTaskPart `yaml:",inline"`
}

// AddStudentToCourseRequest defines model for AddStudentToCourseRequest.
type AddStudentToCourseRequest struct {
	Id string `json:"id"`
}

// AddTaskRequest defines model for AddTaskRequest.
type AddTaskRequest Task

// AddTestingTaskRequest defines model for AddTestingTaskRequest.
type AddTestingTaskRequest struct {
	// Embedded struct due to allOf(#/components/schemas/AddTaskRequest)
	AddTaskRequest `yaml:",inline"`
	// Embedded struct due to allOf(#/components/schemas/TestingTaskPart)
	TestingTaskPart `yaml:",inline"`
}

// AutoCodeCheckingTaskPart defines model for AutoCodeCheckingTaskPart.
type AutoCodeCheckingTaskPart struct {
	Deadline *Deadline   `json:"deadline,omitempty"`
	TestData *[]TestData `json:"testData,omitempty"`
}

// AutoCodeCheckingTaskResponse defines model for AutoCodeCheckingTaskResponse.
type AutoCodeCheckingTaskResponse struct {
	// Embedded struct due to allOf(#/components/schemas/TaskResponse)
	TaskResponse `yaml:",inline"`
	// Embedded struct due to allOf(#/components/schemas/AutoCodeCheckingTaskPart)
	AutoCodeCheckingTaskPart `yaml:",inline"`
}

// AuxiliaryMaterial defines model for AuxiliaryMaterial.
type AuxiliaryMaterial struct {
	Resource     string       `json:"resource"`
	ResourceType ResourceType `json:"resourceType"`
}

// Course defines model for Course.
type Course struct {
	CreatorId   string       `json:"creatorId"`
	Id          string       `json:"id"`
	Period      CoursePeriod `json:"period"`
	Started     bool         `json:"started"`
	TasksNumber int          `json:"tasksNumber"`
	Title       string       `json:"title"`
}

// CoursePeriod defines model for CoursePeriod.
type CoursePeriod struct {
	AcademicEndYear   int      `json:"academicEndYear"`
	AcademicStartYear int      `json:"academicStartYear"`
	Semester          Semester `json:"semester"`
}

// CreateCourseRequest defines model for CreateCourseRequest.
type CreateCourseRequest struct {
	Period  CoursePeriod `json:"period"`
	Started bool         `json:"started"`
	Title   string       `json:"title"`
}

// Deadline defines model for Deadline.
type Deadline struct {
	ExcellentGradeTime openapi_types.Date `json:"excellentGradeTime"`
	GoodGradeTime      openapi_types.Date `json:"goodGradeTime"`
}

// EditCourseRequest defines model for EditCourseRequest.
type EditCourseRequest struct {
	Period  *CoursePeriod `json:"period,omitempty"`
	Started *bool         `json:"started,omitempty"`
	Title   *string       `json:"title,omitempty"`
}

// Error defines model for Error.
type Error struct {
	Details string `json:"details"`
	Slug    string `json:"slug"`
}

// ExtendCourseRequest defines model for ExtendCourseRequest.
type ExtendCourseRequest struct {
	Period  *CoursePeriod `json:"period,omitempty"`
	Started bool          `json:"started"`
	Title   string        `json:"title"`
}

// GetAllAuxiliaryMaterialsResponse defines model for GetAllAuxiliaryMaterialsResponse.
type GetAllAuxiliaryMaterialsResponse []AuxiliaryMaterial

// GetAllCourseCollaboratorsResponse defines model for GetAllCourseCollaboratorsResponse.
type GetAllCourseCollaboratorsResponse []Teacher

// GetAllCourseStudentsResponse defines model for GetAllCourseStudentsResponse.
type GetAllCourseStudentsResponse []Student

// GetAllCoursesResponse defines model for GetAllCoursesResponse.
type GetAllCoursesResponse []Course

// GetCourseResponse defines model for GetCourseResponse.
type GetCourseResponse Course

// ManualCheckingTaskPart defines model for ManualCheckingTaskPart.
type ManualCheckingTaskPart struct {
	Deadline *Deadline `json:"deadline,omitempty"`
}

// ManualCheckingTaskResponse defines model for ManualCheckingTaskResponse.
type ManualCheckingTaskResponse struct {
	// Embedded struct due to allOf(#/components/schemas/TaskResponse)
	TaskResponse `yaml:",inline"`
	// Embedded struct due to allOf(#/components/schemas/ManualCheckingTaskPart)
	ManualCheckingTaskPart `yaml:",inline"`
}

// ResourceType defines model for ResourceType.
type ResourceType string

// Semester defines model for Semester.
type Semester string

// Student defines model for Student.
type Student struct {
	FullName string `json:"fullName"`
	Id       string `json:"id"`
}

// Task defines model for Task.
type Task struct {
	Description string   `json:"description"`
	Title       string   `json:"title"`
	Type        TaskType `json:"type"`
}

// TaskResponse defines model for TaskResponse.
type TaskResponse struct {
	// Embedded struct due to allOf(#/components/schemas/Task)
	Task `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	Number int `json:"number"`
}

// TaskType defines model for TaskType.
type TaskType string

// Teacher defines model for Teacher.
type Teacher struct {
	FullName string `json:"fullName"`
	Id       string `json:"id"`
}

// TestData defines model for TestData.
type TestData struct {
	// property not required in response for student, but required for creation
	InputData *string `json:"inputData,omitempty"`

	// property not required in response for student, but required for creation
	OutputData *string `json:"outputData,omitempty"`
}

// TestPoint defines model for TestPoint.
type TestPoint struct {
	// property not required in response for student, but required for creation
	CorrectVariantNumbers *[]int `json:"correctVariantNumbers,omitempty"`
	Description           string `json:"description"`

	// property indicates that point has single correct variant in response for student
	SingleCorrectVariant *bool    `json:"singleCorrectVariant,omitempty"`
	Variants             []string `json:"variants"`
}

// TestingTaskPart defines model for TestingTaskPart.
type TestingTaskPart struct {
	Points *[]TestPoint `json:"points,omitempty"`
}

// TestingTaskResponse defines model for TestingTaskResponse.
type TestingTaskResponse struct {
	// Embedded struct due to allOf(#/components/schemas/TaskResponse)
	TaskResponse `yaml:",inline"`
	// Embedded struct due to allOf(#/components/schemas/TestingTaskPart)
	TestingTaskPart `yaml:",inline"`
}

// GetAllCoursesParams defines parameters for GetAllCourses.
type GetAllCoursesParams struct {
	// course title substring for filtering
	Title *string `json:"title,omitempty"`
}

// CreateCourseJSONBody defines parameters for CreateCourse.
type CreateCourseJSONBody CreateCourseRequest

// EditCourseJSONBody defines parameters for EditCourse.
type EditCourseJSONBody EditCourseRequest

// GetAllCourseAuxiliaryMaterialsParams defines parameters for GetAllCourseAuxiliaryMaterials.
type GetAllCourseAuxiliaryMaterialsParams struct {
	// resource type for filtering
	ResourceType *ResourceType `json:"resourceType,omitempty"`
}

// AttachAuxiliaryMaterialToCourseJSONBody defines parameters for AttachAuxiliaryMaterialToCourse.
type AttachAuxiliaryMaterialToCourseJSONBody AddAuxiliaryMaterialRequest

// AddCollaboratorToCourseJSONBody defines parameters for AddCollaboratorToCourse.
type AddCollaboratorToCourseJSONBody AddCollaboratorToCourseRequest

// ExtendCourseJSONBody defines parameters for ExtendCourse.
type ExtendCourseJSONBody ExtendCourseRequest

// AddGroupToCourseJSONBody defines parameters for AddGroupToCourse.
type AddGroupToCourseJSONBody AddGroupToCourseRequest

// GetAllCourseStudentsParams defines parameters for GetAllCourseStudents.
type GetAllCourseStudentsParams struct {
	// student full name substring for filtering
	FullName *string `json:"fullName,omitempty"`
}

// AddStudentToCourseJSONBody defines parameters for AddStudentToCourse.
type AddStudentToCourseJSONBody AddStudentToCourseRequest

// GetCourseTasksParams defines parameters for GetCourseTasks.
type GetCourseTasksParams struct {
	// type of task for filtering
	Type *TaskType `json:"type,omitempty"`

	// text for search in tasks title and description
	Text *string `json:"text,omitempty"`
}

// AddTaskToCourseJSONBody defines parameters for AddTaskToCourse.
type AddTaskToCourseJSONBody interface{}

// CreateCourseJSONRequestBody defines body for CreateCourse for application/json ContentType.
type CreateCourseJSONRequestBody CreateCourseJSONBody

// EditCourseJSONRequestBody defines body for EditCourse for application/json ContentType.
type EditCourseJSONRequestBody EditCourseJSONBody

// AttachAuxiliaryMaterialToCourseJSONRequestBody defines body for AttachAuxiliaryMaterialToCourse for application/json ContentType.
type AttachAuxiliaryMaterialToCourseJSONRequestBody AttachAuxiliaryMaterialToCourseJSONBody

// AddCollaboratorToCourseJSONRequestBody defines body for AddCollaboratorToCourse for application/json ContentType.
type AddCollaboratorToCourseJSONRequestBody AddCollaboratorToCourseJSONBody

// ExtendCourseJSONRequestBody defines body for ExtendCourse for application/json ContentType.
type ExtendCourseJSONRequestBody ExtendCourseJSONBody

// AddGroupToCourseJSONRequestBody defines body for AddGroupToCourse for application/json ContentType.
type AddGroupToCourseJSONRequestBody AddGroupToCourseJSONBody

// AddStudentToCourseJSONRequestBody defines body for AddStudentToCourse for application/json ContentType.
type AddStudentToCourseJSONRequestBody AddStudentToCourseJSONBody

// AddTaskToCourseJSONRequestBody defines body for AddTaskToCourse for application/json ContentType.
type AddTaskToCourseJSONRequestBody AddTaskToCourseJSONBody
