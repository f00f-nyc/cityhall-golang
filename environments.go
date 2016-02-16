package cityhall

import (
	"fmt"
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
	set_url := fmt.Sprintf("%s/auth/user/%s/default/", e.parent.Url, e.parent.username)
	json := fmt.Sprintf(`{"env":"%s"}`, defaultEnvironment)
	if _, err := e.parent.createCall("POST", set_url, json); err != nil {
		return err
	}

	e.parent.syncObject.Lock()
	e.defaultEnvironment = defaultEnvironment
	e.parent.syncObject.Unlock()

	return nil
}

func (e *Environments) Create(environment string) error {
	create_url := fmt.Sprintf("%s/auth/env/%s/", e.parent.Url, environment)

	if _, err := e.parent.createCall("POST", create_url, ""); err != nil {
		return err
	}

	return nil
}

func (e *Environments) Get(environment string) (EnvironmentInfo, error) {
	return EnvironmentInfo{}, nil
}
