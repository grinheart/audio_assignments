package main
import ("api/session"
"github.com/gorilla/mux"
"database/sql"
"net/http"
"api/user"
"api/file")

func main() {

	var s session.Session;

	db, err := sql.Open("mysql", "root:1@/coursework")
	if err != nil {
		panic(err)
	}
	u := &user.User{}
	s.Setup(db, u)
	File.Setup(u)
	r := mux.NewRouter()
	r.HandleFunc("/auth", s.Auth)
	r.HandleFunc("/reg", s.Reg)
	r.HandleFunc("/upload", File.Upload)
	http.ListenAndServe(":8080", r)
}