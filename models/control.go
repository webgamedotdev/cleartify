package models

import "fmt"

type (
	DescriptionType string
	AssertOption    func(*AssertOptions)
	AssertOptions   struct {
		SuccessMsg string
		FailMsg    string
	}
	CheckFunc func() (interface{}, error)
	FixAction func() error
)

var (
	General DescriptionType = "General"
	Check   DescriptionType = "Check"
	Fix     DescriptionType = "Fix"
)

// WithSuccessMsg sets the success message for an assertion.
func WithSuccessMsg(msg string) AssertOption {
	return func(opts *AssertOptions) {
		opts.SuccessMsg = msg
	}
}

// WithFailMsg sets the failure message for an assertion.
func WithFailMsg(msg string) AssertOption {
	return func(opts *AssertOptions) {
		opts.FailMsg = msg
	}
}

type Expected struct {
	Comparison string `yaml:"comparison"`
	Value      string `yaml:"value"`
}

type Verification struct {
	Command  string   `yaml:"command"`
	Expected Expected `yaml:"expected"`
}

type Outcome struct {
	FailMessage    string `yaml:"failMessage"`
	SuccessMessage string `yaml:"successMessage"`
}

type ControlDescriptions struct {
	General string `yaml:"general"`
	Check   string `yaml:"check"`
	Fix     string `yaml:"fix"`
}

type Control struct {
	ID           string              `yaml:"id"`
	Title        string              `yaml:"title"`
	Descriptions ControlDescriptions `yaml:"descriptions"`
	Impact       float64             `yaml:"impact"`
	Tags         map[string]string   `yaml:"tags,omitempty"`
	Verify       Verification        `yaml:"verify"`
	Outcome      Outcome             `yaml:"outcome"`
	CheckFunc    CheckFunc
	Fix          FixAction
}

func NewControl(id, title string) *Control {
	return &Control{
		ID:    id,
		Title: title,
		Tags:  make(map[string]string),
	}
}

// AddDescription adds a description to the control
func (c *Control) AddDescription(descType DescriptionType, description string) *Control {
	switch descType {
	case General:
		c.Descriptions.General = description
	case Check:
		c.Descriptions.Check = description
	case Fix:
		c.Descriptions.Fix = description
	}
	return c
}

// SetImpact sets the impact of the control
func (c *Control) SetImpact(impact float64) *Control {
	c.Impact = impact
	return c
}

// AddTag adds a tag to the control
func (c *Control) AddTag(key, value string) *Control {
	c.Tags[key] = value
	return c
}

func (c *Control) Check(checkFunc CheckFunc) *ExpectationBuilder {
	c.CheckFunc = checkFunc
	return &ExpectationBuilder{control: c}
}

// Fix sets the fix action for the control
func (c *Control) FixAction(fix FixAction) *Control {
	c.Fix = fix
	return c
}

func (c *Control) Run(target Target) (bool, error) {
	var passed bool
	actualValue, err := runCheckFunc(c.Verify.Command, target)
	if err != nil {
		return false, fmt.Errorf("Error running check function: %s", err)
	}

	compareOp := c.Verify.Expected.Comparison
	switch compareOp {
	case "Equal":

	case "NotEqual":
		passed = actualValue != c.Verify.Expected.Value
	default:
		return false, fmt.Errorf("unknown comparison operator: %s", compareOp)
	}

	return passed, nil
}

func runCheckFunc(command string, target Target) (string, error) {
	actualValue, err := target.ExecuteCommand(command)
	if err != nil {
		return "", err
	}
	return actualValue, nil
}
