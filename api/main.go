package main
import ("api/session"
"github.com/gorilla/mux"
"database/sql"
"net/http"
"api/user"
"api/file"
"api/task")

func main() {

	var s session.Session;

	db, err := sql.Open("mysql", "root:1@/coursework")
	if err != nil {
		panic(err)
	}
	u := &user.User{}
	s.Setup(db, u)
	task.Setup(db, u)
	File.Setup(u)
	r := mux.NewRouter()
	r.HandleFunc("/auth", s.Auth)
	r.HandleFunc("/reg", s.Reg)
	r.HandleFunc("/upload", File.Upload)
	r.HandleFunc("/task/create", task.Create)
	r.HandleFunc("/task/delete", task.Delete)
	r.HandleFunc("/task/get", task.GetById)
	r.HandleFunc("/task/get_by_student", task.GetByStudentId)
	r.HandleFunc("/task/get_for_student", task.GetForStudent)
	r.HandleFunc("/task/assign", task.Assign)

	http.ListenAndServe(":8080", r)
}