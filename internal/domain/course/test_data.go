package course

import "github.com/pkg/errors"

type TestData struct {
	inputData  string
	outputData string
}

const testDataMaxLen = 1000

var (
	ErrTestInputDataTooLong  = errors.New("test input data too long")
	ErrTestOutputDataTooLong = errors.New("test output data too long")
)

func IsInvalidTestDataError(err error) bool {
	return errors.Is(err, ErrTestInputDataTooLong) ||
		errors.Is(err, ErrTestOutputDataTooLong)
}

func NewTestData(inputData, outputData string) (TestData, error) {
	if len(inputData) > testDataMaxLen {
		return TestData{}, ErrTestInputDataTooLong
	}

	if len(outputData) > testDataMaxLen {
		return TestData{}, ErrTestOutputDataTooLong
	}

	return TestData{
		inputData:  inputData,
		outputData: outputData,
	}, nil
}

func MustNewTestData(inputData, outputData string) TestData {
	td, err := NewTestData(inputData, outputData)
	if err != nil {
		panic(err)
	}

	return td
}

func (td TestData) InputData() string {
	return td.inputData
}

func (td TestData) OutputData() string {
	return td.outputData
}

func (td TestData) IsZero() bool {
	return td == TestData{}
}
