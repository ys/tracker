package main

import (
  "io/ioutil"
  "github.com/mrjones/oauth"
  "log"
  "os"
  "encoding/json"
)

type BodyWeight struct {
  Measures []WeightMeasure `json:"body-weight"`
}

type WeightMeasure struct {
  DateTime string
  Value float64 `json:",string"`
}

type FitbitUser struct {
  DisplayName string
  FullName string
  Nickname string
  Avatar string
  Avatar150 string
}

type UserResult struct {
  User FitbitUser
}

func FitbitUrl() string {
  _, url, err := fitbitClient().GetRequestTokenAndUrl("http://localhost:8080/auth/callback")
  if err != nil {
    log.Fatal(err)
  }
  return url
}

func GetAccessToken(oauthToken string, oauthVerifier string) *oauth.AccessToken {
  accessToken, err := fitbitClient().AuthorizeToken(&oauth.RequestToken{Token: oauthToken}, oauthVerifier)
  if err != nil {
    log.Fatal(err)
  }
  return accessToken
}

func LastMonthWeight(accessToken *oauth.AccessToken) BodyWeight {
  response, err := fitbitClient().Get(
    "http://api.fitbit.com/1/user/-/body/weight/date/today/1m.json",
    map[string]string{},
    accessToken)
  if err != nil {
    log.Fatal(err)
  }
  defer response.Body.Close()

  bits, err := ioutil.ReadAll(response.Body)
  var bodyWeight BodyWeight
  json.Unmarshal(bits, &bodyWeight)
  return bodyWeight
}

func UserProfile(accessToken *oauth.AccessToken) FitbitUser {
  response, err := fitbitClient().Get(
    "http://api.fitbit.com/1/user/-/profile.json",
    map[string]string{},
    accessToken)
  if err != nil {
    log.Fatal(err)
  }
  defer response.Body.Close()

  bits, err := ioutil.ReadAll(response.Body)
  var user UserResult
  json.Unmarshal(bits, &user)
  return user.User
}

func fitbitClient() *oauth.Consumer{
  consumerKey := os.Getenv("FITBIT_CONSUMER_KEY")
  consumerSecret := os.Getenv("FITBIT_CONSUMER_SECRET")
  return oauth.NewConsumer(
    consumerKey,
    consumerSecret,
    oauth.ServiceProvider{
      RequestTokenUrl:   "http://api.fitbit.com/oauth/request_token",
      AuthorizeTokenUrl: "http://www.fitbit.com/oauth/authorize",
      AccessTokenUrl:    "http://api.fitbit.com/oauth/access_token",
    })
}
