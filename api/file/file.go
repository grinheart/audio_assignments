package File

import (
	"net/http"
	"io/ioutil"
	"strconv"
)

func upload(w http.ResponseWriter, r *http.Request, id int) error {
	r.ParseMultipartForm(10 << 20)
	dir := "audio/" + strconv.Itoa(id)
	file, _, err := r.FormFile("audio")
	if (err != nil) {
		return err
	}
	defer file.Close()
	tempFile, err := ioutil.TempFile(dir, "*.wav")
	defer tempFile.Close()
	return err
}

func retrieve() {
	
}