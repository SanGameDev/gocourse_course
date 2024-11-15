package course

import (
	"context"
	"log"
	"time"

	"github.com/SanGameDev/gocourse_domain/domain"
)

type (
	Service interface {
		Create(ctx context.Context, name, startDate, endDate string) (*domain.Course, error)
		GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Course, error)
		Get(ctx context.Context, id string) (*domain.Course, error)
		Delete(ctx context.Context, id string) error
		Update(ctx context.Context, id string, name, startDate, endDate *string) error
		Count(ctx context.Context, filters Filters) (int, error)
	}

	service struct {
		log  *log.Logger
		repo Repository
	}

	Filters struct {
		Name string
	}
)

func NewService(log *log.Logger, repo Repository) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (s service) Create(ctx context.Context, name, startDate, endDate string) (*domain.Course, error) {
	s.log.Println("Creating courses")

	startDateParsed, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		s.log.Println("Error parsing start date: ", err)
		return nil, ErrInvalidStartDate
	}

	endDateParsed, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		s.log.Println("Error parsing end date: ", err)
		return nil, ErrInvalidEndDate
	}

	if startDateParsed.After(endDateParsed) {
		s.log.Println(ErrEndLesserStart)
		return nil, ErrEndLesserStart
	}

	course := domain.Course{
		Name:      name,
		StartDate: startDateParsed,
		EndDate:   endDateParsed,
	}

	if err := s.repo.Create(ctx, &course); err != nil {
		return nil, err
	}

	return &course, nil
}

func (s service) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Course, error) {
	s.log.Println("Getting all courses")
	courses, err := s.repo.GetAll(ctx, filters, offset, limit)
	if err != nil {
		return nil, err
	}
	return courses, nil
}

func (s service) Get(ctx context.Context, id string) (*domain.Course, error) {
	s.log.Println("Getting course with id: ", id)
	course, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return course, nil
}

func (s service) Delete(ctx context.Context, id string) error {
	s.log.Println("Deleting course with id: ", id)
	return s.repo.Delete(ctx, id)
}

func (s service) Update(ctx context.Context, id string, name *string, startDate, endDate *string) error {
	s.log.Println("Updating course with id: ", id)

	var startDateParsed, endDateParsed *time.Time

	course, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}

	if startDate != nil {
		date, err := time.Parse("2006-01-02", *startDate)
		if err != nil {
			s.log.Println("Error parsing start date: ", err)
			return ErrInvalidStartDate
		}

		if date.After(course.EndDate) {
			s.log.Println(ErrEndLesserStart)
			return ErrEndLesserStart
		}

		startDateParsed = &date
	}

	if endDate != nil {
		date, err := time.Parse("2006-01-02", *endDate)
		if err != nil {
			s.log.Println("Error parsing end date: ", err)
			return ErrInvalidEndDate
		}

		if course.EndDate.After(date) {
			s.log.Println(ErrEndLesserStart)
			return ErrEndLesserStart
		}

		endDateParsed = &date
	}

	if err := s.repo.Update(ctx, id, name, startDateParsed, endDateParsed); err != nil {
		return err
	}
	return nil
}

func (s service) Count(ctx context.Context, filters Filters) (int, error) {
	s.log.Println("Counting courses")
	return s.repo.Count(ctx, filters)
}
