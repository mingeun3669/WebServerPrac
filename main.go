package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/http"

	libs "github.com/mingeun3669/L34rn_G0/WebServerPrac/Libs"
)

var uuid map[string]string = libs.Uuid{}

func index(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, nil)
}

func register(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/register.html")
	t.Execute(w, nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/login.html")
	t.Execute(w, nil)
}

// TODO : Implement uuid with Map
func process(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name, email, pwd := r.FormValue("name"), r.FormValue("email"),
		r.FormValue("pwd")

	fmt.Println(name, email, pwd)
	num, _ := libs.SendMail((email))
	_, err := libs.Uuid.Search(uuid, email)
	fmt.Println(num)
	if err != nil {
		t, _ := template.ParseFiles("templates/process.html")
		t.Execute(w, nil)
		libs.Uuid.Add(uuid, email, num)
		hash := sha256.New()
		hash.Write([]byte(email))
		md := hash.Sum(nil)
		mdStr := hex.EncodeToString(md)
		fmt.Println(mdStr)
		cookie := http.Cookie{
			Name:     "email",
			Value:    mdStr,
			HttpOnly: true,
		}
		w.Header().Set("Set-Cookie", cookie.String())
	} else {
		fmt.Fprintln(w, "Already Registered")
	}
}

func process_verify(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// vNum, cNum := strconv.Itoa(), r.FormValue("ver")
	// w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// if vNum == cNum {
	// 	fmt.Fprintln(w, `Register Success!
	// 	<a href="/">Please Back to main page.</a>`)
	// } else {
	// 	fmt.Fprintln(w, `Register Failed!
	// 	<a href="/">Please Back to main page.</a>`)
	// }
}

func main() {
	server := http.Server{
		Addr: "0.0.0.0:8080",
	}
	http.HandleFunc("/", index)
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/process", process)
	http.HandleFunc("/process_verify", process_verify)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	server.ListenAndServe()
}