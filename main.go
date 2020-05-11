package main

import (
	"flag"
	"fmt"
	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
	"github.com/rakyll/statik/fs"
	"io/ioutil"
	"log"
	_ "starRealmsNotify/statik"
	"time"
)

var (
	user      = flag.String("u", "", "Username")
	pass      = flag.String("p", "", "Password")
	tray      = flag.Bool("s", false, "Dont fork to tray (optional)")
	AuthToken string
)

func main() {
	flag.Parse()

	if *user == "" {
		fatalLogger(fmt.Errorf("Please specify a username with -u"))
	}

	if *pass == "" {
		fatalLogger(fmt.Errorf("Please specify a password with -p"))
	}

	var err error

	AuthToken, err = getToken(*user, *pass)
	if err != nil {
		fatalLogger(err)
	}

	if *tray {
		doCheck()
	} else {
		log.Println("Running in tray....")
		systray.Run(onReady, onExit)
	}
}

func onReady() {
	statikFS, err := fs.New()
	if err != nil {
		fatalLogger(err)
	}

	r, err := statikFS.Open("/StarRealms.ico")
	if err != nil {
		fatalLogger(err)
	}
	defer r.Close()


	contents, err := ioutil.ReadAll(r)
	if err != nil {
		fatalLogger(err)
	}
	systray.SetTitle("Star Realm Notifier")
	systray.SetTooltip(fmt.Sprintf("Logged in as %s", *user))
	systray.SetIcon(contents)
	mQuit := systray.AddMenuItem("Quit", "Quits this app")

	go doCheck()

	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

//TODO: Break out this to another file, & systray to another another file.
func doCheck() {
	var currentGames []ActiveGames
	for {
		active, finished, err := getGames(AuthToken)
		if err != nil {
			fatalLogger(err)
		}
		// For every game returned by the API, add it to the tracked list or update the list if its already present.
		for _, newGame := range active {
			gameFound := false
			for _, curGame := range currentGames {
				if newGame.Gameid == curGame.Gameid {
					gameFound = true
					if curGame.Lastupdatedtime.Before(newGame.Lastupdatedtime) {
						// There has been an update. Update our listing.
						log.Printf("Adding game %d to active list\n", curGame.Gameid)
						curGame = newGame
					}
					break
				}
			}
			if !gameFound {
				// Game not in our list, add it.
				log.Printf("Adding game %d to active list\n", newGame.Gameid)
				currentGames = append(currentGames, newGame)
			}
		}

		// If a game is in the finished list & our list, remove it.
		for _, finishedGame := range finished {
			for i, curGame := range currentGames {
				if curGame.Gameid == finishedGame {
					log.Printf("Removing game %d from active list\n", finishedGame)
					currentGames = append(currentGames[:i], currentGames[i+1:]...)
				}
			}
		}

		for i, game := range currentGames {
			// TODO: Make timeout variable
			if (game.Lastupdatedtime.Before(time.Now().Add(time.Second * -300))) && !game.Hasnotified {
				//TODO: Play sound/terminal bell? Option to disable?
				//TODO: Add insults :D
				err := beeep.Notify("Star Realms", fmt.Sprintf("It is your turn in a game with %s", game.Opponentname), "notify.png")
				if err != nil {
					fatalLogger(err)
				}
				log.Printf("Notified for turn in game %d", game.Gameid)
				currentGames[i].Hasnotified = true
			}
		}
		//TODO: Make recheck value customizable with lower limit
		time.Sleep(60 * time.Second)
	}
}

func onExit() {
	log.Println("Exiting...")
}
