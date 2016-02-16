package cityhall

import (
	"io/ioutil"
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

	resp_raw, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}

	resp_bytes, resp_err := ioutil.ReadAll(resp_raw.Body)
	resp_raw.Body.Close()
	if resp_err != nil {
		return resp_err
	}
	if err = isResponseOkay(resp_bytes); err != nil {
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

func (s *Settings) createCall(method string, url string, body string) ([]byte, error) {
	if err := s.ensureLoggedIn(); err != nil {
		return []byte{}, err
	}

	var json []byte
	if len(body) > 0 {
		json = []byte(body)
	} else {
		json = []byte{}
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json))
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp_raw, req_err := s.httpClient.Do(req)
	if req_err != nil {
		return []byte{}, req_err
	}

	resp_bytes, resp_err := ioutil.ReadAll(resp_raw.Body)
	resp_raw.Body.Close()
	if resp_err != nil {
		return []byte{}, resp_err
	}

	if err := isResponseOkay(resp_bytes); err != nil {
		return []byte{}, err
	}

	return resp_bytes, nil
}
