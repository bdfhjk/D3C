package s3nCusers

import (
	"../s3n"
)

func Updatedb() error {
	prx, err := s3n.Db().Prepare("SELECT value FROM configint WHERE name = ?")
	if err!=nil {
		return err
	}
	defer prx.Close()

	var actualversion int
	err = prx.QueryRow("db_version_"+Databasename).Scan(&actualversion)
	if err!=nil {
		return err
	}

	if actualversion<Version {
		return upgradedb(actualversion, Version)
	} else if actualversion>Version {
		return downgradedb(actualversion, Version)
	}

	return nil
}

func upgradedb(from, to int) error {
	prx, err := s3n.Db().Prepare("UPDATE configint SET value=? WHERE name = ?")
	if err!=nil {
		return err
	}
	for i := from+1; i <= to; i++ {
		switch i {
		case 5:
			err = upgradedbto_5()
		}
		if err!=nil {
			return err
		} else {
			_, err = prx.Exec(i, "db_version_"+Databasename)
			if err!=nil {
				return err
			}
		}
	}

	return nil
}

func downgradedb(from, to int) error {
	prx, err := s3n.Db().Prepare("UPDATE configint SET value=? WHERE name = ?")
	if err!=nil {
		return err
	}
	for i := from-1; i >= to; i-- {
		switch i {
		case 5:
			err = downgradedbto_5()
		}
		if err!=nil {
			return err
		} else {
			_, err = prx.Exec(i, "db_version_"+Databasename)
			if err!=nil {
				return err
			}
		}
	}

	return nil
}
