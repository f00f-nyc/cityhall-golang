package cityhall

type Rights int

const (
	None Rights = iota
	Read
	ReadProtected
	Write
	Grant
)

func (p Rights) String() string {
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


