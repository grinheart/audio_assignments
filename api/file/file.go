package File

import (
	"net/http"
	"io/ioutil"
	"strconv"
	"api/user"
	"os"
	"errors"
)

var u *user.User

func Setup(user *user.User) {
	u = user
}

func Upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	id := r.FormValue("id")
	if (id == "") {
		panic(errors.New("Id is empty"))
	}
	dir := "./audio/" + strconv.Itoa(u.GetId()) + "/" + id
	os.Mkdir(dir, 0755) //todo: check and log if dir exists
	file, _, err := r.FormFile("myFile")
	if (err != nil) {
		panic(err) 
	}
	defer file.Close()
	tempFile, err := ioutil.TempFile(dir, "*.wav")
	if (err != nil) {
		panic(err)
	}
	defer tempFile.Close()
}

func retrieve() {

}