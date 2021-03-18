package app

func (app *Application) TablesExists() bool {
	userq, err := app.Database.Con.Query("SELECT COUNT(1) FROM users;")

	if err != nil {
		return false
	}

	defer userq.Close()

	nusers := 0
	users := false
	for userq.Next() {
		userq.Scan(&nusers)
	}

	if nusers > 0 {
		users = true
	}

	licenseq, err := app.Database.Con.Query("SELECT COUNT(1) FROM license;")

	if err != nil {
		return false
	}

	defer licenseq.Close()

	nlicenses := 0
	license := false
	for licenseq.Next() {
		licenseq.Scan(&nlicenses)
	}

	if nlicenses > 0 {
		license = true
	}

	if license && users {
		return true
	}

	return false
}

func (app *Application) LicenseExists() bool {
	rows, err := app.Database.Con.Query("SELECT COUNT(1) FROM license;")

	if err != nil {
		return false
	}

	result := false
	for rows.Next() {
		result = true
	}

	return result
}

func (app *Application) AdminExists() bool {
	rows, err := app.Database.Con.Query("SELECT COUNT(1) FROM users;")

	if err != nil {
		return false
	}

	result := false
	for rows.Next() {
		result = true
	}

	return result
}
