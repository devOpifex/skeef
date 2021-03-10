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
		app.ErrorLog.Panic("Internal: Could not check license")
	}

	req, err := http.NewRequest(
		"POST", "http://localhost:3000/check", bytes.NewBuffer(payloadJSON))

	if err != nil {
		app.ErrorLog.Fatal("Could not ping license endpoint")
	}

	req.Header.Set("skeef-token", "9a525de4769cca6398d0019b909bba84e26a2b80")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		app.ErrorLog.Fatal("Could not ping license endpoint")
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return
	}

	app.ErrorLog.Fatal(`
		Invalid license: 
			1) It could simply be wrong; double check it is correct.
			2) It has expired, go to https://skeef.io to renew
			3) You have deployed the application on more than one machine (wait ~30 minutes to relaunch)
			If you are sure none of these apply contact: john@opifex.org 
	`)
}
