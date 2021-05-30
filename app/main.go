package main

import (
  "os"
  "fmt"
  "log"
  "net/http"
  "encoding/json"
  "strconv"
  "strings"

  "database/sql"
  _ "github.com/go-sql-driver/mysql"
)

func formHandler(w http.ResponseWriter, r *http.Request) {
  if err := r.ParseForm(); err != nil {
    fmt.Fprintf(w, "ParseForm() err: %v", err)
    return
  }
  fmt.Fprintf(w, "POST request succcessful")
  name := r.FormValue("name")
  address := r.FormValue("address")

  fmt.Fprintf(w, "Name = %s\n", name)
  fmt.Fprintf(w, "Address = %s\n", address)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/hello" {
    http.Error(w, "404 not found", http.StatusNotFound)
    return
  }

  if r.Method != "GET" {
    http.Error(w, "Method is not supported.", http.StatusNotFound)
    return
  }

  fmt.Fprintf(w, "Hello!")
}

type User struct {
  Id int `json:"id"`
  Username string `json:"username"`
  Age int `json:"age"`
}

func dbConn() (db *sql.DB) {
  //todo get from env_variables
  dbDriver := os.Getenv("DB_DRIVER")
  dbUser := os.Getenv("DB_USER")
  dbPass := os.Getenv("DB_PASSWORD")
  dbHost := os.Getenv("DB_PROTO") + "(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")"
  dbName := os.Getenv("DB_NAME")
  db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@"+dbHost+"/"+dbName)
  if err != nil {
    panic(err.Error())
  }
  return db
}

func Index(w http.ResponseWriter, r *http.Request) {
  db := dbConn()
  defer db.Close()
  selDB, err := db.Query("SELECT * FROM users ORDER BY id DESC")
  if err != nil {
    panic(err.Error())
  }

  user := User{}
  res := []User{}
  for selDB.Next() {
    var id, age int
    var username string
    err = selDB.Scan(&id, &username, &age)
    if err != nil {
      panic(err.Error())
    }

    user.Id = id
    user.Username = username
    user.Age = age
    res = append(res, user)
  }
  s, err := json.Marshal(res)
  if err != nil {
    panic(err)
  }
  w.Header().Set("Content-Type", "application/json")
  w.Write(s)
}
func UserHandler(w http.ResponseWriter, r *http.Request){
  db := dbConn()
  defer db.Close()
  nId := strings.TrimPrefix(r.URL.Path, "/users/")
  var user User
  switch r.Method {
    case http.MethodGet:
      selDB, err := db.Query("SELECT * FROM users WHERE id=?",nId)
      if err != nil {
        panic(err.Error())
      }
      for selDB.Next() {
        var id, age int
        var username string
        err = selDB.Scan(&id, &username, &age)
        if err != nil {
            panic(err.Error())
        }
        user.Id = id
        user.Username = username
        user.Age = age
      }
      s, err := json.Marshal(user)
      if err != nil {
        panic(err)
      }
      w.Header().Set("Content-Type", "application/json")
      w.Write(s)
    case http.MethodPost:
      err := json.NewDecoder(r.Body).Decode(&user)
      insForm, err := db.Prepare("UPDATE users SET username=?, age=? WHERE id=?")
      if err != nil {
          panic(err.Error())
      }
      insForm.Exec(user.Username, user.Age, nId)
      log.Println("UPDATE: Username: " + user.Username + " | Age: " + strconv.Itoa(user.Age) + " | Id: " + nId)
      w.Header().Add("Content-Type", "application/json")
      w.WriteHeader(http.StatusNoContent)
    case http.MethodDelete:
      insForm, err := db.Prepare("DELETE FROM users WHERE id=?")
      if err != nil {
          panic(err.Error())
      }
      insForm.Exec(nId)
      log.Println("DELETE: Id: " + nId)
      w.Header().Add("Content-Type", "application/json")
      w.WriteHeader(http.StatusNoContent)
  }
}

func Store(w http.ResponseWriter, r *http.Request) {
  db := dbConn()
  defer db.Close()
  if r.Method == "POST" {
    var user User
    err := json.NewDecoder(r.Body).Decode(&user)
    insForm, err := db.Prepare("INSERT INTO users(username, age) VALUES(?,?)")
    if err != nil {
        panic(err.Error())
    }
    insForm.Exec(user.Username, user.Age)
    log.Println("INSERT: UserName: " + user.Username + " | Age: " + strconv.Itoa(user.Age))
  }
  w.Header().Add("Content-Type", "application/json")
  w.WriteHeader(http.StatusCreated)
}

func main() {
  // todo more specific error handling
  // todo package db access methods
  fileServer := http.FileServer(http.Dir("./static"))
  http.Handle("/", fileServer)
  http.HandleFunc("/form", formHandler)
  http.HandleFunc("/hello", helloHandler)
  http.HandleFunc("/users", Index)
  http.HandleFunc("/users/", UserHandler)
  http.HandleFunc("/store", Store)

  fmt.Printf("Starting server at port 8080\n")
  if err := http.ListenAndServe(":8080", nil); err != nil {
      log.Fatal(err)
  }
}
