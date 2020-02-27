package pages

import (
	"fmt"
	"html/template"
	"net/http"
)

func IndexPage(w http.ResponseWriter, r *http.Request) {
	data := IndexData{
		PageTitle: "Bjorn is cool",
		Name:      "Bjorn",
	}
	tmp, e := template.ParseFiles("./templates/index.html")
	if e != nil {
		fmt.Println(e.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tmp.Execute(w, data)
}

func BackupFile(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	b := make([]byte, 100)
	body.Read(b)
	w.WriteHeader(http.StatusOK)
}
