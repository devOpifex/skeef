package db

type License struct {
	Email   string
	License string
}

// InsertLicense Insert the user license
func (DB *Database) InsertLicense(email, license string) error {

	stmt, err := DB.Con.Prepare("INSERT INTO license (email, license) VALUES (?,?);")

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(email, license)

	if err != nil {
		return err
	}

	return nil
}

func (DB *Database) GetLicense() (License, error) {
	res := DB.Con.QueryRow("SELECT license, email FROM license;")

	license := License{}
	res.Scan(&license.License, &license.Email)

	return license, nil
}

func (DB *Database) UpdateLicense(email, license string) error {

	stmt, err := DB.Con.Prepare("UPDATE license SET license = ? WHERE email = ?")

	if err != nil {
		return err
	}

	_, err = stmt.Exec(license, email)

	if err != nil {
		return err
	}

	return nil
}
