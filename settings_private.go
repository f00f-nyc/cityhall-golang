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
	if err = isResponseOkay(raw_resp.Body); err != nil {
		return err
	}

	s.syncObject.Lock()
	s.loggedIn = loggedIn
	s.syncObject.Unlock()
	return nil
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

func (s *Settings) ensureLoggedIn() error {
	if s.loggedIn == loggedOut {
		return cityhallError(fmt.Sprintf("User %s has already been logged out", s.username))
	} else if s.loggedIn == loggedIn {
		return nil
	} else {
		if err := s.login(); err != nil {
			return err
		}
		if err := s.getDefaultEnvironment(); err != nil {
			return err
		}
	}

	return nil
}
