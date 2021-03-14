package File

import (
	"net/http"
	"io/ioutil"
	"strconv"
	"api/session"
	"api/user"
	"encoding/json"
	"os"
	"log"
)

var s *session.Session
var u *user.User


type response struct {
	Errcode int `json:"status"`
	Msg string `json:"message"`
}


func setInternalServerError(w http.ResponseWriter, resp *response) {
	resp.Errcode = -1
	resp.Msg = "Внутренняя ошибка сервера"
	json.NewEncoder(w).Encode(resp)
}

func Setup(_u *user.User, _s *session.Session) {
	s = _s
	u = _u
}

func Upload(w http.ResponseWriter, r *http.Request) {
	s.SetHeaders(w, r)
	r.ParseMultipartForm(10 << 20)
	var resp response;
	id := r.FormValue("id")
	if (id == "") {
		log.Println("Id is empty")
		setInternalServerError(w, &resp);
		return
	}
	dir := "./audio/" + strconv.Itoa(u.GetId()) + "/" + id
	os.Mkdir(dir, 0755) //todo: check and log if dir exists
	file, _, err := r.FormFile("audio")
	if (err != nil) {
		log.Println("Couldn't form audio")
		setInternalServerError(w, &resp);
		return
	}
	defer file.Close()
	tempFile, err := ioutil.TempFile(dir, "*.wav")
	if (err != nil) {
		log.Println("Couldn't save audio")
		setInternalServerError(w, &resp);
		return
	}
	defer tempFile.Close()
	fileBytes, err := ioutil.ReadAll(file)
    if err != nil {
		log.Println("Couldn't read file")
		setInternalServerError(w, &resp);
		return
    }
    // write this byte array to our temporary file
    tempFile.Write(fileBytes)
}

func retrieve() {

}