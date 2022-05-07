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

func (s *MinitestStorage) UpdatedAt() time.Time {
	s.Lock()
	defer s.Unlock()
	return s.updatedAt
}

func (s *MinitestStorage) Get() []model.MinitestRow {
	s.Lock()
	defer s.Unlock()
	return s.rows
}

func (s *MinitestStorage) GetWithFilter(f func(*model.MinitestRow) bool) []model.MinitestRow {
	s.Lock()
	defer s.Unlock()
	var rows []model.MinitestRow
	for _, row := range s.rows {
		if f(&row) {
			rows = append(rows, row)
		}
	}
	return rows
}

func (s *MinitestStorage) GetMinByTime() *model.MinitestRow {
	s.Lock()
	defer s.Unlock()
	var min model.MinitestRow
	now := time.Now()
	min.EndDate = time.Now().AddDate(1, 0, 0)
	for _, row := range s.rows {
		if now.After(row.EndDate) {
			continue
		}
		if min.EndDate.After(row.EndDate) && row.SubmitStatus == model.UnSubmited {
			min = row
		}
	}
	return &min
}

func (s *ReportStorage) UpdatedAt() time.Time {
	s.Lock()
	defer s.Unlock()
	return s.updatedAt
}

func (s *ReportStorage) Get() []model.ReportRow {
	s.Lock()
	defer s.Unlock()
	return s.rows
}

func (s *ReportStorage) GetWithFilter(f func(*model.ReportRow) bool) []model.ReportRow {
	s.Lock()
	defer s.Unlock()
	var rows []model.ReportRow
	for _, row := range s.rows {
		if f(&row) {
			rows = append(rows, row)
		}
	}
	return rows
}

func (s *ReportStorage) GetMinByTime() *model.ReportRow {
	s.Lock()
	defer s.Unlock()
	var min model.ReportRow
	now := time.Now()
	min.EndDate = time.Now().AddDate(1, 0, 0)
	for _, row := range s.rows {
		if now.After(row.EndDate) {
			continue
		}
		emptyTime := time.Time{}
		if min.EndDate.After(row.EndDate) && row.LastSubmitDate.String() == emptyTime.String() {
			min = row
		}
	}
	return &min
}

func (s *ClassEnqStorage) UpdatedAt() time.Time {
	s.Lock()
	defer s.Unlock()
	return s.updatedAt
}

func (s *ClassEnqStorage) Get() []model.ClassEnqRow {
	s.Lock()
	defer s.Unlock()
	return s.rows
}

func (s *ClassEnqStorage) GetWithFilter(f func(*model.ClassEnqRow) bool) []model.ClassEnqRow {
	s.Lock()
	defer s.Unlock()
	var rows []model.ClassEnqRow
	for _, row := range s.rows {
		if f(&row) {
			rows = append(rows, row)
		}
	}
	return rows
}

func (s *ClassEnqStorage) GetMinByTime() *model.ClassEnqRow {
	s.Lock()
	defer s.Unlock()
	var min model.ClassEnqRow
	now := time.Now()
	min.EndDate = time.Now().AddDate(1, 0, 0)
	for _, row := range s.rows {
		if now.After(row.EndDate) {
			continue
		}
		if min.EndDate.After(row.EndDate) && row.SubmitStatus == model.UnSubmited {
			min = row
		}
	}
	return &min
}
