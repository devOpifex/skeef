package db

// CreateUserTable Create user table
func (DB *Database) CreateTableUser() error {

	_, err := DB.Con.Exec("CREATE TABLE users (email VARCHAR(50) NOT NULL PRIMARY KEY, hashed_password CHAR(60) NOT NULL, admin INTEGER);")

	if err != nil {
		return err
	}

	return nil
}

func (DB *Database) CreateTableLicense() error {

	_, err := DB.Con.Exec("CREATE TABLE license (email VARCHAR(50) NOT NULL PRIMARY KEY, license VARCHAR(255) NOT NULL);")

	if err != nil {
		return err
	}

	return nil
}

func (DB *Database) CreateTableTwitterApp() error {

	_, err := DB.Con.Exec("CREATE TABLE twitter_app (api_key VARCHAR(255) NOT NULL, api_secret VARCHAR(255) NOT NULL, access_token VARCHAR(255) NOT NULL, access_secret VARCHAR(255) NOT NULL, id INTEGER);")

	if err != nil {
		return err
	}

	return nil
}

func (DB *Database) CreateTableStreams() error {

	_, err := DB.Con.Exec("CREATE TABLE streams (name VARCHAR(50) NOT NULL PRIMARY KEY, follow VARCHAR(400), track VARCHAR(400), locations VARCHAR(400), active INTEGER, max_edges INTEGER, exclude VARCHAR(254));")

	if err != nil {
		return err
	}

	return nil
}
