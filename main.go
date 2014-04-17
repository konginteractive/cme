package main

import (
	C "github.com/konginteractive/cme/controler"
	M "github.com/konginteractive/cme/model"
	//_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	//"github.com/jinzhu/gorm"
	"html/template"
	"log"
	"net/http"
)

var templates *template.Template

func init() {
	templates = template.Must(template.ParseGlob("./vues/*"))
}

func main() {

	r := mux.NewRouter()

	// listes des rootes
	r.HandleFunc("/", HomeHandler)

	r.HandleFunc("/forum", ForumHandler)
	r.HandleFunc("/forum/nouveau", ForumAddHandler)
	r.HandleFunc("/forum/{category}", ForumCatHandler)

	r.HandleFunc("/eleves", StudentHandler)

	r.HandleFunc("/tutoriels", TutoHandler)
	r.HandleFunc("/tutoriels/nouveau/", TutoAddHandler)

	r.HandleFunc("/actualites", NewsHandler)

	//gestion des fichiers statiques
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

	http.Handle("/", r)
	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)

}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	Render(w, "index", C.HomeView())
}

func ForumHandler(w http.ResponseWriter, r *http.Request) {
	Render(w, C.ForumTempl, C.ForumView())
}

func ForumAddHandler(w http.ResponseWriter, r *http.Request) {
	Render(w, C.ForumAddTempl, C.ForumAddView())
}
func ForumCatHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	category := vars["category"]
	Render(w, C.ForumTempl, C.FormViewCategory(category))
}

func StudentHandler(w http.ResponseWriter, r *http.Request) {
	Render(w, C.UserTempl, C.UserView())
}

func TutoHandler(w http.ResponseWriter, r *http.Request) {
	Render(w, C.TutorialTempl, C.TutorialView())
}

func TutoAddHandler(w http.ResponseWriter, r *http.Request) {
	Render(w, "tutoriel_add", C.HomeView())
}

func NewsHandler(w http.ResponseWriter, r *http.Request) {
	Render(w, C.NewsTempl, C.NewsView())
}

func Render(w http.ResponseWriter, tmpl string, p M.Page) {
	w.Header().Add("Accept-Charset", "utf-8")
	w.Header().Add("Content-Type", "text/html")
	err := templates.ExecuteTemplate(w, tmpl, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
