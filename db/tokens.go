package db

type Tokens struct {
	ApiKey       string
	ApiSecret    string
	AccessToken  string
	AccessSecret string
}

func (DB *Database) InsertTokens(apiKey, apiSecret, accessToken, accessSecret string) error {

	stmt, err := DB.Con.Prepare("INSERT INTO twitter_app (api_key, api_secret, access_token, access_secret, id) VALUES (?,?,?,?,?);")

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(apiKey, apiSecret, accessToken, accessSecret, 1)

	if err != nil {
		return err
	}

	return nil
}

func (DB *Database) UpdateTokens(apiKey, apiSecret, accessToken, accessSecret string) error {

	stmt, err := DB.Con.Prepare("UPDATE twitter_app SET api_key = ?, api_secret = ?, access_token = ?, access_secret = ? WHERE id = 1")

	if err != nil {
		return err
	}

	_, err = stmt.Exec(apiKey, apiSecret, accessToken, accessSecret)

	if err != nil {
		return err
	}

	return nil
}

func (DB *Database) GetTokens() (Tokens, error) {
	res := DB.Con.QueryRow("SELECT api_key, api_secret, access_token, access_secret FROM twitter_app;")

	tk := Tokens{}
	res.Scan(&tk.ApiKey, &tk.ApiSecret, &tk.AccessToken, &tk.AccessSecret)

	return tk, nil
}
