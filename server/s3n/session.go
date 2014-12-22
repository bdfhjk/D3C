package s3n

import (
	"net/http"

	"github.com/Lealen/mysqlstore"
	"github.com/gorilla/sessions"
)

type Session struct {
	S                 *sessions.Session
	needToSaveChanges bool
}

var SessionStore *mysqlstore.MySQLStore

func initsessions(keyPairs ...[]byte) {
	var err error
	SessionStore, err = mysqlstore.NewMySQLStoreFromConnection(Db(), "sessions", GC().Cookies.Path, GC().Cookies.Maxage, keyPairs...)
	if err != nil {
		panic(err)
	}
}

func GetSession(r *http.Request) *Session {
	sess, _ := SessionStore.New(r, "s3n-session")
	return &Session{sess, false}
}

func (s *Session) Save(w http.ResponseWriter, r *http.Request) {
	if s.needToSaveChanges {
		s.needToSaveChanges = false
		SessionStore.Save(r, w, s.S)
	}
}

func (s *Session) Get(key interface{}) interface{} {
	return s.S.Values[key]
}

func (s *Session) Set(key interface{}, value interface{}) {
	s.needToSaveChanges = true
	s.S.Values[key] = value
}

func (s *Session) Clear() {
	for k := range s.S.Values {
		delete(s.S.Values, k)
	}
}
