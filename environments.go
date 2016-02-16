package cityhall

import (
	"fmt"
	"net/http"
	"bytes"
)

type EnvironmentRights struct {
	User string
	Rights Rights
}

type EnvironmentInfo struct {
	Rights []EnvironmentRights
}

type Environments struct {
	defaultEnvironment string
	parent *Settings
}

func (e *Environments) Default() string {
	return e.defaultEnvironment
}

func (e *Environments) SetDefault(defaultEnvironment string) error {
	if err := e.parent.ensureLoggedIn(); err != nil {
		return err
	}

	set_url := fmt.Sprintf("%s/auth/user/%s/default/", e.parent.Url, e.parent.username)
	var json = []byte(fmt.Sprintf(`{"env":"%s"}`, defaultEnvironment))
	req, _ := http.NewRequest("POST", set_url, bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json")

	raw_resp, err := e.parent.httpClient.Do(req)
	if err != nil {
		return err
	}
	if err = isResponseOkay(raw_resp.Body); err != nil {
		return err
	}

	e.parent.syncObject.Lock()
	e.defaultEnvironment = defaultEnvironment
	e.parent.syncObject.Unlock()

	return nil
}

func (e *Environments) Create(environment string) error {
	return nil
}

func (e *Environments) Get(environment string) (EnvironmentInfo, error) {
	return EnvironmentInfo{}, nil
}
