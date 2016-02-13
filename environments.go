package cityhall

type EnvironmentRights struct {
	User string
	Rights Rights
}

type EnvironmentInfo struct {
	Rights []EnvironmentRights
}

type Environments struct {
	defaultEnvironment string
}

func (e *Environments) Default() string {
	return e.defaultEnvironment
}

func (e *Environments) SetDefault(defaultEnvironment string) error {
	return nil
}

func (e *Environments) Create(environment string) error {
	return nil
}

func (e *Environments) Get(environment string) (EnvironmentInfo, error) {
	return EnvironmentInfo{}, nil
}
