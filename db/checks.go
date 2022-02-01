package db

func (db Database) AdminExists() bool {
	rows := db.Con.QueryRow("SELECT COUNT(1) FROM users;")

	var count int
	rows.Scan(&count)

	return count > 0
}

func (db Database) TokensExist() bool {
	rows := db.Con.QueryRow("SELECT COUNT(1) FROM twitter_app;")

	var count int
	rows.Scan(&count)

	return count > 0
}
