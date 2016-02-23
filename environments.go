package cityhall

import (
	"fmt"
	"encoding/json"
)

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
	get_url := fmt.Sprintf("%s/auth/env/%s/", e.parent.Url, environment)
	env_bytes, err := e.parent.createCall("GET", get_url, "")

	if err != nil {
		return EnvironmentInfo{}, err
	}

	var env_resp interface{}
	if err = json.Unmarshal(env_bytes, &env_resp); err != nil {
		return EnvironmentInfo{}, err
	}
	env_map := env_resp.(map[string]interface{})
	users_map := env_map["Users"].(map[string]interface{})
	var ret EnvironmentRights
	for user, rights := range users_map {
		permission := Permission((rights.(float64)))
		ret = append(ret, EnvironmentRight{User: user, Permission: permission})
	}

	return EnvironmentInfo{Rights: ret}, nil
}
