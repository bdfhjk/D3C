package s3nCusers

import (
	"errors"
	"strings"

	"../s3n"
	"code.google.com/p/go.crypto/bcrypt"
)

var (
	ErrNeedToSolveCaptcha = errors.New("s3nCusers: user need to solve captcha")
	ErrUserNotLogIn       = errors.New("User not login")
)

//return user id and error
func Login(loginoremail, password string, counterrors bool) (int, error) {
	var (
		id,
		loginattempt int
	)
	var (
		login,
		passwordhash,
		email string
	)
	row := s3n.Db().QueryRow("SELECT id, login, passwordhash, loginattempt, email FROM "+Databasename+" WHERE LOWER(login) = ? OR LOWER(email) = ?", strings.ToLower(loginoremail), strings.ToLower(loginoremail))
	err := row.Scan(&id, &login, &passwordhash, &loginattempt, &email)
	if err != nil {
		return 0, err
	}

	if loginattempt > 5 && counterrors {
		return 0, ErrNeedToSolveCaptcha
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordhash), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword && counterrors {
			s3n.Db().Exec(`UPDATE `+Databasename+` SET
				loginattempt=loginattempt+1
				WHERE id=?
				`, id)
		}
		return 0, err
	}

	if loginattempt > 0 {
		s3n.Db().Exec(`UPDATE `+Databasename+` SET
			loginattempt=0
			WHERE id=?
			`, id)
	}

	s3n.Log.Printf("%s logged in", login)
	return id, nil
}

func AddUser(login, password, email string, decodeddata map[string]interface{}) error {
	if decodeddata == nil {
		decodeddata = make(map[string]interface{})
	}
	data, err := s3n.G().EncodeMultiNoRestrictions("user", decodeddata)
	if err != nil {
		return err
	}
	passwordhash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost+2)
	if err != nil {
		return err
	}
	if _, err = s3n.Db().Exec(`INSERT INTO `+Databasename+`
		(id, login, passwordhash, loginattempt, email, data) VALUES
		(NULL,?,?,0,?,?);`,
		login, string(passwordhash), email, data); err != nil {
		return err
	}
	return nil
}

func SetUserData(decodeddata map[string]interface{}) error {
	id := decodeddata["id"].(int)
	login := decodeddata["login"].(string)
	email := decodeddata["email"].(string)
	delete(decodeddata, "id")
	delete(decodeddata, "login")
	delete(decodeddata, "email")
	data, err := s3n.G().EncodeMultiNoRestrictions("user", decodeddata)
	if err != nil {
		return err
	}

	if _, err = s3n.Db().Exec(`UPDATE `+Databasename+` SET
	login=?,
	email=?,
	data=?
	WHERE id=?
	`, login, email, data, id); err != nil {
		return err
	}
	return nil
}

func GetUserData(id int) (map[string]interface{}, error) {
	var (
		login,
		email,
		data string
	)
	row := s3n.Db().QueryRow("SELECT login, email, data FROM "+Databasename+" WHERE id = ?", id)
	if err := row.Scan(&login, &email, &data); err != nil {
		return nil, err
	}
	decodeddata := make(map[string]interface{})
	if err := s3n.G().DecodeMultiNoRestrictions("user", data, &decodeddata); err != nil {
		return nil, err
	}
	decodeddata["id"] = id
	decodeddata["login"] = login
	decodeddata["email"] = email
	return decodeddata, nil
}

type AllUsersList struct {
	Id    int
	Login string
}

func GetAllUsers() ([]AllUsersList, error) {
	rows, err := s3n.Db().Query(`SELECT id,login
	FROM ` + Databasename + ` ORDER BY login ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []AllUsersList

	for rows.Next() {
		var (
			id    int
			login string
		)
		if err := rows.Scan(&id, &login); err != nil {
			return nil, err
		}

		result = append(result, AllUsersList{id, login})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
