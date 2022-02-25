package common

import (
	// "fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

// SetUser save id and current time
func SetUser(w http.ResponseWriter, r *http.Request, id int) {
	txt := strconv.Itoa(id)
	txt = Encrypt(SsKey, txt)
	cookie := &http.Cookie{
		Name:     "ss",
		Value:    txt,
		MaxAge:   101556952,
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
	txt = time.Now().Format("2006-01-02 15:04:05")
	txt = Encrypt(T1Key, txt)
	cookie1 := &http.Cookie{
		Name:     "ti",
		Value:    txt,
		MaxAge:   101556952,
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, cookie1)
}

// GetUser check value exsist decrytable not expired
func GetUser(w http.ResponseWriter, r *http.Request) int {
	usrID := 0
	cookie, err := r.Cookie("ss")
	if err != nil {
		return 0
	}
	ss := Decrypt(SsKey, cookie.Value)

	delete := false
	if ss == "" {
		log.Print("ss is wrong: ", err)
		delete = true
	}
	usrID, err = strconv.Atoi(ss)
	if err != nil {
		log.Print("ss error: ", err)
	}
	cookie, err = r.Cookie("ti")
	if err != nil {
		log.Print("No ti Cookie: ", err)
	}
	t1 := Decrypt(T1Key, cookie.Value)

	if t1 == "" {
		log.Print("ti is wrong: ", err)
		delete = true
	}
	stampTime, err := time.Parse("2006-01-02 15:04:05", t1)

	if err != nil {
		log.Print("time.Parse error: ", err)
		delete = true
	}
	t1Add := stampTime.AddDate(1, 0, 0)
	if time.Now().After(t1Add) {
		log.Print("session expired")
		delete = true
	}
	if delete {
		cookie := &http.Cookie{
			Name:     "ss",
			Value:    "",
			MaxAge:   0,
			Secure:   true,
			HttpOnly: true,
			Path:     "/",
		}
		http.SetCookie(w, cookie)
		cookie1 := &http.Cookie{
			Name:     "ti",
			Value:    "",
			MaxAge:   0,
			Secure:   true,
			HttpOnly: true,
			Path:     "/",
		}
		http.SetCookie(w, cookie1)
		return 0
	}
	return usrID
}
