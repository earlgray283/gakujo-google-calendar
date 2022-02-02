package crawle

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("./../../.env"); err != nil {
		if len(os.Getenv("GAKUJO_USERNAME")) == 0 || len(os.Getenv("GAKUJO_PASSWORD")) == 0 {
			log.Fatal(err)
		}
	}
}

func TestCrawle(t *testing.T) {
	username, password := os.Getenv("GAKUJO_USERNAME"), os.Getenv("GAKUJO_PASSWORD")
	opt := DefaultCrawleOption()
	opt.ReportInterval = time.Minute
	opt.MinitestInterval = time.Minute
	opt.ClassenqInterval = time.Minute
	crawler, err := NewCrawler(username, password, opt)
	if err != nil {
		t.Fatal(err)
	}
	crawler.Log.SetOutput(os.Stdout)
	errc := crawler.Start()

	rm := map[time.Time]struct{}{}
	mm := map[time.Time]struct{}{}
	cm := map[time.Time]struct{}{}
	ticker := time.NewTicker(time.Minute)
	timerC := time.After(3 * time.Minute)
	for {
		select {
		case err := <-errc:
			t.Fatal(err)
		case <-ticker.C:
			log.Println(rm, mm, cm)
		case <-timerC:
			goto L1
		}
		rm[crawler.Report.UpdatedAt()] = struct{}{}
		mm[crawler.Minitest.UpdatedAt()] = struct{}{}
		cm[crawler.Classenq.UpdatedAt()] = struct{}{}
	}
L1:

	if len(rm) < 2 {
		t.Fatal("report rows is not enough")
	}
	if len(mm) < 2 {
		t.Fatal("minitest rows is not enough")
	}
	if len(cm) < 2 {
		t.Fatal("classenq rows is not enough")
	}

	t.Log(len(rm), len(mm), len(cm))
}
