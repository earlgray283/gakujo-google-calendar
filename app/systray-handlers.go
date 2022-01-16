package app

// systray 周り

import (
	"fmt"
	"os"
	"log"

	"github.com/getlantern/systray"
)

//iconのデータを読み込む
func readIconData (name string) ([]byte, error) {
	b, err := os.ReadFile(name)
	if err != nil {
  		return b,err
	}
	return b,nil
}

func main() {
	onExit := func() {
	}

	systray.Run(onReady, onExit)
}

func onReady() {

	// We can manipulate the systray in other goroutines
	go func() {
		iconData,err := readIconData("iconwin.ico")
		if err != nil{
			log.Fatal(err)
		}
		systray.SetTemplateIcon(iconData, iconData)
		systray.SetTitle("Google Calender TaskTray")
		systray.SetTooltip("Open the menu")
		
		//GoogleCalenderに登録するためのボタンの作成
		mGCAdder := systray.AddMenuItem("Add to calendar", "ADD to Google Calendar.")

		mGCAdder.SetIcon(iconData)

		//登録と退出の境界線を作る
		systray.AddSeparator()

		//退出ボタンの作成
		mQuit := systray.AddMenuItem("退出", "Quit the whole app")

		for {
			select {
			//学情の情報をGoogleCalenderに登録するためのボタンの動作実行
			case <-mGCAdder.ClickedCh:
				fmt.Println("example::タスクを登録しました。")
				//ScrapingCode Here

			//終了のボタンの動作実行
			case <-mQuit.ClickedCh:
				systray.Quit()
				fmt.Println("タスクトレイアプリを終了します。")
				return
			}
		}
	}()
}

