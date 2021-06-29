package course

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

type Deadline struct {
	excellentGradeTime time.Time
	goodGradeTime      time.Time
}

var (
	ErrZeroExcellentGradeTime      = errors.New("zero excellent grade time")
	ErrZeroGoodGradeTime           = errors.New("zero good grade time")
	ErrExcellentGradeTimeAfterGood = errors.New("excellent grade time after good")
)

func IsInvalidDeadlineError(err error) bool {
	return errors.Is(err, ErrZeroExcellentGradeTime) ||
		errors.Is(err, ErrZeroGoodGradeTime) ||
		errors.Is(err, ErrExcellentGradeTimeAfterGood)
}

func NewDeadline(excellentGradeTime time.Time, goodGradeTime time.Time) (Deadline, error) {
	if excellentGradeTime.IsZero() {
		return Deadline{}, ErrZeroExcellentGradeTime
	}

	if goodGradeTime.IsZero() {
		return Deadline{}, ErrZeroGoodGradeTime
	}

	if excellentGradeTime.After(goodGradeTime) {
		return Deadline{}, ErrExcellentGradeTimeAfterGood
	}

	return Deadline{
		excellentGradeTime: excellentGradeTime,
		goodGradeTime:      goodGradeTime,
	}, nil
}

func MustNewDeadline(excellentGradeTime time.Time, goodGradeTime time.Time) Deadline {
	deadline, err := NewDeadline(excellentGradeTime, goodGradeTime)
	if err != nil {
		panic(err)
	}

	return deadline
}

func (d Deadline) GoodGradeTime() time.Time {
	return d.goodGradeTime
}

func (d Deadline) ExcellentGradeTime() time.Time {
	return d.excellentGradeTime
}

func (d Deadline) IsZero() bool {
	return d == Deadline{}
}

func (d Deadline) String() string {
	return fmt.Sprintf("excellent — %s, good — %s", d.excellentGradeTime, d.goodGradeTime)
}
