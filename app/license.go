package app

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// LicenseCheck Check the license
func (app *Application) LicenseCheck() {

	payload := map[string]string{
		"username": "angela",
		"license":  "e711ccb8d4f00ae6fd9b26dedf477f4238c28d7e",
	}

	payloadJSON, err := json.Marshal(payload)

	if err != nil {
		app.ErrorLog.Fatal("Could not check license")
	}

	req, err := http.NewRequest(
		"POST", "http://localhost:3000/check", bytes.NewBuffer(payloadJSON))

	if err != nil {
		app.ErrorLog.Fatal("Could not verify license")
	}

	req.Header.Set("skeef-token", "9a525de4769cca6398d0019b909bba84e26a2b80")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		app.ErrorLog.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return
	}

	app.ErrorLog.Fatal("Could not verify license")
}
