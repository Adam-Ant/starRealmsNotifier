package main

import (
  "encoding/json"
  "errors"
  "fmt"
  "io/ioutil"
  "net/http"
  "net/url"
)

type Pass struct {
  Token2 string `json:"token2"`
}

func getToken(user string, password string) (string, error) {

  resp, err := http.PostForm("https://srprodv2.whitewizardgames.com/Account/Login", url.Values{
    "username": {user},
    "password": {password}})

  if err != nil {
    return "", err
  }

  defer resp.Body.Close()

  jsonResp, err := ioutil.ReadAll(resp.Body)

  if !IsJSON(jsonResp) {
    return "", errors.New(fmt.Sprintf("Invalid response from server:\n%s", jsonResp))
  }

  var structJson Pass
  err = json.Unmarshal(jsonResp, &structJson)
  if err != nil {
    return "", err
  }

  return structJson.Token2, nil
}
