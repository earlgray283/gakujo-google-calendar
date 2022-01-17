package app

// systray 周り

import (
	"os"

	"github.com/getlantern/systray"
)

func (a *App) OnReady() {
	iconData, err := os.ReadFile("app/icon.ico")
	if err != nil {
		a.Log.Println(err)
		systray.Quit()
	}
	systray.SetIcon(iconData)
	systray.SetTitle("Google Calender Tasktray")
	systray.SetTooltip("Set tasks to your google calendar automaticaly")

	// GoogleCalenderに登録するためのボタンの作成
	mGCAdder := systray.AddMenuItem("Add to calendar", "ADD to Google Calendar.")

	// 登録と退出の境界線を作る
	systray.AddSeparator()

	// quitボタンの作成
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	for {
		select {
		case <-mGCAdder.ClickedCh: // 学情の情報をGoogleCalenderに登録するためのボタンの動作実行
			a.Log.Println("タスクを登録しました。")
		case <-mQuit.ClickedCh: // 終了のボタンの動作実行
			a.Log.Println("タスクトレイアプリを終了します。")
			systray.Quit()
			return
		}
	}
}
