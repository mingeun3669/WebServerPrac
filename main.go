package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/http"

	libs "github.com/mingeun3669/WebServerPrac/Libs"
)

var uuid map[string]string = libs.Uuid{}

/* TODO :
Implement deleting all cookies
Implement uuid check
*/

func deleteCookie(w http.ResponseWriter, name string) {
	c := http.Cookie{
		Name:     name,
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, &c)
}

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

	hash := sha256.New()
	hash.Write([]byte(email))
	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)
	fmt.Println(name, email, pwd)
	_, err := libs.Uuid.Search(uuid, mdStr)
	if err != nil {
		// set cookie
		fmt.Println(mdStr)
		cookie := http.Cookie{
			Name:     "___uuid",
			Value:    mdStr,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)

		num, _ := libs.SendMail((email))
		t, _ := template.ParseFiles("templates/process.html")
		t.Execute(w, nil)
		libs.Uuid.Add(uuid, mdStr, num)
	} else {
		fmt.Fprintln(w, "Already Registered")
	}
}
func process_verify(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	v, err := r.Cookie("___uuid")
	vEmail := v.Value
	vNum, _ := libs.Uuid.Search(uuid, vEmail)
	fmt.Println(vNum)
	if err != nil {
		fmt.Fprintln(w, "Something Error .....")
		return
	}
	cNum := r.FormValue("ver")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if vNum == cNum {
		fmt.Fprintln(w, `Register Success!
		<a href="/">Please Back to main page.</a>`)

	} else {
		fmt.Fprintln(w, `Register Failed!
		<a href="/">Please Back to main page.</a>`)
	}
}

func timeout(w http.ResponseWriter, r *http.Request) {
	v, err1 := r.Cookie("___uuid")
	cs := r.Cookies()

	if err1 != nil {
		for _, c1 := range cs {
			deleteCookie(w, c1.Name)
		}
		fmt.Fprintln(w, "DO NOT EDIT COOKIE'S NAME, FUCK YOU")
		return
	}
	vv := v.Value // mdStr
	err := libs.Uuid.Delete(uuid, vv)
	if err != nil {
		for _, c1 := range cs {
			deleteCookie(w, c1.Name)
		}
		fmt.Fprintln(w, "What ?! Are you HACKER ?!?!")
		return
	}

	// delete cookie
	c := http.Cookie{
		Name:     "___uuid",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, &c)

	fmt.Fprintln(w, "Timed out !")
}

func main() {
	server := http.Server{
		Addr: "0.0.0.0:1234",
	}
	http.HandleFunc("/", index)
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/process", process)
	http.HandleFunc("/process_verify", process_verify)
	http.HandleFunc("/timeout", timeout)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	server.ListenAndServe()
}
