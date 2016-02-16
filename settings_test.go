package cityhall

import (
	"testing"
	"os"
)

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
		t.Errorf("Got an error back creating settings")
	}
	if s.passhash != hash(password) {
		t.Errorf("Expected the password to be hashed, if non-empty")
	}
}

func TestLoginRetrievesDefaultEnvironment(t *testing.T) {
	cityhall := (&mockServer{}).createMockWithLogin(t)
	defer cityhall.Close()

	s, _ := NewSettings(CityHallInfo{Url: cityhall.URL})
	err := s.login()
	if err != nil {
		t.Errorf("Login expected to POST to /auth/ with hostname and empty password: %s", err)
	}
}

func TestGetDefaultEnvironment(t *testing.T) {
	cityhall := (&mockServer{}).createMockWithLogin(t)
	defer cityhall.Close()

	s, _ := NewSettings(CityHallInfo{Url: cityhall.URL})
	err := s.getDefaultEnvironment()
	if err != nil {
		t.Errorf("getDefaultEnvironment should've GET at /auth/user/{hostname}/default/: %s", err)
	}
	if s.Environments.Default() != "dev" {
		t.Errorf("Expected default Environment to be 'dev'")
	}
}
