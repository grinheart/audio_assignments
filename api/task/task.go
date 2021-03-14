package task

import (
	"net/http"
	"net/url"
	"database/sql"
	"api/user"
	"api/session"
	"encoding/json"
	"strconv"
	"log"
	"strings"
)

var db *sql.DB;
var u *user.User;
var s *session.Session;

type Task struct {
	Title string
	Body string
}

type response struct {
	Errcode int `json:"status"`
	Msg string `json:"message"`
	Payload []Task `json:"payload"`
}

func Setup(_db *sql.DB, _u *user.User, _s *session.Session) {
	db = _db
	u = _u
	s = _s
}

//todo: check db set

func params(r *http.Request) url.Values {
	r.ParseForm()
	return r.Form
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

func setInternalServerError(resp *response) {
	resp.Errcode = -1
	resp.Msg = "Внутренняя ошибка сервера"
}

func Create(w http.ResponseWriter, r *http.Request) {
	s.SetHeaders(w, r)
	var resp response
	params := params(r)
	resp.Errcode = 0
	stmt, err := db.Prepare("INSERT INTO tasks (title, body) VALUES (?, ?)")
	if (err != nil) {
		setInternalServerError(&resp)
		log.Println("Error preparing inserting tasks")
		json.NewEncoder(w).Encode(resp)
		return
	}
	res, err := stmt.Exec(params.Get("title"), params.Get("body"))
	if (err != nil) {
		setInternalServerError(&resp)
		log.Fatal("Error executing inserting tasks. Title:\n", params.Get("title"), "\n Body:\n", params.Get("body"))
		
		json.NewEncoder(w).Encode(resp)
		return
	}

	taskId64, err := res.LastInsertId()
	if (err != nil) {
		setInternalServerError(&resp)
		log.Fatal("Error retieveing task's inserted id. Title:\n", params.Get("title"), "\n Body:\n", params.Get("body"))
		json.NewEncoder(w).Encode(resp)
		return
	}

	taskId := int(taskId64)

	assignToStr := strings.Split(params.Get("assign_to"), ",")
	var assignTo []interface{}
	query := "INSERT INTO assignments (task_id, student_id) VALUES"
	for _, idStr := range assignToStr {
		id, err := strconv.Atoi(idStr)
		if (err != nil) {
			log.Fatal("Error parsing ids to assign tasks. Ids param: ", params.Get("assign_to"))
			json.NewEncoder(w).Encode(resp)
			return
		}
		assignTo = append(assignTo, taskId)
		assignTo = append(assignTo, id)
		query += " (?, ?),"
	}
	query = query[:len(query) - 1] + ";"

	//todo: eliminate double code
	_, err = db.Query(query, assignTo...)
	if (err != nil) {
		setInternalServerError(&resp)
		log.Fatal("Error assigning tasks: ", err, "\nquery: " + query)
		json.NewEncoder(w).Encode(resp)
		return
	}

	json.NewEncoder(w).Encode(resp)
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
	resp.Errcode = errorCode(&res, err)
	if (resp.Errcode == 0) {
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
	resp.Errcode = errorCode(&res, err)
	if (resp.Errcode == 0) {
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