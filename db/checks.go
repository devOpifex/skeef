package db

func (db Database) TablesExists() bool {
	userq, err := db.Con.Query("SELECT COUNT(1) FROM users;")

	if err != nil {
		return false
	}

	defer userq.Close()

	users := false
	for userq.Next() {
		users = true
	}

	licenseq, err := db.Con.Query("SELECT COUNT(1) FROM license;")

	if err != nil {
		return false
	}

	defer licenseq.Close()

	license := false
	for licenseq.Next() {
		license = true
	}

	if license && users {
		return true
	}

	return false
}

func (db Database) LicenseExists() bool {
	rows, err := db.Con.Query("SELECT COUNT(1) FROM license;")

	if err != nil {
		return false
	}

	result := false
	for rows.Next() {
		result = true
	}

	return result
}

func (db Database) AdminExists() bool {
	rows, err := db.Con.Query("SELECT COUNT(1) FROM users;")

	if err != nil {
		return false
	}

	result := false
	for rows.Next() {
		result = true
	}

	return result
}
