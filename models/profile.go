package models

import "fmt"

type Profile struct {
	Controls []Control
}

func NewProfile(controls []Control) *Profile {
	return &Profile{
		Controls: controls,
	}
}

func (p *Profile) Run() (*Report, error) {
	report := NewReport()
	for _, control := range p.Controls {
		fmt.Println("Running control:", control.ID)
		msg, err := control.Run()
		if err != nil {
			return &Report{}, err
		}
		report.AddResult(control.ID, msg)
	}
	return report, nil
}
