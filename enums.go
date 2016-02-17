package cityhall

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

type UserRight struct {
	Environment string
	Permission Permission
}

type UserRights []UserRight
