package s3n

import (
	"code.google.com/p/gcfg"
)

const _debug = true

type Config struct {
	S3n struct {
		DisableLogs bool
		Debug       bool
	}
	Database struct {
		Usemysql bool
		Autoconnect bool
		Host     string
		User     string
		Password string
		Dbname   string
	}
	Cache struct {
		On            bool
		Folder        string
		Maxsize       int64
		Maxsizeforall int64
	}
	Cookies struct {
		Path	   string
		Maxage   int
	}
}

var cfg Config

func loadConfig(str string) error {
	err := gcfg.ReadFileInto(&cfg, str)

	if err != nil {
		Log.Printf("CG Can't read config file! Reason: %s", err)
	}

	cfg.S3n.Debug = _debug

	return err
}

func GC() Config {
	return GetConfig()
}

func GetConfig() Config {
	return cfg
}
