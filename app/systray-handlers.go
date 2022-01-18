package app

// systray 周り

import (
	"log"
	"os"
	"time"

	"github.com/getlantern/systray"
	"google.golang.org/api/calendar/v3"
)

func (a *App) OnReady() {
	a.Log.SetOutput(os.Stdout)
	iconData, err := os.ReadFile("app/icon.ico")
	if err != nil {
		a.Log.Println(err)
		systray.Quit()
	}
	systray.SetIcon(iconData)
	systray.SetTitle("Google Calender Tasktray")
	systray.SetTooltip("Set tasks to your google calendar automaticaly")

	// GoogleCalenderに登録するためのボタンの作成
	AllAdder := systray.AddMenuItem("Add every task to calendar", "ADD all")
	ReportAdder := systray.AddMenuItem("Add report to calendar", "ADD report")
	MinitestAdder := systray.AddMenuItem("Add minitest to calendar", "ADD minitest")
	ClassEnqAdder := systray.AddMenuItem("Add classenq to calendar", "ADD classenq")

	// 登録と退出の境界線を作る
	systray.AddSeparator()

	// quitボタンの作成
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	for {
		select {
		case <-AllAdder.ClickedCh:
			CounterSum := 0

			ReportCounter, err := a.registerReport()
			if err != nil {
				log.Fatal(err)
			}
			CounterSum += ReportCounter

			MinitestCounter, err := a.registerMinitest()
			if err != nil {
				log.Fatal(err)
			}
			CounterSum += MinitestCounter

			ClassEnqCounter, err := a.registerClassEnq()
			if err != nil {
				log.Fatal(err)
			}
			CounterSum += ClassEnqCounter

			if CounterSum != 0 {
				a.Log.Println("全ての課題を登録しました。")
			} else {
				a.Log.Println("登録する課題はありませんでした。")
			}

		case <-ReportAdder.ClickedCh:
			counter, err := a.registerReport()
			if err != nil {
				log.Fatal(err)
			}

			if counter != 0 {
				a.Log.Println("レポート課題を登録しました。")
			} else {
				a.Log.Println("登録するレポート課題はありませんでした。")
			}

		case <-MinitestAdder.ClickedCh:
			counter, err := a.registerMinitest()
			if err != nil {
				log.Fatal(err)
			}

			if counter != 0 {
				a.Log.Println("小テスト課題を登録しました。")
			} else {
				a.Log.Println("登録する小テスト課題はありませんでした。")
			}

		case <-ClassEnqAdder.ClickedCh:
			counter, err := a.registerClassEnq()
			if err != nil {
				log.Fatal(err)
			}

			if counter != 0 {
				a.Log.Println("授業アンケートを登録しました。")
			} else {
				a.Log.Println("登録する授業アンケートはありませんでした。")
			}

		case <-mQuit.ClickedCh: // 終了のボタンの動作実行
			a.Log.Println("タスクトレイアプリを終了します。")
			systray.Quit()
			return
		}
	}
}

func (a *App) OnExit() {
	a.appLogFile.Close()
	a.crawlerLogFile.Close()
}

func (a *App) registerReport() (int, error) {
	calendarId := "primary"
	reportRows, _ := a.crawler.Report.Get()
	counter := 0
	for _, row := range reportRows {
		Event := &calendar.Event{
			Summary:     "[" + row.CourseName + "]" + row.Title,
			Location:    "学務情報システム",
			Description: "TaskID: " + row.TaskMetadata.ID,
			Start: &calendar.EventDateTime{
				DateTime: row.EndDate.Add(-time.Hour).Format("2006-01-02T15:04:05+09:00"),
				TimeZone: "Asia/Tokyo",
			},
			End: &calendar.EventDateTime{
				DateTime: row.EndDate.Format("2006-01-02T15:04:05+09:00"),
				TimeZone: "Asia/Tokyo",
			},
		}

		// 未提出だったら・・・・・・・・・・・
		if row.LastSubmitDate.String() == "0001-01-01 00:00:00 +0000 UTC" {
			err := AddSchedule(Event, calendarId, a.srv)
			if err != nil {
				return -1, err
			} else {
				counter += 1
			}
		}
	}
	return counter, nil
}

func (a *App) registerMinitest() (int, error) {
	calendarId := "primary"
	counter := 0

	minitestRows, _ := a.crawler.Minitest.Get()
	for _, row := range minitestRows {
		Event := &calendar.Event{
			Summary:     "[" + row.CourseName + "]" + row.Title,
			Location:    "学務情報システム",
			Description: "TaskID: " + row.TaskMetadata.ID,
			Start: &calendar.EventDateTime{
				DateTime: row.EndDate.Add(-time.Hour).Format("2006-01-02T15:04:05+09:00"),
				TimeZone: "Asia/Tokyo",
			},
			End: &calendar.EventDateTime{
				DateTime: row.EndDate.Format("2006-01-02T15:04:05+09:00"),
				TimeZone: "Asia/Tokyo",
			},
		}

		if row.SubmitStatus == "未提出" {
			if time.Now().Before(row.EndDate) {
				err := AddSchedule(Event, calendarId, a.srv)
				if err != nil {
					return -1, err
				} else {
					counter += 1
				}
			}
		}

	}
	return counter, nil
}

func (a *App) registerClassEnq() (int, error) {
	counter := 0
	calendarId := "primary"

	classEnqRows, _ := a.crawler.Classenq.Get()
	for _, row := range classEnqRows {
		Event := &calendar.Event{
			Summary:     "[" + row.CourseName + "]" + row.Title,
			Location:    "学務情報システム",
			Description: "TaskID: " + row.TaskMetadata.ID,
			Start: &calendar.EventDateTime{
				DateTime: row.EndDate.Add(-time.Hour).Format("2006-01-02T15:04:05+09:00"),
				TimeZone: "Asia/Tokyo",
			},
			End: &calendar.EventDateTime{
				DateTime: row.EndDate.Format("2006-01-02T15:04:05+09:00"),
				TimeZone: "Asia/Tokyo",
			},
		}

		// 未提出だったら・・・・・・・・・・・
		if row.SubmitStatus == "未提出" {
			if time.Now().Before(row.EndDate) {
				err := AddSchedule(Event, calendarId, a.srv)
				if err != nil {
					return -1, nil
				} else {
					counter += 1
				}
			}
		}
	}
	return counter, nil
}
