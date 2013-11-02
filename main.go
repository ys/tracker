package main

import (
  "encoding/json"
  "net/http"
  "os"
  _ "github.com/lib/pq"
  "database/sql"

  "github.com/mrjones/oauth"
  "github.com/gorilla/mux"
  "github.com/joho/godotenv"
)

func db() *sql.DB{
  db, err := sql.Open("postgres", "user=ys dbname=tracker sslmode=disable")
  if err != nil {
    panic("")
  }
  return db
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
  var debug bool = false
  w.Header().Set("Location", FitbitUrl())
  w.WriteHeader(302)
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
  params := r.URL.Query()
  accessToken := GetAccessToken(params.Get("oauth_token"), params.Get("oauth_verifier"))
  user := UserProfile(accessToken)
  insertFirstUser(user, accessToken)
}

func insertFirstUser(user FitbitUser, accessToken *oauth.AccessToken) int {
  var id string
  db().QueryRow("SELECT id FROM users LIMIT 1").Scan(&id)
  if id != nil {
    return 0
  }
  stmt, err := db().Prepare("INSERT INTO users(username, access_token, secret_token) VALUES($1, $2, $3)")
  if err != nil {
    panic(err)
  }
  _, err = stmt.Exec(user.Nickname, accessToken.Token, accessToken.Secret)
  if err != nil {
    panic(err)
  }
  return 1
}

func accessToken() *oauth.AccessToken {
  var access_token string
  var secret_token string
  db().QueryRow("SELECT access_token, secret_token FROM users WHERE id=1").Scan(&access_token, &secret_token)
  return &oauth.AccessToken{Token: access_token, Secret: secret_token}
}

func WeightHandler(w http.ResponseWriter, r *http.Request) {
  bodyWeight := LastMonthWeight(accessToken())
  w.Header().Set("Content-Type", "application/json")
  jsonBytes, _ := json.Marshal(bodyWeight.Measures)
  w.Write(jsonBytes)
}

func main() {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }
  r := mux.NewRouter()
  r.HandleFunc("/auth", AuthHandler)
  r.HandleFunc("/auth/callback", CallbackHandler)
  r.HandleFunc("/weight", WeightHandler)
  http.Handle("/", r)
  http.ListenAndServe(":8080", r)
}

