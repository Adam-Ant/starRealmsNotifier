package main

import (
	"encoding/json"
	"github.com/gen2brain/dlgs"
	"log"
	"os"
)

func IsJSON(str []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(str, &js) == nil
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
