package main
import ("api/session"
"github.com/gorilla/mux"
"github.com/rs/cors"
"database/sql"
"net/http"
"api/user"
"api/file"
"api/task"
"api/general"
"os"
"log")

func main() {

	f, err := os.OpenFile("logfile", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("Server launched")


	var s session.Session;

	db, err := sql.Open("mysql", "root:1@/coursework")
	if err != nil {
		panic(err)
	}
	u := &user.User{}
	s.Setup(db, u)
	task.Setup(db, u, &s)
	File.Setup(u, &s)
	general.Setup(&s, db)
	r := mux.NewRouter()
	c := cors.New(cors.Options{
        AllowedOrigins: []string{"*"},
        AllowCredentials: true,
    })

    handler := c.Handler(r)	

	r.HandleFunc("/auth", s.Auth)
	r.HandleFunc("/admin/auth", s.AdminAuth)
	r.HandleFunc("/reg", s.Reg).Methods("POST")
	r.HandleFunc("/redirect", s.RedirectIfLogged)
	r.HandleFunc("/admin_redirect", s.AdminRedirect)
	r.HandleFunc("/logout", s.Logout).Methods("POST")
	r.HandleFunc("/upload", File.Upload)
	r.HandleFunc("/task/create", task.Create)
	r.HandleFunc("/task/delete", task.Delete)
	r.HandleFunc("/task/get", task.GetById)
	r.HandleFunc("/task/get_by_student", task.GetByStudentId)
	r.HandleFunc("/task/get_for_student", task.GetForStudent)
	r.HandleFunc("/task/check_if_assigned", task.CheckIfAssigned)
	r.HandleFunc("/task/get_one", task.GetById)
	r.HandleFunc("/task/assign", task.Assign)
	r.HandleFunc("/task/all", task.Get)
	r.HandleFunc("/students", general.Students)
	r.PathPrefix("/audio/").
    Handler(http.StripPrefix("/audio/", http.FileServer(http.Dir("./audio"))))

	log.Fatal(http.ListenAndServe(":8080", handler))
}