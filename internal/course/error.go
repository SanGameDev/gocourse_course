package course

import (
	"errors"
	"fmt"
)

var ErrInvalidStartDate = errors.New("invalid start date")
var ErrInvalidEndDate = errors.New("invalid end date")
var ErrNameRequired = errors.New("name is required")
var ErrStartDateRequired = errors.New("start date is required")
var ErrEndDateRequired = errors.New("end date is required")
var ErrEndLesserStart = errors.New("end date mustn't be lesser than start date")

type ErrNotFound struct {
	CourseID string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("course with id '%s' not found or doesn't exist", e.CourseID)
}
