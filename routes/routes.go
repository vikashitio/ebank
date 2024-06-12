package routes

import (
	"ebank/handlers"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var tpl *template.Template
var store = sessions.NewCookieStore([]byte("EindiaBusiness"))

func init() {
	//templates = template.Must(template.ParseGlob("templates/*.html"))
	tpl, _ = template.ParseGlob("static/*.html")
}

func InitRoutes() *mux.Router {

	router := mux.NewRouter()
	router.HandleFunc("/users", handlers.UsersHandler).Methods("GET")
	router.HandleFunc("/login", loginHandler)
	router.HandleFunc("/logout", logoutHandler)
	router.HandleFunc("/registration", registrationHandler)
	router.HandleFunc("/loginPost", handlers.UsersLogin).Methods("POST")
	router.HandleFunc("/registrationPost", handlers.UsersRegistration).Methods("POST")

	//http.Handle("/vv", http.RedirectHandler("http://www.google.com", 302))
	// Add more routes here

	// Serve static files
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	return router
}

type logindata struct {
	Title   string
	Email   string
	Contact string
	Message string
}

var logindata1 logindata

// loginHandler serves form for users to login with
func loginHandler(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "merchant")
	//fmt.Println(session)

	if session.Values["Login-Message"] == "" {
		var msg = ""
		tpl.ExecuteTemplate(w, "login.html", msg)
	} else {
		var msg = session.Values["Login-Message"]
		session.Values["Login-Message"] = ""
		session.Save(r, w)
		tpl.ExecuteTemplate(w, "login.html", msg)
	}

}

// loginHandler serves form for users to login with
func logoutHandler(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "merchant") // Get all Session
	session.Options.MaxAge = -1            //destroy all session
	// Save it before we write to the response/return from the handler.
	session.Save(r, w)

	//fmt.Println("session:", session)
	var msg = "Logout Successfully"
	tpl.ExecuteTemplate(w, "login.html", msg)

}

// Registration serves
func registrationHandler(w http.ResponseWriter, r *http.Request) {

	logindata1 = logindata{
		Title:   "Login Form",
		Email:   "vikashg@itio.in",
		Contact: "+ 977 9852 5862 55",
		Message: "",
	}

	tpl.ExecuteTemplate(w, "registration.html", logindata1)
}
