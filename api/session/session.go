package session

import (
"database/sql"
_ "github.com/go-sql-driver/mysql"
"api/user"
"net/http"
"encoding/json"
"github.com/gorilla/sessions"
"os"
"log"
)

type Session struct {
	u *user.User
	db *sql.DB
	store *sessions.CookieStore
}

type errmsg struct {
	Status int `json:"status"`
	Msg string `json:"message"`
}

func (s *Session) setInternalServerError(w http.ResponseWriter, e *errmsg) {
	e.Status = -1
	e.Msg = "Внутренняя ошибка сервера"
	json.NewEncoder(w).Encode(e)
}

func (s *Session) SetHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Set-Cookie", w.Header().Get("Set-Cookie") + "; SameSite=Strict")
	log.Println("SavedId", s.SavedId(r))
	_, err := s.u.Load(s.SavedId(r))
	log.Println("error loading", err)
}

func (s *Session) user(w http.ResponseWriter, r *http.Request, msg string, status int) {
	var err errmsg
	err.Status = status
	if (status != 0) {
		err.Msg = msg
	} else {
		s.startSession(w, r)
	}
	s.SetHeaders(w, r)
	json.NewEncoder(w).Encode(err)
}

func (s *Session) RedirectIfLogged(w http.ResponseWriter, r *http.Request) {
	s.SetHeaders(w, r)
	var resp errmsg
	if (s.SavedId(r) != 0) {
		resp.Status = 0
	} else {
		resp.Status = 1
	}
	s.SetHeaders(w, r)
	json.NewEncoder(w).Encode(resp)
}

func (s *Session) Reg(w http.ResponseWriter, r *http.Request) {
	s.loadFromReq(r)
	s.user(w, r, "Пользователь уже зарегистрирован", s.u.Reg())
}

func (s *Session) Auth(w http.ResponseWriter, r *http.Request) {
	s.loadFromReq(r)
	s.user(w, r, "Неверный логин или пароль", s.u.Auth())
}

func (s *Session) AdminAuth(w http.ResponseWriter, r *http.Request) {
	var err errmsg
	r.ParseForm()
	res, e := s.db.Query("SELECT * FROM admins WHERE email = ? AND pwd = ?", r.Form.Get("email"), r.Form.Get("pwd"))
	if (e != nil) {
		log.Println("Error checking admin: ", e)
		s.setInternalServerError(w, &err)
	}
	if (res.Next()) {
		log.Println("Start admin session")
		s.startAdminSession(w, r)
	} else {
		err.Status = 1
		err.Msg = "Неверный логин или пароль"
	}
	s.SetHeaders(w, r)
	json.NewEncoder(w).Encode(err)
}

func (s *Session) loadFromReq(r *http.Request) {
	r.ParseForm()
	s.u.Email = r.Form.Get("email")
	s.u.Pwd = r.Form.Get("pwd")
	s.u.Name = r.Form.Get("name")
	log.Println(s.u.Email, " ", s.u.Pwd)
}

func (s *Session) sesId(r *http.Request) (*sessions.Session, error) {
	log.Println(s.store.Get(r, "session"))
	return s.store.Get(r, "session")
}

func (s *Session) AdminRedirect(w http.ResponseWriter, r *http.Request) {
	var err errmsg
	log.Println("AdminRedirect")
	if (!s.IsAdmin(r)) {
		s.setInternalServerError(w, &err)
		return;
	}
	json.NewEncoder(w).Encode(err)
}

func (s *Session) SavedId(r *http.Request) (int) {
	session, _ := s.sesId(r)
	id, ok := session.Values["id"].(int)
	if (!ok) {
		log.Println("not ok")
		return 0
	}
	log.Println("saved id ", id)
	return id
}

func (s *Session) IsAdmin(r *http.Request) (bool) {
	session, _ := s.sesId(r)
	isAdmin, ok := session.Values["admin"].(int)
	log.Println("isAdmin", isAdmin)
	if (!ok) {
		return false
	}
	return isAdmin > 0
}

func (s *Session) Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Set-Cookie", "session=")
	s.SetHeaders(w, r)
}

func (s *Session) startSession(w http.ResponseWriter, r *http.Request) {
	session, _ := s.sesId(r)
	session.Values["id"] = s.u.GetId()
	err := session.Save(r, w)
	if (err != nil) {
		log.Println("saving session failed: ", err)
	}
	s.SetHeaders(w, r)
}

func (s *Session) startAdminSession(w http.ResponseWriter, r *http.Request) {
	session, _ := s.store.Get(r, "session")
	session.Values["admin"] = 1
	err := session.Save(r, w)
	if (err != nil) {
		log.Println("saving session failed: ", err)
	}
	s.SetHeaders(w, r)
}

func (s *Session) Setup(db *sql.DB, u *user.User) {
	s.store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	s.u = u
	s.u.SetDB(db)
	s.db = db
}