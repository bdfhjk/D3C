package s3n

import (
	"compress/gzip"
	"net/http"
	"runtime/debug"
	//"io/ioutil"
	"fmt"
	"log"
	"strconv"
	"strings"
	//"strings"
	"os"
	//"bytes"
	"io"
	"io/ioutil"
)

var logfile *os.File
var Log *log.Logger

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	if "" == w.Header().Get("Content-Type") {
		w.Header().Set("Content-Type", http.DetectContentType(b))
	}
	return w.Writer.Write(b)
}

func Handler(handle func(w http.ResponseWriter, r *http.Request, b *ResponseBuffer, s *Session)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		b, s := BeginLoadPage(w, r)

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			defer func() {
				if e := recover(); e != nil {
					Log.Printf("PANIC RECOVERED:%s\n%s\n", e, debug.Stack())
					Error(w, r, http.StatusInternalServerError, fmt.Errorf("%v", e))
				}
				EndLoadPage(w, r, b, s)
			}()

			handle(w, r, b, s)
		} else {
			w.Header().Set("Content-Encoding", "gzip")
			gz := gzip.NewWriter(w)
			defer gz.Close()
			gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}

			defer func() {
				if e := recover(); e != nil {
					Log.Printf("PANIC RECOVERED:%s\n%s\n", e, debug.Stack())
					Error(w, r, http.StatusInternalServerError, fmt.Errorf("%v", e))
				}
				EndLoadPage(gzr, r, b, s)
			}()

			handle(gzr, r, b, s)
		}
	}
}

func Error(w http.ResponseWriter, r *http.Request, number int, err error) {
	if GC().S3n.Debug && err != nil {
		Log.Printf("<- %s <- %s <- %s[%s] <- Error %d <- %s", r.RemoteAddr, r.Host+r.RequestURI, r.Method, r.Form, number, err.Error())
	} else {
		Log.Printf("<- %s <- %s <- %s[%s] <- Error %d", r.RemoteAddr, r.Host+r.RequestURI, r.Method, r.Form, number)
	}

	switch number {
	case http.StatusForbidden:
		http.Error(w, strconv.Itoa(number)+" forbidden or no permission to access", number)
	case 400:
		http.Error(w, strconv.Itoa(number)+" invalid request form data", number)
	case http.StatusNotFound:
		http.Error(w, strconv.Itoa(number)+" page not found", number)
	case http.StatusInternalServerError:
		http.Error(w, strconv.Itoa(number)+" internal server error", number)
	case http.StatusServiceUnavailable:
		http.Error(w, strconv.Itoa(number)+" service unavailable", number)
	default:
		http.Error(w, strconv.Itoa(number), number)
	}
}

func LogUrl(w http.ResponseWriter, r *http.Request) {
	if GC().S3n.DisableLogs {
		return
	}
	logtoprint := fmt.Sprintf("-> %s -> %s -> %s[%s]", r.RemoteAddr, r.Host+r.RequestURI, r.Method, r.Form)
	for index, element := range r.Header {
		logtoprint = logtoprint + fmt.Sprintf(" -> %s: %s", index, element)
	}
	Log.Printf(logtoprint)
}

func BeginLoadPage(w http.ResponseWriter, r *http.Request) (*ResponseBuffer, *Session) {
	LogUrl(w, r)
	Setheader(w)

	var b *ResponseBuffer
	b = new(ResponseBuffer)
	return b, GetSession(r)
}

func EndLoadPage(w http.ResponseWriter, r *http.Request, b *ResponseBuffer, s *Session) {
	s.Save(w, r)
	//TODO: http://golang.org/pkg/strings/#Replacer
	for {
		newb := strings.Replace(b.b, "  ", " ", -1)
		if newb == b.b {
			break
		}
		b.b = newb
	}
	for {
		newb := strings.Replace(b.b, "\t\t", "\t", -1)
		if newb == b.b {
			break
		}
		b.b = newb
	}
	for {
		newb := strings.Replace(b.b, "\r", "\n", -1)
		if newb == b.b {
			break
		}
		b.b = newb
	}
	for {
		newb := strings.Replace(b.b, "\n ", "\n", -1)
		if newb == b.b {
			break
		}
		b.b = newb
	}
	for {
		newb := strings.Replace(b.b, "\n\t", "\n", -1)
		if newb == b.b {
			break
		}
		b.b = newb
	}
	for {
		if b.b[:1] == "\n" {
			b.b = b.b[1:]
		} else {
			break
		}
	}
	for {
		if len(b.b) <= 2 {
			break
		} else if b.b[:1] == "\n" {
			b.b = b.b[1:]
		} else {
			break
		}
	}
	if len(b.b) > 2 && (b.b[:1] == "<" || b.b[:2] == "\n<" || b.b[:2] == " <" || b.b[:2] == "\t<") {
		fmt.Fprint(w, strings.Join(strings.SplitAfterN(b.b, "\n", 2), `
<!--
==============================================================================
By Mateusz Doroszko aka Lealen bez
==============================================================================
Cześć!
Jaką książkę chcesz dziś przeczytać?
A.. Tak. Już wybrałeś.
A więc zacznijmy.
Spoglądasz teraz na dziwny komentarz, czyż nie? Mam nadzieję, że nie będzie to dla Ciebie utrapieniem.
Pamiętaj, że to co tu widzisz zostało poświęcone przynajmniej kilkoma tygodniami pracy,
	ale były to zabawne i ciekawe tygodnie, ponieważ mogłem napisać to wszystko w Go (golang.org)
	poczynając od frameworka, przez cały serwer, itd.
Najzabawniejsze w tym wszystkim jest szukanie błędów, ale to historia na kolejną opowieść.
Pamiętaj, że w tej historii występują smoki.
Ale koniec czytania na dzisiaj.
Chyba powinieneś Ty również poznawać świat.
Nie powinieneś wyruszyć w drogę?
Na przykład do krainy programowania, do krainy pięknych kodów źródłowych,
	do krainy komentarzy, do krainy wyciekami pamięci płynącymi,
	do krainy czyściciela pamięci, do krainy sprawdzania błędów,
	do krainy struktur.
Tylko nie zapomnij wpaść do mnie czasem, jak za starych dobrych czasów.
Jeśli się zgubisz zapytaj podróżnika o mój adres domowy: http://lealen.pl
Jeśli okaże się, że to za daleko możesz jeszcze przeszukać: https://github.com/Lealen
Bądź wysłać list z doręczeniem automagicznym na adres: by>@<lealen.pl
A teraz ruszaj w świat tego kodu.
Powodzenia w twojej podróży, mam nadzieję, że znajdziesz to czego szukasz.

Lealen

PS Ten komentarz został dodany automatycznie przez serwer, ponieważ wszystko pozostałe może być edytowane. L.
PPS Chodzi mi o to, że style nie muszą być pisane przezemnie, ponieważ framework, który napisałem jest uniwersalny. L.
PPPS Miłego dnia. L.
PPPPS /l33t. L.
==============================================================================
========== |  _  _ | _  _   _ |    ()     /*_ _ | _  _  _   _  _ _  ==========
========== |_(/_(_||(/_| |.|_)|    (_X    \_/(_)|(_|| |(_|.(_)| (_| ==========
===========================|============================_|=======_|===========
-->

`))
	} else {
		fmt.Fprint(w, b.b)
	}
	//b (ResponseBuffer) free memory here
}

func Redirect(w http.ResponseWriter, r *http.Request, s *Session, urlStr string, code int) {
	s.Save(w, r)
	http.Redirect(w, r, urlStr, code)
}

func setlogoutput() {
	logfile, err := os.OpenFile("logs", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(logfile)
	Log = log.New(io.MultiWriter(logfile, os.Stdout), "", log.Lshortfile)
}

const configfile string = "config.gcfg"

func init() {
	setlogoutput()

	if _, err := os.Stat(configfile); os.IsNotExist(err) {
		err = ioutil.WriteFile(configfile, []byte(`; s3n configuration file

[database]
usemysql = true
autoconnect = true
host = 127.0.0.1:3306
user = root
password =
dbname = s3ndb

[cookies]
path = /
maxage = 604800 # 7 days * 24 hours * 60 minutes * 60 seconds
`), 0644)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Plik " + configfile + " został utworzony, uzupełnij go i uruchom program ponownie.")
		os.Exit(1)
	}

	fmt.Println(`
==============================================================================
========== |  _  _ | _  _   _ |    ()     /*_ _ | _  _  _   _  _ _  ==========
========== |_(/_(_||(/_| |.|_)|    (_X    \_/(_)|(_|| |(_|.(_)| (_| ==========
===========================|============================_|=======_|===========
`)

	if loadConfig(configfile) != nil {
		return
	}
}

func Init(keyPairs ...[]byte) {

	InitGorilla(keyPairs...)
	initsessions(keyPairs...)
}

func Exit() {
	logfile.Close()
}

type HeaderFuncType func(w http.ResponseWriter)

var HeaderFunc HeaderFuncType = func(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "SAMEORIGIN")
	w.Header().Set("Connection", "keep-alive")
}

func Setheader(w http.ResponseWriter) {
	HeaderFunc(w)
}
