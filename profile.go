package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/webgamedotdev/cleartify/models"
)

func profileExample() {
	// TODO: use a real target
	var target models.Target
	tmpAccounts := models.NewControl("SV-238196", "Temporary User Account Provisioning").
		AddDescription(models.DescriptionGeneral, "If temporary user accounts remain active when no longer needed...").
		AddDescription(models.DescriptionCheck, "Verify that the Ubuntu operating system expires temporary user accounts within 72 hours...").
		AddDescription(models.DescriptionFix, "If a temporary account must be created, configure the system to terminate the account after a 72-hour time period...").
		SetImpact(0.5).
		Check(createAccountExpirationCheck("tempuser", target)).
		To(models.NotEqual(-1), "Account expiration setting is correct", "Account expiration setting is incorrect")

	fmt.Println("tmpAccounts: \n", tmpAccounts)
}

func getAccountExpiration(accountName string, target models.Target) (interface{}, error) {
	cmdStr := fmt.Sprintf("chage -l %s | grep 'Account expires' | cut -d: -f2", accountName)
	fmt.Println("Running command:", cmdStr)
	out, err := target.ExecuteCommand(cmdStr)
	fmt.Println("Output:", out)
	if err != nil {
		return nil, err
	}

	expirationStr := strings.TrimSpace(out)
	if expirationStr == "never" {
		return -1, nil // -1 can represent 'never expires'
	}

	// Parse expiration date and calculate days until expiration
	daysUntilExpiration, err := parseExpirationDate(expirationStr)
	if err != nil {
		return nil, err
	}
	return daysUntilExpiration, nil
}

func createAccountExpirationCheck(accountName string, target models.Target) models.CheckFunc {
	return func() (interface{}, error) {
		return getAccountExpiration(accountName, target)
	}
}

// fixAccountExpiration sets the expiration time for a given account to 72 hours from now
func fixAccountExpiration(accountName string, target *models.DockerTarget) error {
	expirationDate := time.Now().Add(72 * time.Hour).Format("Jan 02, 2006")
	cmdStr := fmt.Sprintf("sudo chage -E '%s' %s", expirationDate, accountName)
	_, err := target.ExecuteCommand(cmdStr)
	return err
}

// parseExpirationDate parses a date string in a specific format and returns the number of days until that date
func parseExpirationDate(dateStr string) (int, error) {
	// Define the layout of the expected date format
	// This should match the format output by `chage -l`
	// Example layout: "Jan 02, 2006" for a date like "Jul 07, 2023"
	const layout = "Jan 02, 2006"

	expirationDate, err := time.Parse(layout, dateStr)
	if err != nil {
		return 0, fmt.Errorf("error parsing date: %v", err)
	}

	// Calculate the difference in days between now and the expiration date
	daysUntilExpiration := int(expirationDate.Sub(time.Now()).Hours() / 24)

	return daysUntilExpiration, nil
}
