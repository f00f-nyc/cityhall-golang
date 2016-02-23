package cityhall

import (
	"testing"
	"os"
	"fmt"
)

const response_settings_json = `{
		"Response": "Ok",
		"value": "sample value",
		"protect": false
	}`

func TestHash(t *testing.T) {
	hashedEmpty := hash("")
	if hashedEmpty != "" {
		t.Errorf("City Hall convention: empty password should return an empty hash")
	}

	hashedPassword := hash("somepass")
	if hashedPassword == "" {
		t.Errorf("Passwords should be hashed")
	}
}

func TestCreateSettingsWithOnlyURL(t *testing.T) {
	s, err := NewSettings(CityHallInfo{Url: "http://not.a.real.url/api"})
	if err != nil {
		t.Errorf("Got an error back creating settings")
	}
	if s.passhash != "" {
		t.Errorf("Expected an empty password when not specified")
	}
	hostname, _ := os.Hostname()
	if s.username != hostname {
		t.Errorf("Expected the username to be set to the hostname when not specified")
	}
	if s.LoggedIn() {
		t.Errorf("Settings should not automatically log in")
	}
}

func TestCreateSettingsWithPassword(t *testing.T) {
	password := "password"
	s, err := NewSettings(CityHallInfo{Url: "http://not.a.real.url/api", Password: password})
	if err != nil {
		t.Fatal("Got an error back creating settings")
	}
	if s.passhash != hash(password) {
		t.Errorf("Expected the password to be hashed, if non-empty")
	}
}

//this functionality is really being tested in values_test.go
func TestGetValue(t *testing.T) {
	harness := &mockServer{
		Path: "/env/dev/app1/value1/",
		Method: "GET",
		Body: response_settings_json,
	}
	settings := harness.createHarness(t)

	val, err := settings.GetValue("/app1/value1")
	if err != nil {
		t.Fatal("GetValue value returned an error")
	}
	if val != "sample value" {
		t.Fatal("Got back an incorrect value")
	}
}

//this functionality is really being tested in values_test.go
func TestGetValueOverride(t *testing.T) {
	harness := &mockServer{
		Path: "/env/dev/app1/value1/?override=",
		Method: "GET",
		Body: response_settings_json,
	}
	settings := harness.createHarness(t)

	val, err := settings.GetValueOverride("/app1/value1", "")
	if err != nil {
		t.Fatal("GetValueOverride value returned an error")
	}
	if val != "sample value" {
		t.Fatal("Got back an incorrect value")
	}
}

//this functionality is really being tested in values_test.go
func TestGetValueEnvironment(t *testing.T) {
	harness := &mockServer{
		Path: "/env/qa/app1/value1/",
		Method: "GET",
		Body: response_settings_json,
	}
	settings := harness.createHarness(t)

	val, err := settings.GetValueEnvironment("app1/value1", "qa")
	if err != nil {
		t.Fatal("GetValueEnvironment value returned an error")
	}
	if val != "sample value" {
		t.Fatal("Got back an incorrect value")
	}
}

func TestUpdatePassword(t *testing.T) {
	hostname, _ := os.Hostname()
	harness := &mockServer{
		Path: fmt.Sprintf("/auth/user/%s/", hostname),
		Method: "PUT",
		RequestBody: fmt.Sprintf(`{"passhash": "%s"}`, hash("password")),
	}
	settings := harness.createHarness(t)

	err := settings.UpdatePassword("password")
	if err != nil {
		t.Fatal("UpdatePassword returned an error: ", err)
	}
}

func TestLogout(t *testing.T) {
	settings := (&mockServer{
			Path: "/auth/",
			Method: "DELETE",
	}).createHarness(t)

	err := settings.Logout()
	if err == nil {
		t.Error("Logging out when logged in should return an error")
	}

	settings.login()
	settings.Environments.getDefault()

	if !settings.LoggedIn() {
		t.Fatal("We should be already logged in")
	}

	settings.Logout()

	if settings.LoggedIn() {
		t.Fatal("We should be logged out")
	}
	if settings.loggedIn != loggedOut {
		t.Fatal("The state should be loggedOut")
	}
}

func TestSettingsFromUrl(t *testing.T) {
	settings, err := NewSettingsFromUrl("http://not.a.real.url/api/")
	if err != nil {
		t.Fatal("Got an error back from NewSettings")
	}
	if settings.Url != "http://not.a.real.url/api" {
		t.Fatal("Setting an Url should remove trailing slash")
	}
}
