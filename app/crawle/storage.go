package crawle

import (
	"sync"
	"time"

	"github.com/earlgray283/gakujo-google-calendar/gakujo/model"
)

type MinitestStorage struct {
	sync.Mutex
	rows      []model.MinitestRow
	updatedAt time.Time
}

type ReportStorage struct {
	sync.Mutex
	rows      []model.ReportRow
	updatedAt time.Time
}

type ClassEnqStorage struct {
	sync.Mutex
	rows      []model.ClassEnqRow
	updatedAt time.Time
}

func (s *MinitestStorage) Get() ([]model.MinitestRow, time.Time) {
	s.Lock()
	defer s.Unlock()
	return s.rows, s.updatedAt
}

func (s *ReportStorage) Get() ([]model.ReportRow, time.Time) {
	s.Lock()
	defer s.Unlock()
	return s.rows, s.updatedAt
}

func (s *ClassEnqStorage) Get() ([]model.ClassEnqRow, time.Time) {
	s.Lock()
	defer s.Unlock()
	return s.rows, s.updatedAt
}
