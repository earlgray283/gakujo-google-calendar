package main

import (
	"fmt"
	"io/ioutil"
	"time"
	"os"
	"log"

	"github.com/getlantern/systray"
	//"github.com/getlantern/systray/example/icon"
	//"github.com/skratchdot/open-golang/open"
)

//iconのデータを読み込む
func iconDetaLoder () []byte {
	b, err := os.ReadFile("icon/iconwin.ico")
		if err != nil {
  			log.Fatal(err)
		}
	return b
}

func main() {
	onExit := func() {
		now := time.Now()
		ioutil.WriteFile(fmt.Sprintf(`on_exit_%d.txt`, now.UnixNano()), []byte(now.String()), 0644)
	}

	systray.Run(onReady, onExit)
}

func onReady() {

	// We can manipulate the systray in other goroutines
	go func() {
		iconData := iconDetaLoder()
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
				//ScrapingCode Here

			//終了のボタンの動作実行
			case <-mQuit.ClickedCh:
				systray.Quit()
				fmt.Println("Quit2 now...")
				return
			}
		}
	}()
}

