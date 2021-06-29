package course

import "github.com/pkg/errors"

type TestPoint struct {
	description           string
	variants              []string
	correctVariantNumbers []int
}

const testPointDescriptionMaxLen = 500

var (
	ErrTestPointDescriptionTooLong     = errors.New("test point description too long")
	ErrEmptyTestPointVariants          = errors.New("empty test point variants")
	ErrEmptyTestPointCorrectVariants   = errors.New("test points has no correct variants")
	ErrTooMuchTestPointCorrectVariants = errors.New("test point has too much correct variants")
	ErrInvalidTestPointVariantNumber   = errors.New("invalid test point variant number")
)

func IsInvalidTestPointError(err error) bool {
	return errors.Is(err, ErrTestPointDescriptionTooLong) ||
		errors.Is(err, ErrEmptyTestPointVariants) ||
		errors.Is(err, ErrEmptyTestPointCorrectVariants) ||
		errors.Is(err, ErrTooMuchTestPointCorrectVariants) ||
		errors.Is(err, ErrInvalidTestPointVariantNumber)
}

func NewTestPoint(description string, variants []string, correctVariantNumbers []int) (TestPoint, error) {
	if len(description) > testPointDescriptionMaxLen {
		return TestPoint{}, ErrTestPointDescriptionTooLong
	}

	if len(variants) == 0 {
		return TestPoint{}, ErrEmptyTestPointVariants
	}

	if len(correctVariantNumbers) == 0 {
		return TestPoint{}, ErrEmptyTestPointCorrectVariants
	}

	if len(correctVariantNumbers) > len(variants) {
		return TestPoint{}, ErrTooMuchTestPointCorrectVariants
	}

	for _, n := range correctVariantNumbers {
		if n >= len(variants) || n < 0 {
			return TestPoint{}, ErrInvalidTestPointVariantNumber
		}
	}

	variantsCopy := make([]string, len(variants))
	copy(variantsCopy, variants)

	correctVariantNumbersCopy := make([]int, len(correctVariantNumbers))
	copy(correctVariantNumbersCopy, correctVariantNumbers)

	return TestPoint{
		description:           description,
		variants:              variantsCopy,
		correctVariantNumbers: correctVariantNumbersCopy,
	}, nil
}

func MustNewTestPoint(description string, variants []string, correctVariantNumbers []int) TestPoint {
	tp, err := NewTestPoint(description, variants, correctVariantNumbers)
	if err != nil {
		panic(err)
	}

	return tp
}

func (tp TestPoint) Description() string {
	return tp.description
}

func (tp TestPoint) Variants() []string {
	variantsCopy := make([]string, len(tp.variants))
	copy(variantsCopy, tp.variants)

	return variantsCopy
}

func (tp TestPoint) CorrectVariantNumbers() []int {
	correctVariantNumbersCopy := make([]int, len(tp.correctVariantNumbers))
	copy(correctVariantNumbersCopy, tp.correctVariantNumbers)

	return correctVariantNumbersCopy
}

func (tp TestPoint) IsZero() bool {
	return tp.description == "" && len(tp.variants) == 0 && len(tp.correctVariantNumbers) == 0
}
