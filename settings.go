package cityhall

import (
	"net/http"
	"sync"
	"net/http/cookiejar"
	"os"
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
		loggedIn: false,
	}

	if settings.cookieJar, err = cookiejar.New(nil); err != nil {
		return nil, err
	}
	settings.httpClient = &http.Client{
		Jar: settings.cookieJar,
	}


	return settings, nil
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
