package s3nCusers

import (
	"net/http"

	"../s3n"
	"github.com/gorilla/schema"
)

const minloginchars int8 = 4

type LoginForm struct {
	Login    string
	Password string
}

func LoginFormHandlerWithoutCaptcha(w http.ResponseWriter, r *http.Request, b *s3n.ResponseBuffer, s *s3n.Session) {
	if err := r.ParseForm(); err != nil {
		s3n.Error(w, r, 400, err)
		return
	}

	loginForm := new(LoginForm)
	decoder := schema.NewDecoder()

	if err := decoder.Decode(loginForm, r.PostForm); err != nil {
		s3n.Error(w, r, 400, err)
		return
	}

	id, err := Login(loginForm.Login, loginForm.Password, false)
	if err != nil {
		if err == ErrNeedToSolveCaptcha {
			s3n.Redirect(w, r, s, "/captcha?back=login", 302)
			return
		}
		s3n.Redirect(w, r, s, "/login?err", 302)
		return
	}

	s.Set("id", id)
	s3n.Redirect(w, r, s, "/", 302)
}

func LoginFormHandler(w http.ResponseWriter, r *http.Request, b *s3n.ResponseBuffer, s *s3n.Session) {
	if err := r.ParseForm(); err != nil {
		s3n.Error(w, r, 400, err)
		return
	}

	loginForm := new(LoginForm)
	decoder := schema.NewDecoder()

	if err := decoder.Decode(loginForm, r.PostForm); err != nil {
		s3n.Error(w, r, 400, err)
		return
	}

	id, err := Login(loginForm.Login, loginForm.Password, true)
	if err != nil {
		if err == ErrNeedToSolveCaptcha {
			s3n.Redirect(w, r, s, "/captcha?back=login", 302)
			return
		}
		s3n.Redirect(w, r, s, "/login?err", 302)
		return
	}

	s.Set("id", id)
	s3n.Redirect(w, r, s, "/", 302)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request, b *s3n.ResponseBuffer, s *s3n.Session) {
	s.Clear()
	s3n.Redirect(w, r, s, "/", 302)
}


type RegisterForm struct {
	Login       string
	Password    string
	Email       string
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

	//data := make(map[string]interface{})
	//data["something"] = registerForm.Something
	err := AddUser(registerForm.Login, registerForm.Password, registerForm.Email, nil /*data*/)
	if err != nil {
		s3n.Log.Print(err)
	}
	s3n.Redirect(w, r, s, "/", 302)
}


/*
func TestHandler(w http.ResponseWriter, r *http.Request, b *s3n.ResponseBuffer, s *s3n.Session) {
	s3n.Log.Printf("\n%#v\n", s.S)

	if id := s.Get("id"); id != nil {
		b.Printf("id-%d", id)
		data, err := GetUserData(id.(int))
		if err != nil {
			for index, value := range data {
				b.Print(index)
				b.Print("-")
				b.Println(value.(string))
			}
		}
	}
}

func TestSetHandler(w http.ResponseWriter, r *http.Request, b *s3n.ResponseBuffer, s *s3n.Session) {
	s.Set("id", 0)
}*/
