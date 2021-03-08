package session

import (
"database/sql"
_ "github.com/go-sql-driver/mysql"
"api/user"
"net/http"
"encoding/json"
"github.com/gorilla/sessions"
"os"
)

type Session struct {
	u *user.User
	db *sql.DB
	store *sessions.CookieStore
}

type errmsg struct {
	Msg int `json:"status"`
}

func (s *Session) Reg(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var err errmsg
	s.loadFromReq(r)
	if (s.u.Reg() != nil) {
		err.Msg = -1
	} else {
		s.u.Auth()
		s.startSession(w, r)
	}
	json.NewEncoder(w).Encode(err)
}

func (s *Session) Auth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	s.loadFromReq(r)
	 var err errmsg
	if (s.u.Email == "" || s.u.Pwd == "") {
		err.Msg = -1
	} else {
		auth_success, auth_err := s.u.Auth()
		if (auth_err != nil) {
			err.Msg = -1
		} else {
			if (auth_success) {
				err.Msg = 0
				s.startSession(w, r)
			} else {
				err.Msg = 1
			}
		}
	}
	json.NewEncoder(w).Encode(err)
}

func (s *Session) loadFromReq(r *http.Request) {
	params := r.URL.Query()
	s.u.Email = params.Get("email")
	s.u.Pwd = params.Get("pwd")
	s.u.Name = params.Get("name")
}

func (s *Session) startSession(w http.ResponseWriter, r *http.Request) {
	session, _ := s.store.Get(r, "session")
	session.Values["id"] = s.u.GetId()
	session.Save(r, w)
}

func (s *Session) Setup(db *sql.DB, u *user.User) {
	s.store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	s.u = u
	s.u.SetDB(db)
}