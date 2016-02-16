package cityhall

import (
	"testing"
	"os"
	"fmt"
)

func TestSetDefaultEnvironment(t *testing.T) {
	hostname, _ := os.Hostname()
	harness := &mockServer{
			Path: fmt.Sprintf("/auth/user/%s/default/", hostname),
			Method: "POST",
			RequestBody: `{"env":"qa"}`,
		}
	settings := harness.createHarness(t)
	test_func := func () error { return settings.Environments.SetDefault("qa") }

	err := test_func()
	if err != nil {
		t.Fatal("SetDefault returned an error")
	}
	if settings.loggedIn != loggedIn {
		t.Errorf("Calls should automatically log the user in")
	}
	if settings.Environments.Default() != "qa" {
		t.Errorf("Default environment should be updated")
	}

	harness.testBadResultFailsGracefully(test_func)
	harness.testCallFailsWhenLoggedOut(test_func)
}

func TestCreateEnvironment(t *testing.T) {
	harness := &mockServer{
		Path: fmt.Sprintf("/auth/env/qa/"),
		Method: "POST",
	}
	settings := harness.createHarness(t)
	test_func := func () error { return settings.Environments.Create("qa") }

	err := test_func()
	if err != nil {
		t.Fatal("CreateEnvironment returned an error")
	}
	if settings.loggedIn != loggedIn {
		t.Errorf("Calls should automatically log the user in")
	}

	harness.testBadResultFailsGracefully(test_func)
	harness.testCallFailsWhenLoggedOut(test_func)
}
