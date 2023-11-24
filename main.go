package main

import (
	"fmt"
	"strings"
)

// DescriptionType represents the type of a control description.
type DescriptionType int

// Enum values for DescriptionType.
const (
	DescriptionGeneral DescriptionType = iota // Default value, starts at 0
	DescriptionCheck                          // Automatically becomes 1
	DescriptionFix                            // Automatically becomes 2
)

type Control struct {
	ID           string
	Title        string
	Descriptions map[DescriptionType][]string
	Check        func() bool
	Fix          func()
	Impact       float64
	Tags         map[string]string
}

func NewControl(id, title string) *Control {
	return &Control{
		ID:           id,
		Title:        title,
		Descriptions: make(map[DescriptionType][]string),
		Tags:         make(map[string]string),
	}
}

// AddDescription adds a description of the specified type to the control.
func (c *Control) AddDescription(descriptionType DescriptionType, description string) *Control {
	c.Descriptions[descriptionType] = append(c.Descriptions[descriptionType], description)
	return c
}

func (c *Control) SetCheck(check func() bool) *Control {
	c.Check = check
	return c
}

func (c *Control) SetFix(fix func()) *Control {
	c.Fix = fix
	return c
}

func (c *Control) SetImpact(impact float64) *Control {
	c.Impact = impact
	return c
}

func (c *Control) AddTag(key, value string) *Control {
	c.Tags[key] = value
	return c
}

func (c *Control) AssertVar(v interface{}) *Control {
	c.Check = func() bool {
		return v != nil
	}
	return c
}

func (c *Control) AssertEq(v interface{}) *Control {
	c.Check = func() bool {
		return v == nil
	}
	return c
}

func (c *Control) FailMsg(msg string) *Control {
	c.Check = func() bool {
		fmt.Println(msg)
		return false
	}
	return c
}

func (c *Control) SuccessMsg(msg string) *Control {
	c.Check = func() bool {
		fmt.Println(msg)
		return true
	}
	return c
}

func (c *Control) Run() (string, error) {
	if c.Check == nil {
		return "", fmt.Errorf("Check function is not defined")
	}
	return "", nil
}

type Profile struct {
	Controls []*Control
}

func NewProfile(controls []*Control) *Profile {
	return &Profile{
		Controls: controls,
	}
}

func (p *Profile) Run() (*Report, error) {
	report := NewReport()
	for _, control := range p.Controls {
		msg, err := control.Run()
		if err != nil {
			return nil, err
		}
		report.AddResult(control.ID, msg)
	}
	return report, nil
}

type Report struct {
	Results map[string]string
}

func NewReport() *Report {
	return &Report{
		Results: make(map[string]string),
	}
}

func (r *Report) AddResult(id, msg string) {
	r.Results[id] = msg
}

func (r *Report) Generate(format string) error {
	switch strings.ToLower(format) {
	case "json":
		return r.generateJSON()
	case "yaml":
		return r.generateYAML()
	case "html":
		return r.generateHTML()
	case "markdown":
		return r.generateMarkdown()
	default:
		return fmt.Errorf("Unknown format: %s", format)
	}
}

func (r *Report) generateJSON() error {
	return nil
}

func (r *Report) generateYAML() error {
	return nil
}

func (r *Report) generateHTML() error {
	return nil
}

func (r *Report) generateMarkdown() error {
	return nil
}

func main() {
	tmpAccounts := NewControl("SV-238196", "Temporary User Account Provisioning").
		AddDescription(DescriptionGeneral, "If temporary user accounts remain active when no longer needed...").
		AddDescription(DescriptionCheck, "Verify that the Ubuntu operating system expires temporary user accounts within 72 hours...").
		AddDescription(DescriptionFix, "If a temporary account must be created, configure the system to terminate the account after a 72-hour time period...").
		SetImpact(0.5).
		AddTag("severity", "medium").
		AddTag("category", "authentication").
		AssertVar("test").
		FailMsg("test failed").
		SuccessMsg("test succeeded")

	profile := NewProfile([]*Control{tmpAccounts})
	report, err := profile.Run()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(report)
}
