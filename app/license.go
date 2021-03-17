package app

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type response struct {
	Success bool   `json:"success"`
	Reason  string `json:"reason"`
}

var Valid = true

// LicenseCheck Check the license
func (app *Application) LicenseCheck() {

	payload := map[string]string{
		"email":   "jcoenep@gmail.com",
		"license": "f36f66f454e44dd02a1f40a9d65ef2f649db21a1",
	}

	payloadJSON, err := json.Marshal(payload)

	if err != nil {
		Valid = false
		app.ErrorLog.Panic("Internal: Could not check license")
	}

	req, err := http.NewRequest(
		"POST", "http://localhost:3000/check", bytes.NewBuffer(payloadJSON))

	if err != nil {
		Valid = false
		app.ErrorLog.Fatal("Could not ping license endpoint")
	}

	req.Header.Set("skeef-token", "9a525de4769cca6398d0019b909bba84e26a2b80")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		Valid = false
		app.ErrorLog.Fatal("Could not ping license endpoint")
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		Valid = true
		return
	}

	var result response

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		Valid = true
		app.ErrorLog.Fatal("Could not check license")
	}

	if result.Success {
		Valid = true
		return
	}

	Valid = false

	app.ErrorLog.Fatal(result.Reason)

}
