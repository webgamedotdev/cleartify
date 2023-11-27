package models

import (
	"fmt"
	"strings"
)

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
	report := make(map[string]interface{})
	for id, msg := range r.Results {
		report[id] = msg
	}
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
