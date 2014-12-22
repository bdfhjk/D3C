package s3n

import (
	"html"
	"net/http"
	//	"database/sql"
	//	_ "github.com/go-sql-driver/mysql"
)

type typeuser struct {
	id    uint32
	login string
	name  string
	email string
	guest bool
	admin bool
	moder bool
	ip    string
}

func inituser(w http.ResponseWriter, r *http.Request) *typeuser {
	var u *typeuser
	u = new(typeuser)
	u.id = 0
	u.login = ""
	u.name = ""
	u.email = ""
	u.guest = true
	u.admin = false
	u.moder = false
	u.ip = r.RemoteAddr

	sessid, err := r.Cookie("s3n-session")
	if err == nil {
		u.loaduser(w, r, sessid.Value)
	}
	return u
}

func (u *typeuser) loaduser(w http.ResponseWriter, r *http.Request, sessid string) {
	dbp, err := Db().Prepare("SELECT * FROM s3n_sessions WHERE sessid = ?")
	if err != nil {
		Error(w, r, 503, err)
	}
	defer dbp.Close()

	var und int
	err = dbp.QueryRow(html.EscapeString(sessid)).Scan(&und)
	if err != nil {
		return
	}

}
