package crawle

import (
	"sync"
	"time"

	"github.com/earlgray283/gakujo-google-calendar/gakujo/model"
)

type MinitestStorage struct {
	sync.Mutex
	Rows []model.MinitestRow
	UpdatedAt    time.Time
}

type ReportStorage struct {
	sync.Mutex
	Rows []model.ReportRow
	UpdatedAt  time.Time
}

type ClassEnqStorage struct {
	sync.Mutex
	Rows []model.ClassEnqRow
	UpdatedAt    time.Time
}
