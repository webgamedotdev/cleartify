package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/webgamedotdev/cleartify/models"
	"gopkg.in/yaml.v2"
)

func parseArgs() (profile, target string, err error) {
	if len(os.Args) < 3 {
		return "", "", fmt.Errorf("Usage: %s <profile> <target>", os.Args[0])
	}
	return os.Args[1], os.Args[2], nil
}

type App struct {
	profile models.Profile
	target  models.Target
}

func NewApp(profile models.Profile, target models.Target) *App {
	return &App{profile, target}
}

func (app *App) RunProfile() (*models.Report, error) {
	report := models.NewReport()
	for _, control := range app.profile.Controls {
		fmt.Println("Running control:", control.ID)
		msg, err := control.Run(app.target)
		if err != nil {
			return &models.Report{}, err
		}
		fmt.Println("Control result:", msg)
		// report.AddResult(control.ID, msg)
	}
	return report, nil
}

func main() {

	profilePath := flag.String("profile", "", "Path to the profile YAML file")
	platform := flag.String("platform", "", "Platform to execute on (e.g., docker://container_id)")
	// reportPath := flag.String("report", "profile.json", "Path to save the report")

	// Parse command line flags
	flag.Parse()

	// Load the profile from the YAML file
	profileData, err := os.ReadFile(*profilePath)
	if err != nil {
		fmt.Printf("Error reading profile file: %v\n", err)
		os.Exit(1)
	}

	var profile models.Profile
	err = yaml.Unmarshal(profileData, &profile)
	if err != nil {
		fmt.Printf("Error parsing profile YAML: %v\n", err)
		os.Exit(1)
	}

	// Determine the target based on the platform
	var target models.Target
	if strings.HasPrefix(*platform, "docker://") {
		containerID := strings.TrimPrefix(*platform, "docker://")
		target = &models.DockerTarget{ContainerID: containerID}
	} else {
		// Handle other platforms or default to local
		target = &models.LocalTarget{}
	}

	// Create and configure the application
	app := NewApp(profile, target)

	// Run the profile and generate the report
	report, err := app.RunProfile()
	if err != nil {
		fmt.Printf("Error running profile: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Profile execution complete. Report: {%s}", report)

	// Save the report to a file
	// err = report.SaveToFile(*reportPath)
	// if err != nil {
	// 	fmt.Printf("Error saving report: %v\n", err)
	// 	os.Exit(1)
	// }
	//
	// fmt.Println("Profile execution complete. Report saved to:", *reportPath)
}
