package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// Site represents a website to monitor
type Site struct {
	Name      string
	URL       string
	Status    string
	LastCheck time.Time
	Details   string
}

// checkDowndetector checks if a site is reported down on downdetector.com
func checkDowndetector(site *Site) error {
	url := fmt.Sprintf("https://downdetector.com/status/%s/", strings.ToLower(site.Name))
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error checking downdetector: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}

	// Simple check for common error indicators
	content := string(body)
	if strings.Contains(content, "reported problems") || strings.Contains(content, "issues detected") {
		site.Status = "DOWN"
		site.Details = "Issues reported on Downdetector"
	} else {
		site.Status = "UP"
		site.Details = "No issues reported on Downdetector"
	}
	site.LastCheck = time.Now()
	return nil
}

// checkIsItDown checks if a site is reported down on isitdownrightnow.com
func checkIsItDown(site *Site) error {
	url := fmt.Sprintf("https://www.isitdownrightnow.com/%s.html", strings.ToLower(site.Name))
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error checking isitdownrightnow: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}

	// Simple check for common error indicators
	content := string(body)
	if strings.Contains(content, "is down") || strings.Contains(content, "has issues") {
		site.Status = "DOWN"
		site.Details = "Issues reported on IsItDownRightNow"
	} else {
		site.Status = "UP"
		site.Details = "No issues reported on IsItDownRightNow"
	}
	site.LastCheck = time.Now()
	return nil
}

func printStatus(site Site) {
	statusColor := "\033[32m" // Green for UP
	if site.Status == "DOWN" {
		statusColor = "\033[31m" // Red for DOWN
	}
	resetColor := "\033[0m"

	fmt.Printf("\n%s%s%s\n", statusColor, strings.Repeat("=", 50), resetColor)
	fmt.Printf("Website: %s\n", site.Name)
	fmt.Printf("Status: %s%s%s\n", statusColor, site.Status, resetColor)
	fmt.Printf("Details: %s\n", site.Details)
	fmt.Printf("Last Check: %s\n", site.LastCheck.Format(time.RFC3339))
	fmt.Printf("URL: %s\n", site.URL)
	fmt.Printf("%s%s%s\n", statusColor, strings.Repeat("=", 50), resetColor)
}

func main() {
	// Sites to monitor
	sites := []Site{
		{Name: "Uber", URL: "https://www.uber.com"},
		{Name: "Netflix", URL: "https://www.netflix.com"},
		{Name: "Amazon", URL: "https://www.amazon.com"},
		{Name: "Google", URL: "https://www.google.com"},
		{Name: "Facebook", URL: "https://www.facebook.com"},
	}

	fmt.Println("Starting Website Status Monitor...")
	fmt.Println("Press Ctrl+C to stop\n")

	// Monitor sites
	for {
		fmt.Printf("\nChecking sites at %s\n", time.Now().Format(time.RFC3339))

		for i := range sites {
			// Check both services
			err1 := checkDowndetector(&sites[i])
			err2 := checkIsItDown(&sites[i])

			if err1 != nil || err2 != nil {
				log.Printf("Error checking %s: %v, %v", sites[i].Name, err1, err2)
				continue
			}

			// Print status
			printStatus(sites[i])
		}

		fmt.Printf("\nWaiting 5 minutes before next check...\n")
		time.Sleep(5 * time.Minute)
	}
}
