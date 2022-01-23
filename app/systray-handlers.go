package app

// systray 周り

import (
	"fmt"
	"os"
	"strconv"
	"time"

	calendar "github.com/earlgray283/gakujo-google-calendar/app/google-calendar-api"
	"github.com/earlgray283/gakujo-google-calendar/assets"
	"github.com/getlantern/systray"
	"github.com/go-co-op/gocron"
	"github.com/skratchdot/open-golang/open"
)

const dateTimeNotSubmitted = "0001-01-01 00:00:00 +0000 UTC" //未提出の課題を選別するためのconst

const googleCalendarURL = "https://calendar.google.com/calendar/" //GoogleCalendarのURL(string参照)
const gakujoURL = "https://gakujo.shizuoka.ac.jp/portal/"

func (a *App) OnReady() {
	a.Log.SetOutput(os.Stdout)
	a.crawler.Log.SetOutput(os.Stdout)
	//タスクトレイ大元の設定
	systray.SetIcon(assets.IconGakujo)
	systray.SetTitle("")
	systray.SetTooltip("Set tasks to your google calendar automatically.")

	// 直近の課題
	systray.AddMenuItem("Recent task: ", "Recent task")
	a.recentTaskItem = systray.AddMenuItem("Recent task: ", "task within 1 day")
	a.startRecentTaskUpdater()

	// 未提出課題数
	a.unSubmittedItem = systray.AddMenuItem("UnSubmitted task(s): ", "your task(s) to do")

	systray.AddSeparator()

	//GoogleCalendar, Gakujo を開く
	calendarOpener := systray.AddMenuItem("Open Google Calendar in Web Site", "Open Google Calendar")
	gakujoOpener := systray.AddMenuItem("Open Gakujo Portal Site", "Open Gakujo")

	systray.AddSeparator()

	// GoogleCalenderに登録するためのボタンの作成
	allAdder := systray.AddMenuItem("Add every task to calendar", "update your google calendar")

	// 最終同期
	a.lastSyncItem = systray.AddMenuItem("Last sync: ", "last updated at.")
	systray.AddSeparator()

	// quitボタンの作成
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	// 定期実行する
	if err := a.autoAddSchedule(); err != nil {
		a.Log.Println(err)
		a.Log.Println("定期実行に失敗しました")
		systray.Quit()
	}

	for {
		select {
		case <-calendarOpener.ClickedCh:
			if err := a.openWebSite(googleCalendarURL); err != nil {
				a.Log.Println(err)
				systray.Quit()
			}
			a.Log.Println("Google Calendarを開きます。")

		case <-gakujoOpener.ClickedCh:
			if err := a.openWebSite(gakujoURL); err != nil {
				a.Log.Println(err)
				systray.Quit()
			}
			a.Log.Println("学情ポータルサイトを開きます。")

		case <-allAdder.ClickedCh:
			count, err := a.registAll()
			if err != nil {
				a.Log.Println(err)
				a.Log.Println("タスクの登録に失敗しました")
				systray.Quit()
			}
			if count != 0 {
				a.Log.Println("全ての課題を登録しました。")
			} else {
				a.Log.Println("登録する課題はありませんでした。")
			}
		case <-mQuit.ClickedCh:
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
	return open.Run(url)
}

func (a *App) startRecentTaskUpdater() {
	s := gocron.NewScheduler(time.Local)

	_, _ = s.Every(time.Hour).Do(func() {
		now := time.Now()
		newTitle, deadline := "", now.AddDate(1, 0, 0)
		classenq := a.crawler.Classenq.GetMinByTime()
		fmt.Println(classenq.Title)
		if deadline.After(classenq.EndDate) {
			newTitle, deadline = "["+classenq.CourseName+"]"+classenq.Title, classenq.EndDate
		}
		report := a.crawler.Report.GetMinByTime()
		fmt.Println(report.Title)
		if deadline.After(report.EndDate) {
			newTitle, deadline = "["+report.CourseName+"]"+report.Title, report.EndDate
		}
		minitest := a.crawler.Minitest.GetMinByTime()
		fmt.Println(minitest.Title)
		if deadline.After(minitest.EndDate) {
			newTitle, deadline = "["+minitest.CourseName+"]"+minitest.Title, minitest.EndDate
		}
		a.recentTaskItem.SetTitle(newTitle)
		if deadline.Sub(now) < time.Hour*24 {
			a.recentTaskItem.SetIcon(assets.IconAlert)
		} else {
			a.recentTaskItem.SetIcon(nil)
		}
	})

	s.StartAsync()
}

func (a *App) registReport() (int, error) {
	calendarId := "primary"
	reportRows, _ := a.crawler.Report.Get()
	counter := 0
	for _, row := range reportRows {
		event := calendar.NewGakujoEvent("["+row.CourseName+"]"+row.Title, row.EndDate)
		if row.EndDate.After(time.Now()) {
			if row.LastSubmitDate.String() == dateTimeNotSubmitted {
				if err := calendar.AddSchedule(event, calendarId, a.srv); err != nil {
					return -1, err
				}
				counter += 1
			}
		}
	}
	return counter, nil
}

func (a *App) registMinitest() (int, error) {
	calendarId := "primary"
	counter := 0
	minitestRows, _ := a.crawler.Minitest.Get()
	for _, row := range minitestRows {
		event := calendar.NewGakujoEvent("["+row.CourseName+"]"+row.Title, row.EndDate)
		if row.SubmitStatus == "未提出" {
			if time.Now().Before(row.EndDate) {
				if err := calendar.AddSchedule(event, calendarId, a.srv); err != nil {
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
		event := calendar.NewGakujoEvent("["+row.CourseName+"]"+row.Title, row.EndDate)
		if row.SubmitStatus == "未提出" {
			if time.Now().Before(row.EndDate) {
				if err := calendar.AddSchedule(event, calendarId, a.srv); err != nil {
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

	a.updateItems()

	return cntSum, nil
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

func (a *App) updateItems() {
	cnt := strconv.Itoa(a.countUnSubmitted())
	timeNow := time.Now().Format("2006-01-02 15:04:05")
	a.unSubmittedItem.SetTitle("UnSubmitted task(s): " + cnt)
	a.unSubmittedItem.SetTooltip("you have " + cnt + "task(s) to do.")
	a.lastSyncItem.SetTitle("Last sync: " + timeNow)
	a.lastSyncItem.SetTooltip("your google calendar was updated at " + timeNow)
}

func (a *App) countUnSubmitted() int {
	cnt := 0
	reportRows, _ := a.crawler.Report.Get()
	minitestRows, _ := a.crawler.Minitest.Get()
	classEnqRows, _ := a.crawler.Classenq.Get()

	for _, row := range reportRows {
		if row.EndDate.After(time.Now()) {
			if row.LastSubmitDate.String() == dateTimeNotSubmitted {
				cnt++
			}
		}
	}
	for _, row := range minitestRows {
		if row.SubmitStatus == "未提出" {
			if time.Now().Before(row.EndDate) {
				cnt++
			}
		}
	}
	for _, row := range classEnqRows {
		if row.SubmitStatus == "未提出" {
			if time.Now().Before(row.EndDate) {
				cnt++
			}
		}
	}
	return cnt
}
