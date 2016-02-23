package cityhall

import (
	"testing"
)

const response_value_json = `{
		"Response": "Ok",
		"value": "sample value",
		"protect": false
	}`
const response_children_json = `{
		"Response": "Ok",
		"path": "/value1/",
		"children": [
			{
				"id": 302,
				"name": "value1",
				"override": "",
				"path": "/app1/domainA/feature_1/value1/",
				"protect": false,
				"value": "1000"
			},
			{
				"id": 552,
				"name": "value1",
				"override": "cityhall",
				"path": "/app1/domainA/feature_1/value1/",
				"protect": false,
				"value": "2"
			}
					]
	}`
const response_history_json = `{
		"Response": "Ok",
		"History": [
			{
				"active": false,
				"override": "",
				"id": 12,
				"value": "999",
				"datetime": "2015-01-01T01:01:00.000Z",
				"protect": false,
				"name": "value1",
				"author": "cityhall"
			},
			{
				"active": false,
				"override": "",
				"id": 12,
				"value": "1000",
				"datetime": "2015-01-02T01:01:00.000Z",
				"protect": false,
				"name": "value1",
				"author": "cityhall"
			}
		]
	}`

func areEqual(t *testing.T, a string, b string) {
	if a != b {
		t.Fatalf("Expected %s but got %s", a, b)
	}
}

func TestSanitizePath(t *testing.T) {
	areEqual(t, "/", sanitizePath(""))
	areEqual(t, "/", sanitizePath("/"))
	areEqual(t, "/val1/", sanitizePath("val1"))
	areEqual(t, "/val1/", sanitizePath("/val1"))
	areEqual(t, "/val1/", sanitizePath("val1/"))
	areEqual(t, "/val1/", sanitizePath("/val1/"))
	areEqual(t, "/val1/val2/", sanitizePath("val1/val2"))
}

func TestGetRawNoOverride(t *testing.T) {
	harness := &mockServer{
		Path: "/env/dev/value1/",
		Method: "GET",
		Body: response_value_json,
	}
	settings := harness.createHarness(t)
	test_func := func () error { _, err := settings.Values.GetRaw("dev", "value1", make(map[string]string)); return err }

	val, err := settings.Values.GetRaw("dev", "value1", make(map[string]string));
	if err != nil {
		t.Fatal("GetRaw value returned an error")
	}
	if len(val) == 0 {
		t.Fatal("Expected response to contain json")
	}
	if settings.loggedIn != loggedIn {
		t.Errorf("Calls should automatically log the user in")
	}

	harness.testBadResultFailsGracefully(test_func)
	harness.testCallFailsWhenLoggedOut(test_func)
}

func TestGetRawOverride(t *testing.T) {
	harness := &mockServer{
		Path: "/env/qa/value1/?override=cityhall&",
		Method: "GET",
		Body: response_value_json,
	}
	settings := harness.createHarness(t)
	args := make(map[string]string)
	args["override"] = "cityhall"
	test_func := func () error { _, err := settings.Values.GetRaw("qa", "value1", args); return err }

	val, err := settings.Values.GetRaw("qa", "value1", args);
	if err != nil {
		t.Fatal("GetRaw value returned an error")
	}
	if len(val) == 0 {
		t.Fatal("Expected response to contain json")
	}
	if settings.loggedIn != loggedIn {
		t.Errorf("Calls should automatically log the user in")
	}

	harness.testBadResultFailsGracefully(test_func)
	harness.testCallFailsWhenLoggedOut(test_func)
}

func TestGetRawHistoryOverride(t *testing.T) {
	harness := &mockServer{
		Path: "/env/qa/value1/?override=cityhall&viewhistory=true&",
		Method: "GET",
		Body: response_history_json,
	}

	settings := harness.createHarness(t)
	args := make(map[string]string)
	args["override"] = "cityhall"
	args["viewhistory"] = "true"
	test_func := func () error { _, err := settings.Values.GetRaw("qa", "value1", args); return err }

	val, err := settings.Values.GetRaw("qa", "value1", args);
	if err != nil {
		t.Fatal("GetRaw value returned an error")
	}
	if len(val) == 0 {
		t.Fatal("Expected response to contain json")
	}
	if settings.loggedIn != loggedIn {
		t.Errorf("Calls should automatically log the user in")
	}

	harness.testBadResultFailsGracefully(test_func)
	harness.testCallFailsWhenLoggedOut(test_func)
}

func TestGetRawChildren(t *testing.T) {
	harness := &mockServer{
		Path: "/env/qa/app1/domainA/feature_1/?viewchildren=true&",
		Method: "GET",
		Body: response_children_json,
	}
	settings := harness.createHarness(t)
	args := make(map[string]string)
	args["viewchildren"] = "true"
	test_func := func () error { _, err := settings.Values.GetRaw("qa", "app1/domainA/feature_1/", args); return err }

	val, err := settings.Values.GetRaw("qa", "/app1/domainA/feature_1", args);
	if err != nil {
		t.Fatal("GetRaw value returned an error")
	}
	if len(val) == 0 {
		t.Fatal("Expected response to contain json")
	}
	if settings.loggedIn != loggedIn {
		t.Errorf("Calls should automatically log the user in")
	}

	harness.testBadResultFailsGracefully(test_func)
	harness.testCallFailsWhenLoggedOut(test_func)
}

func TestSetRawRequiresThingToSet(t *testing.T) {
	settings := (&mockServer{}).createHarness(t)
	err := settings.Values.SetRaw("dev", "value1", Value{Value: nil, Protect: nil}, "")
	if err == nil {
		t.Fatal("Expected an error from SetRaw(), are attempting to execute a no-op")
	}
}

func TestSetRawValueNoOverride(t *testing.T) {
	harness := &mockServer{
		Path: "/env/dev/value1/?override=",
		Method: "POST",
		RequestBody: `{"value": "some value"}`,
	}
	settings := harness.createHarness(t)
	s := "some value"
	value := Value{Value: &s}
	test_func := func () error { return settings.Values.SetRaw("dev", "value1", value, "")}

	err := test_func()
	if err != nil {
		t.Fatal("SetRaw value returned an error")
	}
	if settings.loggedIn != loggedIn {
		t.Errorf("Calls should automatically log the user in")
	}

	harness.testBadResultFailsGracefully(test_func)
	harness.testCallFailsWhenLoggedOut(test_func)
}

func TestSetRawProtectNoOverride(t *testing.T) {
	harness := &mockServer{
		Path: "/env/dev/value1/?override=",
		Method: "POST",
		RequestBody: `{"protect": true}`,
	}
	settings := harness.createHarness(t)
	p := true
	value := Value{Protect: &p}
	test_func := func () error { return settings.Values.SetRaw("dev", "value1", value, "")}

	err := test_func()
	if err != nil {
		t.Fatal("SetRaw value returned an error")
	}
	if settings.loggedIn != loggedIn {
		t.Errorf("Calls should automatically log the user in")
	}

	harness.testBadResultFailsGracefully(test_func)
	harness.testCallFailsWhenLoggedOut(test_func)
}

func TestSetRawValueAndProtectNoOverride(t *testing.T) {
	harness := &mockServer{
		Path: "/env/dev/value1/?override=",
		Method: "POST",
		RequestBody: `{"value": "some value", "protect": true}`,
	}
	settings := harness.createHarness(t)
	s := "some value"
	p := true
	value := Value{Value: &s, Protect: &p}
	test_func := func () error { return settings.Values.SetRaw("dev", "value1", value, "")}

	err := test_func()
	if err != nil {
		t.Fatal("SetRaw value returned an error")
	}
	if settings.loggedIn != loggedIn {
		t.Errorf("Calls should automatically log the user in")
	}

	harness.testBadResultFailsGracefully(test_func)
	harness.testCallFailsWhenLoggedOut(test_func)
}

func TestSetRawWithOverride(t *testing.T) {
	harness := &mockServer{
		Path: "/env/dev/value1/?override=cityhall",
		Method: "POST",
		RequestBody: `{"value": "some value", "protect": true}`,
	}
	settings := harness.createHarness(t)
	s := "some value"
	p := true
	value := Value{Value: &s, Protect: &p}
	test_func := func () error { return settings.Values.SetRaw("dev", "value1", value, "cityhall")}

	err := test_func()
	if err != nil {
		t.Fatal("SetRaw value returned an error")
	}
	if settings.loggedIn != loggedIn {
		t.Errorf("Calls should automatically log the user in")
	}

	harness.testBadResultFailsGracefully(test_func)
	harness.testCallFailsWhenLoggedOut(test_func)
}





