package cityhall

import (
	"net/http"
	"net/http/httptest"
	"fmt"
	"testing"
)

type MockServer struct {
	Path string
	Method string
	Status_Code int
	Body string
}

func Log(str string) {
	fmt.Printf("%s\n", str)
}

func (m *MockServer) Success(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
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
	}))
}
