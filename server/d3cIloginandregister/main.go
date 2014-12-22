package d3cIloginandregister

import (
	"net/http"
	"github.com/gorilla/schema"

	"../s3n"
	"../s3nCusers"
)

/*
func LoginHandler(w http.ResponseWriter, r *http.Request, b *s3n.ResponseBuffer, s *s3n.Session) {
	data := make(map[string]interface{})
	//data["something"] = something
	err := b.ExecuteTemplate("login", data)
	if err != nil {
		s3n.Log.Print(err)
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request, b *s3n.ResponseBuffer, s *s3n.Session) {
	data := make(map[string]interface{})
	//data["something"] = something
	err := b.ExecuteTemplate("register", data)
	if err != nil {
		s3n.Log.Print(err)
	}
}
*/

type RegisterForm struct {
	FirstName,
	SecondName,
	Email,
	Password string
}

func RegisterFormHandler(w http.ResponseWriter, r *http.Request, b *s3n.ResponseBuffer, s *s3n.Session) {
	if err := r.ParseForm(); err != nil {
		s3n.Error(w, r, 400, err)
		return
	}

	registerForm := new(RegisterForm)
	decoder := schema.NewDecoder()

	if err := decoder.Decode(registerForm, r.PostForm); err != nil {
		s3n.Error(w, r, 400, err)
		return
	}

	data := make(map[string]interface{})
	data["FirstName"] = registerForm.FirstName
	data["SecondName"] = registerForm.SecondName
	err := s3nCusers.AddUser(registerForm.Email, registerForm.Password, registerForm.Email, data)
	if err != nil {
		s3n.Log.Print(err)
	}
	s3n.Redirect(w, r, s, "/", 302)
}
