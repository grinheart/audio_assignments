package main
import ("api/session"
"github.com/gorilla/mux"
"database/sql"
"net/http")

func main() {

	var s session.Session;

	db, err := sql.Open("mysql", "root:1@/coursework")
	if err != nil {
		panic(err)
	}
	s.Setup(db)
	r := mux.NewRouter()
	r.HandleFunc("/auth", s.Auth)
	r.HandleFunc("/reg", s.Reg)
	http.ListenAndServe(":8080", r)
}