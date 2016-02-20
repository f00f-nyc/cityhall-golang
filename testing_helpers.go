package cityhall

import (
	"net/http"
	"net/http/httptest"
	"fmt"
	"testing"
	"io/ioutil"
	"encoding/json"
	"os"
	"strings"
)

type mockServer struct {
	Path string
	Method string
	Status_Code int
	Body string
	RequestBody string
	User string

	CityHall *httptest.Server
	Test *testing.T
	Settings *Settings
}

func log(str string) {
	fmt.Printf("%s\n", str)
}

func logAny(object interface{}) {
	fmt.Printf("%s\n", object)
}

func (m *mockServer) createMockWithLogin(t *testing.T) *httptest.Server {
	if m.User == "" {
		m.User, _ = os.Hostname()
	}

	return httptest.NewServer(
		http.HandlerFunc(
			func (w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/auth/" && r.Method == "POST" {
					//mock logging in.
					request_bytes, err := ioutil.ReadAll(r.Body)
					r.Body.Close()
					if err != nil {
						t.Errorf("Caught an error while attempting to read request bytes: %s", err.Error())
					}

					var request interface{}
					if err = json.Unmarshal(request_bytes, &request); err != nil {
						t.Errorf("Caught an error while attempting to read request json: %s", err.Error())
					}
					m := request.(map[string]interface{})

					if m["username"] == "" {
						t.Errorf("Request to login should include username")
					}
					fmt.Fprintln(w, "{\"Response\": \"Ok\"}")
				} else if r.URL.Path == fmt.Sprintf("/auth/user/%s/default/", m.User) && r.Method == "GET" {
					fmt.Fprintln(w, "{\"Response\": \"Ok\", \"value\":\"dev\"}")
				} else {
					if len(m.Path) > 0 {
						url := m.Path
						url_end := strings.IndexAny(m.Path, "?")
						if url_end >= 0 {
							url = m.Path[0:url_end]
						}
						if url != r.URL.Path {
							t.Errorf("Invalid path.  Expected '%s', but got '%s'", url, r.URL.Path)
						}
						if url_end >= 0 {
							split_1 := func(c rune) bool {
								return c == '&'
							}
							split_2 := func(c rune) bool {
								return c == '='
							}
							url_params := m.Path[url_end+1:]
							params := strings.FieldsFunc(url_params, split_1)
							expected := r.URL.Query()

							if len(expected) != len(params) {
								t.Fatalf("Expected a different number of params")
							}

							for i := 0; i < len(params); i++ {
								param := strings.FieldsFunc(params[i], split_2)
								if val, ok := expected[param[0]]; ok {
									if val[0] != param[1] {
										t.Fatalf("For param %s, expected '%s' but got '%s'", param[0], param[1], val)
									}
								} else {
									t.Fatal("Expected param ", param[0])
								}
							}
						}
					}

					if len(m.Method) > 0 {
						if m.Method != r.Method {
							t.Errorf("Invalid method. Expected '%s', but got '%s'", m.Method, r.Method)
						}
					}

					if len(m.RequestBody) > 0 {
						request_bytes, _ := ioutil.ReadAll(r.Body)
						r.Body.Close()
						if m.RequestBody != string(request_bytes) {
							t.Errorf("Invalid request body to %s: Expected '%s' but got '%s'", r.URL.Path, m.RequestBody, string(request_bytes))
						}
					}

					if len(m.Body) > 0 {
						fmt.Fprintln(w, m.Body)
					} else {
						fmt.Fprintln(w, `{"Response": "Ok"}`)
					}
				}
			}))
}

func (m *mockServer) createHarness(t *testing.T) *Settings {
	var err error
	m.Test = t
	m.CityHall = m.createMockWithLogin(t)
	m.Settings, err = NewSettings(CityHallInfo{Url: m.CityHall.URL})
	if err != nil {
		t.Errorf("Got an error back creating settings")
	}
	return m.Settings
}

func (m *mockServer) testBadResultFailsGracefully(call func() error) {
	m.CityHall.Close()

	// all functions will fail
	m.CityHall = httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, `{"Response": "Failure", "Message": "Unspecified test error"}`)
			}))
	m.Settings.Url = m.CityHall.URL
	err := call()
	if err == nil {
		m.Test.Fatal("testBadResultFailsGracefully: Expected the call to generate an error, but got none")
	}
	m.CityHall.Close()
}

func (m *mockServer) testCallFailsWhenLoggedOut(call func() error) {
	m.CityHall = httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				m.Test.Fatal("Should not have had any calls to server")
				fmt.Fprintln(w, `{"Response": "Failure", "Message": "Unspecified test error"}`)
			}))

	m.Settings.loggedIn = loggedOut
	err := call()
	if err == nil {
		m.Test.Fatal("testCallFailsWhenLoggedOut: Have logged out, and exepected the call to generate an error, but got none")
	}
	m.CityHall.Close()
}
