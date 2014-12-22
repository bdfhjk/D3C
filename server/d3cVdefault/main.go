package d3cVdefault

import (
	"net/http"

	"../s3n"
	"../s3nCusers"
)

func DefaultHandler(w http.ResponseWriter, r *http.Request, b *s3n.ResponseBuffer, s *s3n.Session) {
	data := make(map[string]interface{})
	var err error

	data["loggedin"]=false

	if id := s.Get("id"); id != nil {
		userdata, err := s3nCusers.GetUserData(id.(int))
		if err == nil {
			data["LoggedIn"]=true
			data["FirstName"]=userdata["FirstName"]
			data["SecondName"]=userdata["SecondName"]
		}
	}/* else {
		s3n.Redirect(w, r, s, "/login", 302)
		return
	}*/

	err = b.ExecuteTemplate("default", data)
	if err != nil {
		s3n.Log.Print(err)
	}
}

func GameHandler(w http.ResponseWriter, r *http.Request, b *s3n.ResponseBuffer, s *s3n.Session) {
	data := make(map[string]interface{})
	var err error

	data["loggedin"]=false

	if id := s.Get("id"); id != nil {
		userdata, err := s3nCusers.GetUserData(id.(int))
		if err == nil {
			data["LoggedIn"]=true
			data["FirstName"]=userdata["FirstName"]
			data["SecondName"]=userdata["SecondName"]
		} else {
			s3n.Redirect(w, r, s, "/", 302)
			return
		}
	} else {
		s3n.Redirect(w, r, s, "/", 302)
		return
	}

	err = b.ExecuteTemplate("game", data)
	if err != nil {
		s3n.Log.Print(err)
	}
}
