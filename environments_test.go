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
		Path: "/auth/env/qa/",
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

func TestGetEnvironment(t *testing.T) {
	harness := &mockServer{
		Path: "/auth/env/dev/",
		Method: "GET",
		Body: `{
			"Response": "Ok",
			"Users": {
            	"test_user": 4,
            	"some_user": 1,
            	 "cityhall": 4
			}
        }`,
	}
	settings := harness.createHarness(t)
	test_func := func () error { _, err := settings.Environments.Get("dev"); return err }

	env, err := settings.Environments.Get("dev")
	if err != nil {
		t.Fatal("GetEnvironment returned an error")
	}
	if len(env.Rights) != 3 {
		t.Fatal("Response does not match expected.")
	}
	if settings.loggedIn != loggedIn {
		t.Errorf("Calls should automatically log the user in")
	}

	harness.testBadResultFailsGracefully(test_func)
	harness.testCallFailsWhenLoggedOut(test_func)
}
