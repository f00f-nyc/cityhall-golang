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
	if p & None == None {
		return "None"
	} else if p & Read == Read {
		return "Read"
	} else if p & ReadProtected == ReadProtected {
		return "ReadProtected"
	} else if p & Write == Write {
		return "Write"
	} else if p & Grant == Grant {
		return "Grant"
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
