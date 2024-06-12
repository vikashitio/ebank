package handlers

import (
	"ebank/models"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("EindiaBusiness"))

var tpl *template.Template

func init() {
	//templates = template.Must(template.ParseGlob("templates/*.html"))
	tpl, _ = template.ParseGlob("static/*.html")
}

type Sub struct {
	Username string
	Data     string
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := models.GetUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
	//fmt.Println(users)

}

func UsersLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()                      // Parses the request body
	username := r.Form.Get("username") // Fetch Value of username
	password := r.Form.Get("password") // Fetch Value of password

	message, _ := models.GetLogeddetails(username, password)

	//fmt.Println(message)
	session, err := store.Get(r, "merchant")
	if message.Alert == "" {
		//fmt.Println("Client ID:", message.ID)
		//fmt.Println("Name:", message.Name)
		//fmt.Println("Email:", message.Email)
		//fmt.Println("Status:", message.Status)
		// Store Session Variable

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set some session values.
		session.Values["Merchant-Name"] = message.Name
		session.Values["Merchant-ID"] = message.ID
		session.Values["Merchant-Email"] = message.Email
		session.Values["Merchant-Status"] = message.Status
		//session.Values["Login-Message"] = "Login Done"

		fmt.Println("Client ID 11:", session.Values["Merchant-ID"])
		fmt.Println("Name 22:", session.Values["Merchant-Name"])
		fmt.Println("Email 33:", session.Values["Merchant-Email"])
		fmt.Println("Status 44:", session.Values["Merchant-Status"])

		// Save it before we write to the response/return from the handler.
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", 302)
	} else {
		session.Values["Login-Message"] = message.Alert
		// Save it before we write to the response/return from the handler.
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//fmt.Println("Status 44:", session.Values["Login-Message"])
		http.Redirect(w, r, "/login", 302)

	}

}

func UsersRegistration(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()                // Parses the request body
	name := r.Form.Get("name")   // Fetch Value of name
	email := r.Form.Get("email") // Fetch Value of email

	message, err := models.UsersRegistration(name, email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//fmt.Println(message)

	fmt.Println("Client ID:", message.ID)
	fmt.Println("Name:", message.Name)
	fmt.Println("Email:", message.Email)
	fmt.Println("Status:", message.Status)

	//http.Redirect(w, r, "/", 302)
}
