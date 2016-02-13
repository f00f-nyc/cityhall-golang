package cityhall

import (
	"time"
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

}

func (v *Values) GetRaw(environment string, path string, args map[string]string) (string, error) {
	return "", nil
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
