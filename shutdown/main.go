package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
}

func (app *App) Start() {
	fmt.Printf("App started.\n")
	time.Sleep(10 * time.Second)
}

func (app *App) Stop() {
	fmt.Printf("App stopped.\n")
}

func main() {
	app := &App{}
	exitCh := make(chan os.Signal)
	signal.Notify(exitCh,
		syscall.SIGTERM, // terminate: stopped by `kill -9 PID`
		syscall.SIGINT,  // interrupt: stopped by Ctrl + C
	)

	go func() {
		defer func() {
			exitCh <- syscall.SIGTERM // send terminate signal when application stop naturally
		}()
		app.Start() // start the application
	}()

	<-exitCh   // blocking until receive exit signal
	app.Stop() // stop the application
}
