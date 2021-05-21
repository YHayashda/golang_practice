package main

import (
  "fmt"
  "log"
  "net/http"
  "encoding/json"
  "strconv"

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
  dbDriver := "mysql"
  dbUser := "db_user"
  dbPass := "example"
  dbHost := "tcp(db:3306)"
  dbName := "sample"
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
func Show(w http.ResponseWriter, r *http.Request){
  db := dbConn()
  defer db.Close()
  // todo update to pathparameter, how to get pathparameter :id
  nId := r.URL.Query().Get("id")
  selDB, err := db.Query("SELECT * FROM users WHERE id=?",nId)
  if err != nil {
    panic(err.Error())
  }
  user := User{}
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
  // todo CRUD操作を書く
  // todo DBのconfigを環境変数から(ソースコードにベタガキはやめる)
  fileServer := http.FileServer(http.Dir("./static"))
  http.Handle("/", fileServer)
  http.HandleFunc("/form", formHandler)
  http.HandleFunc("/hello", helloHandler)
  http.HandleFunc("/users", Index)
  http.HandleFunc("/show", Show)
  http.HandleFunc("/store", Store)

  fmt.Printf("Starting server at port 8080\n")
  if err := http.ListenAndServe(":8080", nil); err != nil {
      log.Fatal(err)
  }
}
