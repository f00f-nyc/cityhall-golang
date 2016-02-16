package cityhall

import (
	"net/http"
	"fmt"
	"bytes"
)

type cityhallError string

func (c cityhallError) Error() string {
	return string(c)
}

func (s *Settings) login() error {
	auth_url := s.Url + "/auth/"
	var json = []byte(fmt.Sprintf(`{"username":"%s", "passhash":"%s"}`, s.username, s.passhash))
	req, _ := http.NewRequest("POST", auth_url, bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json")

	raw_resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	return isResponseOkay(raw_resp.Body)
}

func (s *Settings) getDefaultEnvironment() error {
	default_env_url := s.Url + "/auth/user/" + s.username + "/default/"
	req, _ := http.NewRequest("GET", default_env_url, bytes.NewBuffer([]byte{}))
	req.Header.Set("Content-Type", "application/json")

	raw_resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}

	value, valueOk := getValueFromResponse(raw_resp.Body)
	if valueOk != nil {
		return valueOk
	}
	s.Environments.defaultEnvironment = value
	return nil
}
