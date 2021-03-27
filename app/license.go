package app

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type LicenseResponse struct {
	Success bool   `json:"success"`
	Reason  string `json:"reason"`
}

type request struct {
	Email   string
	License string
	Ping    bool
}

// LicenseCheck Check the license
func (app *Application) LicenseCheck(ping bool) LicenseResponse {

	if app.License.License == "" {
		license, err := app.Database.GetLicense()

		if err != nil {
			return LicenseResponse{false, "Could not fetch license from database"}
		}

		app.License = license
	}

	payload := request{
		Email:   app.License.Email,
		License: app.License.License,
		Ping:    ping,
	}

	payloadJSON, err := json.Marshal(payload)

	if err != nil {
		app.ErrorLog.Println("Could not serialise payload")
		return LicenseResponse{false, "Could not serialise payload"}
	}

	req, err := http.NewRequest(
		"POST", "http://localhost:3000/check", bytes.NewBuffer(payloadJSON))

	if err != nil {
		app.ErrorLog.Println("Could not ping license endpoint")
		return LicenseResponse{false, "Could not ping license endpoint"}
	}

	req.Header.Set("skeef-token", "9a525de4769cca6398d0019b909bba84e26a2b80")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		app.ErrorLog.Println("Could not ping license endpoint")
		return LicenseResponse{false, "Could not ping license endpoint"}
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return LicenseResponse{false, "Could not ping license endpoint"}
	}

	var result LicenseResponse

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		app.ErrorLog.Println("Could not parse response from license endpoint")
		return LicenseResponse{false, "Could not parse response from license endpoint"}
	}

	app.ErrorLog.Println(result.Reason)

	return result
}

func (app *Application) LicenseValidity() {
	app.LicenseResponse = app.LicenseCheck(app.Streaming)

	if app.LicenseResponse.Success {
		return
	}

	app.Database.PauseAllStreams()
	app.Streaming = false
}
