package s3n

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strings"
)

func Include(w http.ResponseWriter, r *http.Request, file string) string {
	str, err := ReadFile(file, true)

	if err != nil {
		return ""
	}

	return str
}

func Require(w http.ResponseWriter, r *http.Request, file string) string {
	str, err := ReadFile(file, true)

	if err != nil {
		Error(w, r, 503, err)
		return ""
	}

	return str
}

func IsDir(file string) (bool, error) {
	stat, err := os.Stat(file)
	if err != nil {
		return false, err
	} else {
		return stat.Mode().IsDir(), nil
	}
}

func FileHandler(w http.ResponseWriter, r *http.Request) { //TODO: FIX IT!
	r.ParseForm()
	LogUrl(w, r)
	//log.Printf("-> %s -> %s -> %s[%s]", r.RemoteAddr, r.RequestURI, r.Method, r.Form)
	file := r.URL.Path

	filetype := 0

	if file == "/favicon.ico" {
		file = file[1:]
	} else if r.URL.Path[1:len("assets/")+1] == "assets/" {
		file = r.URL.Path[len("assets/")+1:]
		filetype = 1
	} else if r.URL.Path[1:len("uploads/")+1] == "uploads/" {
		file = r.URL.Path[len("uploads/")+1:]
		filetype = 2
	} else {
		Log.Print("Valid names for s3n.FileHandler: favicon.ico, assets/*, uploads/*")
		Error(w, r, 501, nil)
		return
	}

	if filetype == 1 || filetype == 0 {
		isdir, err := IsDir("assets/" + file)
		if err != nil {
			Log.Print(err)
			Error(w, r, 404, nil)
			return
		} else if isdir {
			Log.Print("err file 1")
			Error(w, r, 404, nil)
			return
		}
	} else if filetype == 2 {
		isdir, err := IsDir("uploads/" + file)
		if err != nil {
			Log.Print(err)
			Error(w, r, 404, nil)
			return
		} else if isdir {
			Log.Print("err file 1")
			Error(w, r, 404, nil)
			return
		}
	}

	tmp := strings.Split(file, "/")
	if len(tmp) >= 1 {
		tmp = strings.Split(file, ".")
		//ext := tmp[len(tmp)-1]
		if filetype == 1 || filetype == 0 {
			ViewFile("assets/"+file, w, r, true)
		} else {
			ViewFile("uploads/"+file, w, r, true)
		}
	} else {
		Log.Print("err file 2")
		Error(w, r, 404, nil)
		return
	}
}

//View File, used by FileHandler
func ViewFile(file string, w http.ResponseWriter, r *http.Request, cacheit bool) {
	//	if cacheit {
	//		ServeCachedFile(w, r, file)
	//	} else {
	http.ServeFile(w, r, file)
	//	}
}

/**
 * Read Text File
 */
func ReadFile(file string, cacheit bool) (string, error) {
	if cacheit {
		CacheFile(file)
	}
	buf := bytes.NewBuffer(nil)
	f, err := OpenCachedFile(file)
	io.Copy(buf, f)
	f.Close()
	content := string(buf.Bytes())
	return strings.Trim(content, "\n"), err
}
