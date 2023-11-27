package models

import "fmt"

// DescriptionType represents the type of a control description.
type (
	DescriptionType int
	AssertOption    func(*AssertOptions)
	AssertOptions   struct {
		SuccessMsg string
		FailMsg    string
	}
	CheckFunc func() (interface{}, error)
)

// Enum values for DescriptionType.
const (
	DescriptionGeneral DescriptionType = iota // Default value, starts at 0
	DescriptionCheck                          // Automatically becomes 1
	DescriptionFix                            // Automatically becomes 2
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

type Control struct {
	ID           string
	Title        string
	Descriptions map[DescriptionType]string
	Impact       float64
	Tags         map[string]string
	Expectation  *Expectation
	Condition    func() bool
	CheckFunc    CheckFunc
	Fix          *FixAction
}

type FixAction func() error

func NewControl(id, title string) *Control {
	return &Control{
		ID:           id,
		Title:        title,
		Descriptions: make(map[DescriptionType]string),
		Tags:         make(map[string]string),
	}
}

// AddDescription adds a description to the control
func (c *Control) AddDescription(descType DescriptionType, description string) *Control {
	c.Descriptions[descType] = description
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
func (c *Control) FixAction(fixAction *FixAction) *Control {
	c.Fix = fixAction
	return c
}

func (c *Control) Run() (string, error) {
	actualValue, err := c.CheckFunc()
	fmt.Println("Actual value:", actualValue)
	if err != nil {
		return "", err
	}

	passed := c.Expectation.Comparator(actualValue)
	if passed {
		return c.Expectation.SuccessMsg, nil
	} else {
		return c.Expectation.FailMsg, nil
	}
}
