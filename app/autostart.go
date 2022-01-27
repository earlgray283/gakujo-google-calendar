package app

import (
	"os"

	autostart "github.com/emersion/go-autostart"
	"github.com/getlantern/systray"
)

type AutoStarter struct {
	item *systray.MenuItem
	app  *autostart.App
}

func NewAutoStartApp() *AutoStarter {
	execPath, _ := os.Executable()
	as := &AutoStarter{
		item: systray.AddMenuItem("", ""),
		app: &autostart.App{
			Name: "gakujo-google-calendar",
			Exec: []string{execPath},
		},
	}
	as.setTitle()
	return as
}

func (a *AutoStarter) StartAsync() chan error {
	errC := make(chan error)
	go func() {
		for {
			<-a.item.ClickedCh
			if err2 := a.handle(); err2 != nil {
				errC <- err2
				return
			}
		}
	}()
	return errC
}

func (a *AutoStarter) setTitle() {
	if a.app.IsEnabled() {
		a.item.SetTitle("自動起動をオフ")
		a.item.SetTooltip("autostart: enabled")
	} else {
		a.item.SetTitle("自動起動をオン")
		a.item.SetTooltip("autostart: disabled")
	}
}

func (a *AutoStarter) handle() error {
	defer a.setTitle()
	if a.app.IsEnabled() {
		return a.app.Disable()
	} else {
		return a.app.Enable()
	}
}
