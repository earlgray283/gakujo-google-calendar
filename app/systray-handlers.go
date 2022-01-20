package app

// systray 周り

import (
	"os"
	"time"

	"github.com/getlantern/systray"
	"github.com/go-co-op/gocron"
	"github.com/skratchdot/open-golang/open"
	"google.golang.org/api/calendar/v3"
)

const dateTimeNotSubmited = "0001-01-01 00:00:00 +0000 UTC"//未提出の課題を選別するめのconst
const formatAdded = "2006-01-02 15:04.05"
const googleCalendarURL = "https://calendar.google.com/calendar/"//GoogleCalendarのURL(string参照)

func (a *App) OnReady() {
	a.Log.SetOutput(os.Stdout)
	iconData, err := os.ReadFile("app/icon.ico")
	if err != nil {
		a.Log.Println(err)
		systray.Quit()
	}
	//タスクトレイ大元の設定
	systray.SetIcon(iconData)
	systray.SetTitle("Google Calender Tasktray")
	systray.SetTooltip("Set tasks to your google calendar automaticaly")

	//GoogleCalendarをWebSiteで開くためのボタン
	WebSiteOpener := systray.AddMenuItem("Open Google Calendar in Web Site", "Open Google Calendar")
	systray.AddSeparator()

	// GoogleCalenderに登録するためのボタンの作成
	AllAdder := systray.AddMenuItem("Add every task to calendar", "ADD all")
	ReportAdder := systray.AddMenuItem("Add report to calendar", "ADD report")
	MinitestAdder := systray.AddMenuItem("Add minitest to calendar", "ADD minitest")
	ClassEnqAdder := systray.AddMenuItem("Add classenq to calendar", "ADD classenq")

	// 登録と退出の境界線を作る
	systray.AddSeparator()

	// quitボタンの作成
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	// 定期実行する
	err = a.autoAddSchedule()
	if err != nil {
		systray.Quit()
	}

	for {
		select {
		case <-WebSiteOpener.ClickedCh:
			err := a.openWebSite(googleCalendarURL)
			if err != nil {
				systray.Quit()
			}
			a.Log.Println("Google Calendarを開きます。")
		case <-AllAdder.ClickedCh:
			count, err := a.registAll()
			if err != nil {
				systray.Quit()
			}

			if count != 0 {
				a.Log.Println("全ての課題を登録しました。")
			} else {
				a.Log.Println("登録する課題はありませんでした。")
			}

		case <-ReportAdder.ClickedCh:
			count, err := a.registReport()
			if err != nil {
				systray.Quit()
			}

			if count != 0 {
				a.Log.Println("レポート課題を登録しました。")
			} else {
				a.Log.Println("登録するレポート課題はありませんでした。")
			}

		case <-MinitestAdder.ClickedCh:
			count, err := a.registMinitest()
			if err != nil {
				systray.Quit()
			}

			if count != 0 {
				a.Log.Println("小テスト課題を登録しました。")
			} else {
				a.Log.Println("登録する小テスト課題はありませんでした。")
			}

		case <-ClassEnqAdder.ClickedCh:
			count, err := a.registClassEnq()
			if err != nil {
				systray.Quit()
			}

			if count != 0 {
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

func (a *App) openWebSite(url string) error {
	err := open.Run(url)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) registReport() (int, error) {
	calendarId := "primary"
	reportRows, _ := a.crawler.Report.Get()
	counter := 0
	for _, row := range reportRows {
		event := newEvent("["+row.CourseName+"]"+row.Title, row.EndDate)

		// 未提出だったら・・・・・・・・・・・
		if row.LastSubmitDate.String() == dateTimeNotSubmited {
			err := AddSchedule(event, calendarId, a.srv)
			if err != nil {
				return -1, err
			}
			counter += 1
		}
	}
	return counter, nil
}

func (a *App) registMinitest() (int, error) {
	calendarId := "primary"
	counter := 0
	minitestRows, _ := a.crawler.Minitest.Get()
	for _, row := range minitestRows {
		event := newEvent("["+row.CourseName+"]"+row.Title, row.EndDate)

		if row.SubmitStatus == "未提出" {
			if time.Now().Before(row.EndDate) {
				err := AddSchedule(event, calendarId, a.srv)
				if err != nil {
					return -1, err
				}
				counter += 1
			}
		}

	}
	return counter, nil
}

func (a *App) registClassEnq() (int, error) {
	counter := 0
	calendarId := "primary"
	classEnqRows, _ := a.crawler.Classenq.Get()
	for _, row := range classEnqRows {
		event := newEvent("["+row.CourseName+"]"+row.Title, row.EndDate)

		// 未提出だったら・・・・・・・・・・・
		if row.SubmitStatus == "未提出" {
			if time.Now().Before(row.EndDate) {
				err := AddSchedule(event, calendarId, a.srv)
				if err != nil {
					return -1, nil
				}
				counter += 1
			}
		}
	}
	return counter, nil
}

func (a *App) registAll() (int, error) {
	cntSum := 0

	ReportCounter, err := a.registReport()
	if err != nil {
		return 0, err
	}
	cntSum += ReportCounter

	MinitestCounter, err := a.registMinitest()
	if err != nil {
		return 0, err
	}
	cntSum += MinitestCounter

	ClassEnqCounter, err := a.registClassEnq()
	if err != nil {
		return 0, err
	}
	cntSum += ClassEnqCounter

	return cntSum, nil
}

func newEvent(title string, t time.Time) *calendar.Event { // タイトルと日時を入れると Event 型を返す
	added := time.Now().Format(formatAdded)
	Event := &calendar.Event{
		Summary:     title,
		Location:    "学務情報システム",
		Description: "学情カレンダーから追加された予定です。\n学務情報システム: https://gakujo.shizuoka.ac.jp/portal/\nAdded: " + added,
		Start: &calendar.EventDateTime{
			DateTime: t.Add(-time.Hour).Format(formatAdded),
			TimeZone: "Asia/Tokyo",
		},
		End: &calendar.EventDateTime{
			DateTime: t.Format(formatAdded),
			TimeZone: "Asia/Tokyo",
		},
	}
	return Event
}

func (a *App) autoAddSchedule() error { // 定期実行
	s := gocron.NewScheduler(time.Local)

	_, _ = s.Every(30).Minutes().Do(func() {
		counter, err := a.registAll()
		if err != nil {
			systray.Quit()
		}

		if counter != 0 {
			a.Log.Println("すべての予定を登録しました。")
		} else {
			a.Log.Println("登録する予定はありませんでした。")
		}
	})
	s.StartAsync()
	return nil
}
