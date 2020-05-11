package main

import (
	"github.com/gen2brain/dlgs"
	"log"
	"os"
	"time"
)

type ActiveGames struct {
	Gameid          int       `json:"gameid"`
	Opponentname    string    `json:"opponentname"`
	Actionneeded    bool      `json:"actionneeded"`
	Lastupdatedtime time.Time `json:"lastupdatedtime"`
	Hasnotified     bool
}

func fatalLogger(err error) {
	msg := err.Error()
	if *tray {
		log.Println(msg)
	} else {
		_, err := dlgs.Warning("Star Realms Notifier", msg)
		if err != nil {
			panic(err)
		}
	}
	os.Exit(1)
}
