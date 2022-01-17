package crawle

import "time"

type CrawleOption struct {
	ReportInterval   time.Duration
	MinitestInterval time.Duration
	ClassenqInterval time.Duration
	RetryCount       int
}

func DefaultCrawleOption() *CrawleOption {
	return &CrawleOption{
		ReportInterval:   3 * time.Hour,
		MinitestInterval: 3 * time.Hour,
		ClassenqInterval: 3 * time.Hour,
		RetryCount:       5,
	}
}
