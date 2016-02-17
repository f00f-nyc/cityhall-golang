package cityhall

type UserInfo struct {
	UserRights Permissions
}

type Users struct {
}

func (u *Users) Get(username string) (UserInfo, error) {
	return 	UserInfo{}, nil
}

func (u *Users) CreateUser(username string, password string) error {
	return nil
}

func (u *Users) DeleteUser(username string) error {
	return nil
}

func (u *Users) Grant(username string, environment string, rights Permission) error {
	return nil
}
