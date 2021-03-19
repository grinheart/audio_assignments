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
	"io/ioutil"
)

var db *sql.DB;
var u *user.User;
var s *session.Session;

type Task struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Body string `json:"body"`
	Audio []string `json:"audio"`
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

func setInternalServerError(w http.ResponseWriter, resp *response) {
	resp.Errcode = -1
	resp.Msg = "Внутренняя ошибка сервера"
	json.NewEncoder(w).Encode(resp)
}

func Create(w http.ResponseWriter, r *http.Request) {
	var resp response
	if (!s.IsAdmin(r)) {
		setInternalServerError(w, &resp)
		return;
	} 
	s.SetHeaders(w, r)
	params := params(r)
	resp.Errcode = 0
	stmt, err := db.Prepare("INSERT INTO tasks (title, body) VALUES (?, ?)")
	if (err != nil) {
		log.Println("Error preparing inserting tasks")
		setInternalServerError(w, &resp)
		return
	}
	res, err := stmt.Exec(params.Get("title"), params.Get("body"))
	if (err != nil) {
		log.Println("Error executing inserting tasks. Title:\n", params.Get("title"), "\n Body:\n", params.Get("body"))
		setInternalServerError(w, &resp)
		return
	}

	taskId64, err := res.LastInsertId()
	if (err != nil) {
		setInternalServerError(w, &resp)
		log.Println("Error retieveing task's inserted id. Title:\n", params.Get("title"), "\n Body:\n", params.Get("body"))
		return
	}

	taskId := int(taskId64)

	if (params.Get("assign_to") == "") {
		json.NewEncoder(w).Encode(resp)
		return
	}

	assignToStr := strings.Split(params.Get("assign_to"), ",")
	var assignTo []interface{}
	query := "INSERT INTO assignments (task_id, student_id) VALUES"
	for _, idStr := range assignToStr {
		id, err := strconv.Atoi(idStr)
		if (err != nil) {
			log.Println("Error parsing ids to assign tasks. Ids param: ", params.Get("assign_to"))
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
		log.Println("Error assigning tasks: ", err, "\nquery: " + query)
		setInternalServerError(w, &resp)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	var resp response
	if (!s.IsAdmin(r)) {
		setInternalServerError(w, &resp)
		return;
	} 
	id := params(r).Get("id")
	_, err := db.Query("DELETE FROM tasks WHERE id = ?", id)
	if (err != nil) {
		panic(err)
	}
}

func GetById(w http.ResponseWriter, r *http.Request) {
	s.SetHeaders(w, r)
	id, err := strconv.Atoi(params(r).Get("id"))
	var resp response
	if (err != nil) {
		log.Println("Parsing id err while retrieving a task for student:", err)
		setInternalServerError(w, &resp)
		return
	}
	res, err := db.Query("SELECT id, title, body FROM tasks WHERE id = ?", id)
	if (err != nil) {
		log.Println("Error retrieving a task:", err)
		setInternalServerError(w, &resp)
		return
	} 
	if (!res.Next()) {
		log.Println("Error retrieving a task - no such id: ", id)
		setInternalServerError(w, &resp)
		return
	}
	resp.Payload = append(resp.Payload, Task{Id: id})
	res.Scan(&resp.Payload[0].Id, &resp.Payload[0].Title, &resp.Payload[0].Body)
	json.NewEncoder(w).Encode(resp)
}

func getTasks(w http.ResponseWriter, r *http.Request, query string, queryArgs []interface{}, errmsg string) (response) {
	s.SetHeaders(w, r)
	log.Printf("args")
	log.Println(queryArgs...)
	res, err := db.Query(query, queryArgs...)
	var resp response
	if (err != nil) {
		log.Println(errmsg, "; error:", err)
		setInternalServerError(w, &resp)
		return resp
	}
	resp.Errcode = errorCode(&res, err)
	if (resp.Errcode == 0) {
		for i := 0;;i++ {
			resp.Payload = append(resp.Payload, Task{})
			res.Scan(&resp.Payload[i].Id, &resp.Payload[i].Title, &resp.Payload[i].Body)
			resp.Payload[i].Audio = []string{}
			if (!res.Next()) {
				break;
			}
		}
	}
	return resp
}

func getByStudentId(w http.ResponseWriter, r *http.Request, id int) {
	resp := getTasks(
		w, r,
		"SELECT tasks.id, tasks.title, tasks.body FROM tasks, assignments WHERE tasks.id = assignments.task_id AND assignments.student_id = ?",
		[]interface{}{id},
		"Couldn't retrieve student's tasks")
	for i, task := range resp.Payload {
		id := strconv.Itoa(id)
		task_id := strconv.Itoa(task.Id)
		path := "audio/" + id + "/" + task_id
		files, err := ioutil.ReadDir("./" + path)
		if err != nil {
			log.Println("couldn't open audio dir ", path)
		}
		for _, f := range files {
            resp.Payload[i].Audio = append(resp.Payload[i].Audio, path + "/" + f.Name())
    	}
	}
	json.NewEncoder(w).Encode(resp)
}

func GetByStudentId(w http.ResponseWriter, r *http.Request) {
	var resp response
	if (!s.IsAdmin(r)) {
		setInternalServerError(w, &resp)
		return;
	} 
	id, _ := strconv.Atoi(params(r).Get("id"))
	getByStudentId(w, r, id)
	
}

func GetForStudent(w http.ResponseWriter, r *http.Request) {
	s.SetHeaders(w, r)
	getByStudentId(w, r, u.GetId())
}

func Get(w http.ResponseWriter, r *http.Request) {
	var resp response
	if (!s.IsAdmin(r)) {
		setInternalServerError(w, &resp)
		return;
	} 
	resp = getTasks(
		w, r,
		"SELECT tasks.id, tasks.title, tasks.body FROM tasks",
		[]interface{}{},
		"Couldn't retrieve tasks")

	json.NewEncoder(w).Encode(resp)
}

func CheckIfAssigned(w http.ResponseWriter, r *http.Request) {
	var resp response
	s.SetHeaders(w, r)
	params := params(r)
	task_id, err := strconv.Atoi(params.Get("task_id"))
	if (err != nil) {
		log.Println("Invalid task_id: ", params.Get("task_id"))
		setInternalServerError(w, &resp)
		return
	}
	student_id := s.SavedId(r)
	if (student_id == 0) {
		resp.Errcode = 1
		json.NewEncoder(w).Encode(resp)
		return
	}
	res, err := db.Query("SELECT id FROM assignments WHERE task_id = ? AND student_id = ?", task_id, student_id)
	if (!res.Next()) {
		resp.Errcode = 1
	}
	json.NewEncoder(w).Encode(resp)
}

func Assign(w http.ResponseWriter, r *http.Request) {
	if (!s.IsAdmin(r)) {
		return;
	} 
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