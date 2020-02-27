package pages

import (
	"net/http"
)

func ServeScripts() {
	fs := http.FileServer(http.Dir("scripts/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
}
