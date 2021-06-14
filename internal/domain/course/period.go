package course

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

type Semester uint8

const (
	FirstSemester Semester = iota + 1
	SecondSemester
)

func (s Semester) String() string {
	switch s {
	case FirstSemester:
		return "1st semester"
	case SecondSemester:
		return "2nd semester"
	}
	return "%!Semester(" + strconv.Itoa(int(s)) + ")"
}

type Period struct {
	academicStartYear int
	academicEndYear   int
	semester          Semester
}

var (
	ErrStartYearAfterEnd      = errors.New("academic start year after end")
	ErrYearDurationOverYear   = errors.New("academic year duration over year")
	ErrStartYearEqualsEndYear = errors.New("academic start year equals end year")
)

func NewPeriod(academicStartYear, academicEndYear int, semester Semester) (Period, error) {
	if academicStartYear > academicEndYear {
		return Period{}, ErrStartYearAfterEnd
	}
	if academicEndYear-academicStartYear > 1 {
		return Period{}, ErrYearDurationOverYear
	}
	if academicStartYear == academicEndYear {
		return Period{}, ErrStartYearEqualsEndYear
	}
	return Period{
		academicStartYear: academicStartYear,
		academicEndYear:   academicEndYear,
		semester:          semester,
	}, nil
}

func MustNewPeriod(academicStartYear, academicEndYear int, semester Semester) Period {
	p, err := NewPeriod(academicStartYear, academicEndYear, semester)
	if err != nil {
		panic(err)
	}
	return p
}

func (p Period) IsZero() bool {
	return p == Period{}
}

func (p Period) AcademicStartYear() int {
	return p.academicStartYear
}

func (p Period) AcademicEndYear() int {
	return p.academicEndYear
}

func (p Period) Semester() Semester {
	return p.semester
}

func (p Period) next() Period {
	return Period{
		academicStartYear: p.academicEndYear,
		academicEndYear:   p.academicEndYear + 1,
		semester:          p.semester,
	}
}

func (p Period) String() string {
	return fmt.Sprintf("%d-%d %s", p.academicStartYear, p.academicEndYear, p.semester)
}
