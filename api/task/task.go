package task

import (
	"net/http"
	"net/url"
	"database/sql"
	"api/user"
	"encoding/json"
	"strconv"
)

var db *sql.DB;
var u *user.User;

type Task struct {
	Title string
	Body string
}

type response struct {
	Msg int `json:"status"`
	Payload []Task `json:"payload"`
}

func Setup(_db *sql.DB, _u *user.User) {
	db = _db
	u = _u
}

//todo: check db set

func params(r *http.Request) url.Values {
	return r.URL.Query()
} 

func errorCode(res **sql.Rows, err error) int {
	if (err != nil) {
		return -1
	}
	if (!(*res).Next()) {
		return 1
	}
	return 0
}

func Create(w http.ResponseWriter, r *http.Request) {
	params := params(r)
	stmt, err := db.Prepare("INSERT INTO tasks (title, body) VALUES (?, ?)")
	if (err != nil) {
		panic(err)
	}
	_, err = stmt.Exec(params.Get("title"), params.Get("body"))
	if (err != nil) {
		panic(err)
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	id := params(r).Get("id")
	_, err := db.Query("DELETE FROM tasks WHERE id = ?", id)
	if (err != nil) {
		panic(err)
	}
}

func GetById(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(params(r).Get("id"))
	res, err := db.Query("SELECT title, body FROM tasks WHERE id = ?", id)
	if (err != nil) {
		panic(err)
	} 
	var resp response
	resp.Msg = errorCode(&res, err)
	if (resp.Msg == 0) {
		resp.Payload = append(resp.Payload, Task{})
		res.Scan(&resp.Payload[0].Title, &resp.Payload[0].Body)
	}
	json.NewEncoder(w).Encode(resp)
}

func getByStudentId(w http.ResponseWriter, r *http.Request, id int) {
	res, err := db.Query("SELECT tasks.title, tasks.body FROM tasks, assignments WHERE tasks.id = assignments.task_id AND assignments.student_id = ?", id)
	if (err != nil) {
		panic(err)
	}
	var resp response
	resp.Msg = errorCode(&res, err)
	if (resp.Msg == 0) {
		for i := 0;;i++ {
			resp.Payload = append(resp.Payload, Task{})
			res.Scan(&resp.Payload[i].Title, &resp.Payload[i].Body)
			if (!res.Next()) {
				break;
			}
		}
	}
	json.NewEncoder(w).Encode(resp)
}

func GetByStudentId(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(params(r).Get("id"))
	getByStudentId(w, r, id)
	
}

func GetForStudent(w http.ResponseWriter, r *http.Request) {
	getByStudentId(w, r, u.GetId())
}

func Assign(w http.ResponseWriter, r *http.Request) {
	params := params(r)
	task_id := params.Get("task_id")
	student_id := params.Get("student_id")
	stmt, err := db.Prepare("INSERT INTO assignments (task_id, student_id) VALUES (?, ?)")
	if (err != nil) {
		panic(err)
	}
	_, err = stmt.Exec(task_id, student_id)
	if (err != nil) {
		panic(err)
	}
}