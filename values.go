package cityhall

import (
	"time"
	"strings"
	"fmt"
)

type Value struct {
	Value string
	Protect bool
}

type Entry struct {
	Id int
	Name string
	Value string
	Author string
	DateTime time.Time
	Active bool
	Protect bool
	Override string
}

type History struct {
	Entries []Entry
}

type Child struct {
	Id int
	Name string
	Override string
	Path string
	Value string
	Protect bool
}

type Children struct {
	Path string
	SubChildren []Child
}

type Values struct {
	parent *Settings
}

func sanitizePath(path string) string {
	if len(path) == 0 {
		return "/"
	}

	var ret string
	ret = path
	if !strings.HasPrefix(path, "/") {
		ret = "/" + path
	}
	if !strings.HasSuffix(path, "/") {
		ret = ret + "/"
	}
	return ret
}

func (v *Values) GetRaw(environment string, path string, args map[string]string) (string, error) {
	get_url := fmt.Sprintf("%s/env/%s%s", v.parent.Url, environment, sanitizePath(path))

	if len(args) > 0 {
		get_url = get_url + "?"

		for key, value := range args {
			get_url = fmt.Sprintf("%s%s=%s&", get_url, key, value)
		}
	}

	env_bytes, err := v.parent.createCall("GET", get_url, "")

	if err != nil {
		return "", err
	}

	return string(env_bytes), nil
}

func (v *Values) SetRaw(environment string, path string, value Value, override string) error {
	return nil
}

func (v *Values) SetValue(environment string, path string, value string, override string) error {
	return nil
}

func (v *Values) SetProtect(environment string, path string, protect bool, override string) error {
	return nil
}

func (v *Values) Delete(environment string, path string, override string) error {
	return nil
}

func (v *Values) GetHistory(environment string, path string, override string) (History, error) {
	return History{}, nil
}

func (v *Values) GetChildren(environment string, path string) (Children, error) {
	return Children{}, nil
}
