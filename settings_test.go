package cityhall

import (
	"testing"
)

func TestServer(t *testing.T) {
	cityhall := (&MockServer{Path:"/auth/", Method:"POST"}).Success(t)
	defer cityhall.Close()

	s, _ := NewSettings(CityHallInfo{Url: cityhall.URL})
	err := s.login()
	if err != nil {
		Log(err.Error())
		t.Fatal("got error back from GetSettings")
	}

	// Log("logged in!")
}

