package general

import (
	"net/http"
	"database/sql"
	"api/session"
	"log"
	"encoding/json"
)

var s *session.Session;
var db *sql.DB;

type student struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

type response struct {
	Errcode int `json:"status"`
	Payload []student `json:"payload"`
}

func Setup(_s *session.Session, _db *sql.DB) {
	s = _s
	db = _db
}

func Students(w http.ResponseWriter, r *http.Request) {
	s.SetHeaders(w ,r)
	var resp response;

	res, err := db.Query("SELECT id, name FROM users")
	if (err != nil) {
		resp.Errcode = -1
		log.Println("Couldn't retrieve student ids: ", err)
	} else {
		for res.Next() {
			var stu student;
			res.Scan(&stu.Id, &stu.Name)
			resp.Payload = append(resp.Payload, stu)
		}
	}

	json.NewEncoder(w).Encode(resp)
}