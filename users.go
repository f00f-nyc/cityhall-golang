package cityhall

import (
	"fmt"
	"encoding/json"
)

type UserInfo struct {
	Rights UserRights
}

type Users struct {
	parent *Settings
}

func (u *Users) Get(username string) (UserInfo, error) {
	get_url := fmt.Sprintf("%s/auth/user/%s/", u.parent.Url, username)
	env_bytes, err := u.parent.createCall("GET", get_url, "")

	if err != nil {
		return UserInfo{}, err
	}

	var user_resp interface{}
	if err = json.Unmarshal(env_bytes, &user_resp); err != nil {
		return UserInfo{}, err
	}

	user_map := user_resp.(map[string]interface{})
	envs_map := user_map["Rights"].(map[string]interface{})
	var ret UserRights
	for env, rights := range envs_map {
		permission := Permission((rights.(float64)))
		ret = append(ret, UserRight{Environment: env, Permission: permission})
	}

	return UserInfo{Rights: ret}, nil
}

func (u *Users) CreateUser(username string, password string) error {
	create_url := fmt.Sprintf("%s/auth/user/%s/", u.parent.Url, username)
	json := fmt.Sprintf(`{"passhash": "%s"}`, hash(password))

	if _, err := u.parent.createCall("POST", create_url, json); err != nil {
		return err
	}

	return nil
}

func (u *Users) DeleteUser(username string) error {
	create_url := fmt.Sprintf("%s/auth/user/%s/", u.parent.Url, username)

	if _, err := u.parent.createCall("DELETE", create_url, ""); err != nil {
		return err
	}

	return nil
}

func (u *Users) Grant(username string, environment string, rights Permission) error {
	create_url := fmt.Sprintf("%s/auth/grant/", u.parent.Url)
	json := fmt.Sprintf(`{"env": "%s", "user": "%s", "rights": %v}`, environment, username, int(rights))

	if _, err := u.parent.createCall("POST", create_url, json); err != nil {
		return err
	}

	return nil
}
