package cityhall

import (
	"net/http"
	"net/http/httptest"
	"fmt"
	"testing"
	"io/ioutil"
	"encoding/json"
	"os"
)

type mockServer struct {
	Path string
	Method string
	Status_Code int
	Body string
	RequestBody string
	User string
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
						if m.Path != r.URL.Path {
							t.Errorf("Invalid path.  Expected '%s', but got '%s'", m.Path, r.URL.Path)
						}
					}

					if len(m.Method) > 0 {
						if m.Method != r.Method {
							t.Errorf("Invalid method. Expected '%s', but got '%s'", m.Method, r.Method)
						}
					}

					if len(m.Body) > 0 {
						fmt.Fprintln(w, m.Body)
					} else {
						fmt.Fprintln(w, "{\"Response\": \"Ok\"}")
					}
				}
			}))
}

func (m *mockServer) createSettingsAndServer(t *testing.T) (*Settings, *httptest.Server) {
	cityhall := m.createMockWithLogin(t)
	s, err := NewSettings(CityHallInfo{Url: cityhall.URL})
	if err != nil {
		t.Errorf("Got an error back creating settings")
	}
	return s, cityhall
}
