package common

import (
	// "fmt"
	"log"
	"net/http"
	"time"

	"github.com/catinello/base62"
	// "strconv"
)

// MakeCSRF save current unix time and encrypt
func MakeCSRF(w http.ResponseWriter, r *http.Request) string {
	u62 := base62.Encode(int(time.Now().Unix()))
	page := Encrypt(CsrfKey, u62)
	cookie := &http.Cookie{
		Name:     "xr",
		Value:    u62,
		Path:     "/",
		MaxAge:   0,
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	return page
}

// CheckCSRF check value exists, can be decrypt, not expired
func CheckCSRF(r *http.Request, page string, w http.ResponseWriter) bool {
	cookie, err := r.Cookie("xr")
	u62 := ""
	if err == nil {
		u62 = cookie.Value
	}
	if Decrypt(CsrfKey, page) != u62 || u62 == "" {
		log.Print("CSRF page cookie is not match")
		return false
	}
	n, err := base62.Decode(u62)
	if err != nil {
		log.Print(err)
	}
	tm := time.Unix(int64(n), 0)
	tAdd := tm.Add(2 * time.Hour)
	if time.Now().After(tAdd) {
		log.Print("CSRF is expired")
		return false
	}
	// var w http.ResponseWriter
	cookie = &http.Cookie{
		Name:     "xr",
		Value:    "",
		Path:     "/",
		MaxAge:   0,
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	return true
}
