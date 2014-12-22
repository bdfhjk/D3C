package main

import (
	"net/http"

	"../s3n"
	"../d3cVdefault"
	"../d3cIloginandregister"
	"../s3nCusers"
)

func init() {
	s3n.HeaderFunc = func(w http.ResponseWriter) {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		w.Header().Set("Connection", "keep-alive")
	}

	s3n.Init([]byte("12345678901234561234567890123456"), []byte("1234567890123456"))
}

func main() {
	defer s3n.Exit()

	s3n.Loadtemplates()

	s3n.G().HandleFunc("/", s3n.Handler(d3cVdefault.DefaultHandler))
	s3n.G().HandleFunc("/game", s3n.Handler(d3cVdefault.GameHandler))

	s3n.G().HandleFunc("/form/login", s3n.Handler(s3nCusers.LoginFormHandler))
	s3n.G().HandleFunc("/form/register", s3n.Handler(d3cIloginandregister.RegisterFormHandler))
	s3n.G().HandleFunc("/logout", s3n.Handler(s3nCusers.LogoutHandler))

	s3n.G().HandleFunc("/favicon.ico", s3n.FileHandler)
	s3n.G().PathPrefix("/assets/").HandlerFunc(s3n.FileHandler)
	s3n.G().PathPrefix("/uploads/").HandlerFunc(s3n.FileHandler)

	s3n.G().Handle("/")

	http.ListenAndServe(":18888", s3n.G().DefaultClearHandler())
}
