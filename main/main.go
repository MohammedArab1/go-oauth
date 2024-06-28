package main

import (
	"fmt"
	// "net/http"

	// "html/template"
	"os"

	// "github.com/go-chi/chi/v5"
	// "github.com/go-chi/chi/v5/middleware"
)

func main() {

	fs := os.DirFS("static/")
	fmt.Println(fs.Open("index.html"))
	// someHtml := `
	// <h1>You've made it to the main page!</h1>
	// `

	// r := chi.NewRouter()
	// r.Use(middleware.Logger)
	// r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 	// http.FileServer(http.Dir("./static"))
	// 	tmpl := template.Must(template.New("myTemplate").Parse(someHtml))
	// 	tmpl.Execute(w, "some data")
	// })
	// http.ListenAndServe(":3000", r)
}
