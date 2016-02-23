package cityhall

import (
	"testing"
)

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
	err := s.Environments.getDefault()
	if err != nil {
		t.Fatal("getDefaultEnvironment should've GET at /auth/user/{hostname}/default/: %s", err)
	}
	if s.Environments.Default() != "dev" {
		t.Errorf("Expected default Environment to be 'dev'")
	}
}

func TestEnsureLoggedIn(t *testing.T) {
	harness := &mockServer{}
	s := harness.createHarness(t)
	defer harness.CityHall.Close()

	err := s.ensureLoggedIn()
	if err != nil {
		t.Fatal("ensureLoggedIn returned an error")
	}

	if s.loggedIn != loggedIn {
		t.Fatal("The end result of ensureLoggedIn should be that it is logged in")
	}
	if s.Environments.Default() != "dev" {
		t.Fatal("The end result of ensureLoggedIn should be that it retrieves the default environment")
	}

	s.loggedIn = loggedOut
	err = s.ensureLoggedIn()
	if err == nil {
		t.Errorf("Expected an error when calling on a settings that's been logged out")
	}
}
