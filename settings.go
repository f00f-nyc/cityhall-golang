package cityhall

import (
	"net/http"
	"io/ioutil"
	"sync"
	"net/url"
	"encoding/json"
	"io"
)

type CityHallInfo struct {
	Url string
	Username string
	Password string
}

type Settings struct {
	Url string
	Users Users
	Values Values
	Environments Environments

	username string
	passhash string
	loggedIn bool
	syncObject sync.RWMutex
}

func NewSettings(info CityHallInfo) (*Settings, error) {
	return &Settings{
		Url: info.Url,
		username: info.Username,
		passhash: info.Password,
		loggedIn: false,
	}, nil
}

func (s *Settings) GetValueFull(path string, environment string, override string) (string, error) {
	return "", nil
}

func (s *Settings) GetValue(path string) (string, error) {
	return s.GetValueFull(path, "", "")
}

func (s *Settings) GetValueEnvironment(path string, environment string) (string, error) {
	return s.GetValueFull(path, environment, "")
}

func (s *Settings) GetValueOverride(path string, override string) (string, error) {
	return s.GetValueFull(path, "", override)
}

func (s *Settings) LoggedIn() bool {
	return s.loggedIn
}

func (s *Settings) UpdatePassword(password string) error {
	return nil
}

func (s *Settings) Logout() error {
	return nil
}

///////////////////////////////////////
type cityhallError string

func (c cityhallError) Error() string {
	return string(c)
}

func isOkay(body io.ReadCloser) error {
	resp_bytes, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	type response struct {
		Response, Message string
	}
	var resp response
	if err = json.Unmarshal(resp_bytes, &resp); err != nil {
		return err
	}

	if resp.Response != "Ok" {
		err = cityhallError(resp.Message)
	}
	return err
}

func (s *Settings) login() error {
	auth_url := s.Url + "/auth/"
	raw_resp, err := http.PostForm(auth_url, url.Values{"username": {s.username}, "passhash": {s.passhash}})
	if err != nil {
		return err
	}
//	if resp, err := ioutil.ReadAll(raw_resp.Body); err != nil {
//		return err
//	}
	return isOkay(raw_resp.Body)
}
