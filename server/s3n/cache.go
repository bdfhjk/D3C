package s3n

import (
	//	"github.com/golang/groupcache"
	"os"
	//	"os/exec"
	"errors"
	"net/http"
	"strings"
	//	"fmt"
	"io"
)

func cp(src, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	return d.Close()
}

func CacheFile(file string) error {
	if strings.Contains(file, "/../") {
		return errors.New("File adress cannot contain ../")
	} else if GC().Cache.On {
		if _, err := os.Stat(GC().Cache.Folder + file); err == nil {
			return nil
		} else if _, err := os.Stat(file); err != nil {
			return err
		} else {
			tmp := strings.Split(file, "/")
			err = os.MkdirAll(GC().Cache.Folder+strings.Join(tmp[:len(tmp)-1], "/"), os.FileMode(0700))
			cp(file, GC().Cache.Folder+file)
		}
	}

	return nil
}

func OpenCachedFile(file string) (*os.File, error) {
	if _, err := os.Stat(GC().Cache.Folder + file); err == nil {
		fileop, err := os.Open(GC().Cache.Folder + file)
		return fileop, err
	} else {
		fileop, err := os.Open(file)
		return fileop, err
	}
}

func ServeCachedFile(w http.ResponseWriter, r *http.Request, file string) {
	_, err := os.Stat(GC().Cache.Folder + file)
	stat2, err2 := os.Stat(file)

	if err != nil && err2 == nil && GC().Cache.On && !stat2.Mode().IsDir() {
		CacheFile(file)
		_, err = os.Stat(GC().Cache.Folder + file)
	}

	if err != nil || !GC().Cache.On || stat2.Mode().IsDir() {
		if err2 != nil || stat2.Mode().IsDir() {
			Error(w, r, 403, err2)
			return
		} else {
			http.ServeFile(w, r, file)
		}
	} else {
		http.ServeFile(w, r, GC().Cache.Folder+file)
	}
}
