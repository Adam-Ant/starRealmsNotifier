package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)




func UnmarshalJSON(b []byte) ([]ActiveGames, []int, error){
	type Activetemp struct {
		Activegames []ActiveGames `json:"activegames"`
	}

	type Finishedtemp struct {
		Finishedgames []struct {
			Gameid int `json:"gameid"`
		} `json:"finishedgames"`
	}

	activetmp := Activetemp{}

	err := json.Unmarshal(b, &activetmp)
	if err != nil {
		return []ActiveGames{}, nil, err
	}

	active := []ActiveGames{}

	// Break the un-needed double loop
	active = activetmp.Activegames

	for _, game := range active {
		// All incoming games haven't been notified yet. Passed through to main function, notifier sets this to true.
		game.Hasnotified = false
	}

	finished := Finishedtemp{}

	err = json.Unmarshal(b, &finished)
	if err != nil {
		return []ActiveGames{}, nil, err
	}

	var finishedInt []int

	for _, game := range finished.Finishedgames {
		finishedInt = append(finishedInt, game.Gameid)
	}

	return active, finishedInt, nil
}

func getGames (token string) ([]ActiveGames, []int, error){
	client := &http.Client{Timeout: time.Second * 5}
	req, err := http.NewRequest("GET", "https://srprodv2.whitewizardgames.com/NewGame/ListActivitySortable", nil)
	if err != nil {
		return []ActiveGames{}, nil, err
	}
	req.Header.Add("coreversion", "28")
	req.Header.Add("Auth", token)
	resp, err := client.Do(req)

	if err != nil {
		return []ActiveGames{}, nil, err
	}

	defer resp.Body.Close()

	jsonResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []ActiveGames{}, nil, err
	}

	if !IsJSON(jsonResp) {
		return []ActiveGames{}, nil, errors.New(fmt.Sprintf("Invalid response from server:\n%s", jsonResp))
	}

	return UnmarshalJSON(jsonResp)
}