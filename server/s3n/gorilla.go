package s3n

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

type GorillaT struct {
	GorillaRouter        *mux.Router
	Codecs               []securecookie.Codec
	CodecsNoRestrictions []securecookie.Codec
}

var Gorilla GorillaT

func G() *GorillaT {
	return &Gorilla
}

func CodecsFromPairs(keyPairs ...[]byte) []securecookie.Codec {
	codecs := make([]securecookie.Codec, len(keyPairs)/2+len(keyPairs)%2)
	for i := 0; i < len(keyPairs); i += 2 {
		var blockKey []byte
		if i+1 < len(keyPairs) {
			blockKey = keyPairs[i+1]
		}
		codecs[i/2] = securecookie.New(keyPairs[i], blockKey)
		codecs[i/2].(*securecookie.SecureCookie).MaxAge(0)
	}
	return codecs
}

func InitGorilla(keyPairs ...[]byte) {
	if len(keyPairs) > 0 {
		Gorilla = GorillaT{mux.NewRouter(), securecookie.CodecsFromPairs(keyPairs...), CodecsFromPairs(keyPairs...)}
	} else {
		Gorilla = GorillaT{mux.NewRouter(), nil, nil}
	}
}

func (g GorillaT) PathPrefix(tpl string) *mux.Route {
	return g.GorillaRouter.PathPrefix(tpl)
}

func (g GorillaT) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route {
	return g.GorillaRouter.HandleFunc(path, f)
}

func (g GorillaT) Handle(path string) {
	http.Handle(path, g.GorillaRouter)
}

func (g GorillaT) Vars(r *http.Request) map[string]string {
	return mux.Vars(r)
}

func (g GorillaT) EncodeMulti(name string, value interface{}) (string, error) {
	return securecookie.EncodeMulti(name, value, g.Codecs...)
}

func (g GorillaT) DecodeMulti(name string, value string, dst interface{}) error {
	return securecookie.DecodeMulti(name, value, dst, g.Codecs...)
}

func (g GorillaT) EncodeMultiNoRestrictions(name string, value interface{}) (string, error) {
	return securecookie.EncodeMulti(name, value, g.CodecsNoRestrictions...)
}

func (g GorillaT) DecodeMultiNoRestrictions(name string, value string, dst interface{}) error {
	return securecookie.DecodeMulti(name, value, dst, g.CodecsNoRestrictions...)
}

func (g GorillaT) DefaultClearHandler() http.Handler {
	return context.ClearHandler(http.DefaultServeMux)
}
