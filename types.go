package cityhall

import (
	"time"
)

type Permission int

const (
	None Permission = iota
	Read
	ReadProtected
	Write
	Grant
)

func (p Permission) String() string {
	switch p {
		case None: return "None"
		case Read: return "Read"
		case ReadProtected: return "ReadProtected"
		case Write: return "Write"
		case Grant: return "Grant"
	}

	return "Unknown"
}

type Permissions []Permission

type state int

const (
	notYetLoggedIn state = iota
	loggedIn
	loggedOut
)

type EnvironmentRight struct {
	User string
	Permission Permission
}

type EnvironmentRights []EnvironmentRight

type EnvironmentInfo struct {
	Rights EnvironmentRights
}

type UserRight struct {
	Environment string
	Permission Permission
}

type UserRights []UserRight

type UserInfo struct {
	Rights UserRights
}

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

