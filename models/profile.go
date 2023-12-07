package models

type Profile struct {
	Name        string    `yaml:"name"`
	Description string    `yaml:"description"`
	Controls    []Control `yaml:"controls"`
}

func NewProfile(controls []Control) *Profile {
	return &Profile{
		Controls: controls,
	}
}
