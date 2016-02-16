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

	err := settings.Environments.SetDefault("qa")
	if err != nil {
		t.Fatal("SetDefault returned an error: ")
	}
	if settings.loggedIn != loggedIn {
		t.Errorf("Calls should automatically log the user in")
	}
	if settings.Environments.Default() != "qa" {
		t.Errorf("Default environment should be updated")
	}

	harness.testBadResultFailsGracefully(func () error { return settings.Environments.SetDefault("qa") })
	harness.testCallFailsWhenLoggedOut(func () error { return settings.Environments.SetDefault("qa") })
}
