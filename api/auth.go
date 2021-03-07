package main

import (
"database/sql"
_ "github.com/go-sql-driver/mysql"
"api/user"
"net/http"
"encoding/json"
"github.com/gorilla/mux"
"github.com/gorilla/sessions"
"os"
)

var u *user.User
var db *sql.DB
var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

type errmsg struct {
	Msg int `json:"status"`
}

func Reg(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var err errmsg
	loadFromReq(r)
	if (u.Reg() != nil) {
		err.Msg = -1
	} else {
		u.Auth()
		startSession(w, r)
	}
	json.NewEncoder(w).Encode(err)
}

func Auth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	loadFromReq(r)
	 var err errmsg
	if (u.Email == "" || u.Pwd == "") {
		err.Msg = -1
	} else {
		auth_success, auth_err := u.Auth()
		if (auth_err != nil) {
			err.Msg = -1
		} else {
			if (auth_success) {
				err.Msg = 0
				startSession(w, r)
			} else {
				err.Msg = 1
			}
		}
	}
	json.NewEncoder(w).Encode(err)
}

func loadFromReq(r *http.Request) {
	params := r.URL.Query()
	u.Email = params.Get("email")
	u.Pwd = params.Get("pwd")
	u.Name = params.Get("name")
}

func startSession(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	session.Values["id"] = u.GetId()
	session.Save(r, w)
}

func main() {
	r := mux.NewRouter()
	u = &user.User{}

	db, err := sql.Open("mysql", "root:1@/coursework")
	if err != nil {
		panic(err)
	}
	u.SetDB(db)
	
	r.HandleFunc("/auth", Auth).Methods("GET")
	r.HandleFunc("/reg", Reg)
	http.ListenAndServe(":8080", r)
}