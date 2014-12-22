package s3n

import (
	"database/sql"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func Db() *sql.DB {
	var err error
	if db != nil {
		err = db.Ping()
	}
	if db == nil || err != nil {
		ConnectToDatabase()
	}
	return db
}

func ConnectToDatabase() {
	connectToDatabase(0)
}

func connectToDatabase(i int) {
	if GC().Database.Usemysql {
		var err error
		tmp := GC().Database.User + ":" + GC().Database.Password + "@tcp(" + GC().Database.Host + ")/" + GC().Database.Dbname
		if GC().S3n.Debug {
			Log.Printf("DB Attempt to connect to mysql: %s", tmp)
		} else {
			Log.Printf("DB Attempt to connect to mysql: %s", strings.NewReplacer(GC().Database.Password, "#######").Replace(tmp))
		}
		db, err = sql.Open("mysql", tmp)
		if err != nil {
			if GC().S3n.Debug {
				Log.Printf("DB Error when connecting to dabatase: %s", err)
			} else {
				Log.Print("DB Error when connecting to dabatase")
			}
			if i == 2 {
				//Error(w,r,503, err)
				panic(err.Error()) //tymczasowy panic
				//return
			} else {
				connectToDatabase(i + 1)
				return
			}
		}

		_, err = db.Query("SET NAMES utf8")
		if err != nil {
			if GC().S3n.Debug {
				Log.Printf("DB Error when connecting to dabatase: %s", err)
			} else {
				Log.Print("DB Error when connecting to dabatase")
			}
			if i == 2 {
				//Error(w,r,503, err)
				panic(err.Error()) //tymczasowy panic
				//return
			} else {
				connectToDatabase(i + 1)
				return
			}
		}

		Log.Printf("DB Connected to database")
	}

	checkandinstall()
}

func checkandinstall() {
	if _, err := Db().Exec("SELECT 1 FROM configint"); err != nil {
		Log.Print("Installing database")
		install()
	}
}

func install() {
	if _, err := Db().Exec(`CREATE TABLE configint
	(
		name varchar(64),
		value int(11)
	)
	DEFAULT CHARACTER SET utf8
	DEFAULT COLLATE utf8_general_ci;`); err != nil {
		panic(err)
	}
	if _, err := Db().Exec(`CREATE TABLE configstring
	(
		name varchar(64),
		value text
	)
	DEFAULT CHARACTER SET utf8
	DEFAULT COLLATE utf8_general_ci;`); err != nil {
		panic(err)
	}
}
