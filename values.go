package cityhall

import (
	"time"
	"strings"
	"fmt"
	"encoding/json"
)

type Value struct {
	Value *string
	Protect *bool
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

func (v *Values) urlFromItems(environment string, path string, args map[string]string) string {
	ret_url := fmt.Sprintf("%s/env/%s%s", v.parent.Url, environment, sanitizePath(path))

	if len(args) > 0 {
		ret_url = ret_url + "?"

		for key, value := range args {
			ret_url = fmt.Sprintf("%s%s=%s&", ret_url, key, value)
		}
	}

	return ret_url
}

func (v *Values) GetRaw(environment string, path string, args map[string]string) (string, error) {
	get_url := v.urlFromItems(environment, path, args)

	env_bytes, err := v.parent.createCall("GET", get_url, "")

	if err != nil {
		return "", err
	}

	return string(env_bytes), nil
}

func (v *Values) SetRaw(environment string, path string, value Value, override string) error {
	args := make(map[string]string)
	args["override"] = override
	post_url := v.urlFromItems(environment, path, args)
	json_str := ""

	if value.Protect == nil && value.Value == nil {
		return cityhallError("Must set Protect and/or Value")
	} else if value.Protect == nil && value.Value != nil {
		json_str = fmt.Sprintf(`{"value": "%s"}`, *value.Value)
	} else if value.Protect != nil && value.Value == nil {
		json_str = fmt.Sprintf(`{"protect": %v}`, *value.Protect)
	} else {
		json_str = fmt.Sprintf(`{"value": "%s", "protect": %v}`, *value.Value, *value.Protect)
	}

	_, err := v.parent.createCall("POST", post_url, json_str)

	return err
}

func (v *Values) SetValue(environment string, path string, value string, override string) error {
	return v.SetRaw(environment, path, Value{Value:&value}, override)
}

func (v *Values) SetProtect(environment string, path string, protect bool, override string) error {
	return v.SetRaw(environment, path, Value{Protect:&protect}, override)
}

func (v *Values) Delete(environment string, path string, override string) error {
	args := make(map[string]string)
	args["override"] = override
	delete_url := v.urlFromItems(environment, path, args)
	_, err := v.parent.createCall("DELETE", delete_url, "")
	return err
}

func (v *Values) GetHistory(environment string, path string, override string) (History, error) {
	args := make(map[string]string)
	args["viewhistory"] = "true"
	args["override"] = override
	json_str, err := v.GetRaw(environment, path, args)
	if err != nil {
		return History{}, err
	}

	//History.History isn't as readable, so this temporary type will be used for unmarshaling
	type history_resp struct {
		History []Entry
	}
	var response history_resp
	if err = json.Unmarshal([]byte(json_str), &response); err != nil {
		return History{}, err
	}
	return History{Entries:response.History}, nil
}

func (v *Values) GetChildren(environment string, path string) (Children, error) {
	args := make(map[string]string)
	args["viewchildren"] = "true"
	json_str, err := v.GetRaw(environment, path, args)
	if err != nil {
		return Children{}, err
	}

	//Children.Children isn't as readable, so this temporary type will be used for unmarshaling
	type children_resp struct {
		Path string
		Children []Child
	}
	var response children_resp
	if err = json.Unmarshal([]byte(json_str), &response); err != nil {
		return Children{}, err
	}

	return Children{Path:response.Path, SubChildren:response.Children}, nil
}
