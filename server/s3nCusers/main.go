package s3nCusers

import (
	"../s3n"
	//"code.google.com/p/go.crypto/bcrypt"
)

const Databasename string = "users"
const Version int = 4

func init() {
	if err := checkandinstall(); err != nil {
		panic(err)
	}
	if err := Updatedb(); err != nil {
		s3n.Log.Print(err)
	}
}

func checkandinstall() error {
	row := s3n.Db().QueryRow("SELECT count(*) FROM configint WHERE name = ?", "db_version_"+Databasename)
	var i int
	err := row.Scan(&i)
	if err != nil {
		return err
	}
	if i == 0 {
		return install()
	}
	return nil
}

func install() error {
	prx, err := s3n.Db().Prepare("INSERT INTO configint (name, value) VALUES (?,?);")
	if err != nil {
		return err
	}
	defer prx.Close()
	_, err = prx.Query("db_version_"+Databasename, Version)
	if err != nil {
		return err
	}
	if _, err := s3n.Db().Exec(`CREATE TABLE ` + Databasename + `
	(
		id int(11) NOT NULL AUTO_INCREMENT,
		login varchar(64),
		passwordhash varchar(256),
		loginattempt int(3),
		email varchar(128),
		data longblob,
		PRIMARY KEY (id)
	)
	DEFAULT CHARACTER SET utf8
	DEFAULT COLLATE utf8_general_ci
	AUTO_INCREMENT=1;`); err != nil {
		return err
	}
	return nil
}

func Uninstall() error {
	return nil
}
