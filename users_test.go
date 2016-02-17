package cityhall

import (
	"testing"
)

func TestGetUser(t *testing.T) {
	harness := &mockServer{
		Path: "/auth/user/test_user/",
		Method: "GET",
		Body: `{
			"Response": "Ok",
			"Rights": {
            	"dev": 4,
            	"auto": 1,
				"users": 1
			}
        }`,
	}
	settings := harness.createHarness(t)
	test_func := func () error { _, err := settings.Users.Get("test_user"); return err }

	user, err := settings.Users.Get("test_user")
	if err != nil {
		t.Fatal("GetUser returned an error")
	}
	if len(user.Rights) != 3 {
		t.Fatal("Response does not match expected.")
	}
	if settings.loggedIn != loggedIn {
		t.Errorf("Calls should automatically log the user in")
	}

	harness.testBadResultFailsGracefully(test_func)
	harness.testCallFailsWhenLoggedOut(test_func)
}

func TestCreateUser(t *testing.T) {
	harness := &mockServer{
		Path: "/auth/user/a_new_user/",
		Method: "POST",
		RequestBody: `{"passhash": ""}`,
	}
	settings := harness.createHarness(t)
	test_func := func () error { return settings.Users.CreateUser("a_new_user", "") }

	err := test_func()
	if err != nil {
		t.Fatal("CreateUser returned an error")
	}
	if settings.loggedIn != loggedIn {
		t.Errorf("Calls should automatically log the user in")
	}

	harness.testBadResultFailsGracefully(test_func)
	harness.testCallFailsWhenLoggedOut(test_func)
}

func TestDeleteUser(t *testing.T) {
	harness := &mockServer{
		Path: "/auth/user/a_new_user/",
		Method: "DELETE",
	}
	settings := harness.createHarness(t)
	test_func := func () error { return settings.Users.DeleteUser("a_new_user") }

	err := test_func()
	if err != nil {
		t.Fatal("DeleteUser returned an error")
	}
	if settings.loggedIn != loggedIn {
		t.Errorf("Calls should automatically log the user in")
	}

	harness.testBadResultFailsGracefully(test_func)
	harness.testCallFailsWhenLoggedOut(test_func)
}

func TestGrantUser(t *testing.T) {
	harness := &mockServer{
		Path: "/auth/grant/",
		Method: "POST",
		RequestBody: `{"env": "dev", "user": "a_new_user", "rights": 2}`,
	}
	settings := harness.createHarness(t)
	test_func := func () error { return settings.Users.Grant("a_new_user", "dev", ReadProtected) }

	err := test_func()
	if err != nil {
		t.Fatal("GrantUser returned an error")
	}
	if settings.loggedIn != loggedIn {
		t.Errorf("Calls should automatically log the user in")
	}

	harness.testBadResultFailsGracefully(test_func)
	harness.testCallFailnLoggedOut(test_func)
}
