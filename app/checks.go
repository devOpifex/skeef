package app

func (app *Application) TablesExists() bool {
	userq, err := app.Database.Con.Query("SELECT COUNT(1) FROM users;")

	if err != nil {
		return false
	}

	defer userq.Close()

	users := false
	for userq.Next() {
		users = true
	}

	licenseq, err := app.Database.Con.Query("SELECT COUNT(1) FROM licrenses;")

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
