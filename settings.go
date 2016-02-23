package cityhall

import (
	"net/http"
	"sync"
	"net/http/cookiejar"
	"os"
	"fmt"
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
	loggedIn state
	syncObject sync.RWMutex
	cookieJar *cookiejar.Jar
	httpClient *http.Client
}

func NewSettings(info CityHallInfo) (*Settings, error) {
	var username, hostname string
	var err error

	if info.Username == "" {
		if hostname, err = os.Hostname(); err != nil {
			return nil, err
		}
		username = hostname
	} else {
		username = info.Username
	}

	settings := &Settings{
		Url: info.Url,
		username: username,
		passhash: hash(info.Password),
		loggedIn: notYetLoggedIn,
	}

	if settings.cookieJar, err = cookiejar.New(nil); err != nil {
		return nil, err
	}
	settings.httpClient = &http.Client{
		Jar: settings.cookieJar,
	}
	settings.Environments.parent = settings
	settings.Users.parent = settings
	settings.Values.parent = settings

	return settings, nil
}

func (s *Settings) GetValueFull(path string, environment string, override string) (string, error) {
	args := make(map[string]string)
	args["override"] = override
	json, err := s.Values.GetRaw(environment, path, args)
	if err != nil {
		return "", err
	}
	val, err_val := getValueFromResponse([]byte(json))
	if err_val != nil {
		return "", err_val
	}
	return val, nil
}

func (s *Settings) GetValue(path string) (string, error) {
	args := make(map[string]string)
	if err_login := s.ensureLoggedIn(); err_login != nil {
		return "", err_login
	}
	json, err := s.Values.GetRaw(s.Environments.Default(), path, args)
	if err != nil {
		return "", err
	}
	val, err_val := getValueFromResponse([]byte(json))
	if err_val != nil {
		return "", err_val
	}
	return val, nil
}

func (s *Settings) GetValueEnvironment(path string, environment string) (string, error) {
	args := make(map[string]string)
	json, err := s.Values.GetRaw(environment, path, args)
	if err != nil {
		return "", err
	}
	val, err_val := getValueFromResponse([]byte(json))
	if err_val != nil {
		return "", err_val
	}
	return val, nil
}

func (s *Settings) GetValueOverride(path string, override string) (string, error) {
	if err_login := s.ensureLoggedIn(); err_login != nil {
		return "", err_login
	}
	return s.GetValueFull(path, s.Environments.Default(), override)
}

func (s *Settings) LoggedIn() bool {
	return s.loggedIn == loggedIn
}

func (s *Settings) UpdatePassword(password string) error {
	if err_login := s.ensureLoggedIn(); err_login != nil {
		return err_login
	}

	update_url := fmt.Sprintf("%s/auth/user/%s/", s.Url, s.username)
	update_json := fmt.Sprintf(`{"passhash": "%s"}`, hash(password))
	_, err := s.createCall("PUT", update_url, update_json);
	return err
}

func (s *Settings) Logout() error {
	if s.loggedIn != loggedIn {
		return cityhallError("Cannot log out when not already logged in")
	}

	logout_url := fmt.Sprintf("%s/auth/", s.Url)
	_, err := s.createCall("DELETE", logout_url, "")
	if err != nil {
		return err
	}

	s.syncObject.Lock()
	s.loggedIn = loggedOut
	s.syncObject.Unlock()

	return nil
}
