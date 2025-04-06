package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// SystemInfo represents system and network information
type SystemInfo struct {
	IPAddress   string
	PingLatency string
	LastCheck   time.Time
}

// Service represents a cloud service status
type Service struct {
	Name      string
	Provider  string
	Status    string
	LastCheck time.Time
	Details   string
	Region    string
}

// Provider represents a cloud provider
type Provider struct {
	Name     string
	URL      string
	Services []Service
}

// StatusReport represents the overall status report
type StatusReport struct {
	Timestamp  time.Time
	SystemInfo SystemInfo
	Providers  []Provider
}

// HTML template for the status report
const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Hyperscaler Status Dashboard</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 20px;
            background-color: #f5f5f5;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
        }
        .header {
            background-color: #333;
            color: white;
            padding: 20px;
            border-radius: 5px;
            margin-bottom: 20px;
        }
        .system-info {
            background-color: #fff;
            padding: 15px;
            border-radius: 5px;
            margin-bottom: 20px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .system-info h2 {
            margin-top: 0;
            color: #333;
        }
        .system-info p {
            margin: 5px 0;
        }
        .status-table {
            width: 100%;
            border-collapse: collapse;
            background-color: white;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            margin-bottom: 20px;
        }
        .status-table th, .status-table td {
            padding: 12px;
            text-align: left;
            border-bottom: 1px solid #ddd;
        }
        .status-table th {
            background-color: #f8f9fa;
            font-weight: bold;
        }
        .status-table tr:hover {
            background-color: #f5f5f5;
        }
        .status-up {
            color: #4caf50;
            font-weight: bold;
        }
        .status-down {
            color: #f44336;
            font-weight: bold;
        }
        .timestamp {
            color: #666;
            font-size: 14px;
        }
        .refresh-button {
            background-color: #4caf50;
            color: white;
            padding: 10px 20px;
            border: none;
            border-radius: 5px;
            cursor: pointer;
            margin-top: 20px;
        }
        .refresh-button:hover {
            background-color: #45a049;
        }
        .provider-header {
            background-color: #e9ecef;
            padding: 10px;
            margin-top: 20px;
            border-radius: 5px;
            font-size: 18px;
            font-weight: bold;
        }
        .provider-link {
            color: #0066cc;
            text-decoration: none;
            font-size: 14px;
            margin-left: 10px;
        }
        .provider-link:hover {
            text-decoration: underline;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Hyperscaler Status Dashboard</h1>
            <p class="timestamp">Last Updated: {{.Timestamp.Format "2006-01-02 15:04:05 MST"}}</p>
        </div>
        <div class="system-info">
            <h2>System Information</h2>
            <p><strong>IP Address:</strong> {{.SystemInfo.IPAddress}}</p>
            <p><strong>Google Ping Latency:</strong> {{.SystemInfo.PingLatency}}</p>
            <p><strong>Last Network Check:</strong> {{.SystemInfo.LastCheck.Format "15:04:05 MST"}}</p>
        </div>
        {{range .Providers}}
        <div class="provider-header">
            {{.Name}}
            <a href="{{.URL}}" target="_blank" class="provider-link">View Official Status Page</a>
        </div>
        <table class="status-table">
            <thead>
                <tr>
                    <th>Service</th>
                    <th>Status</th>
                    <th>Region</th>
                    <th>Details</th>
                    <th>Last Check</th>
                </tr>
            </thead>
            <tbody>
                {{range .Services}}
                <tr>
                    <td>{{.Name}}</td>
                    <td class="{{if eq .Status "UP"}}status-up{{else}}status-down{{end}}">{{.Status}}</td>
                    <td>{{.Region}}</td>
                    <td>{{.Details}}</td>
                    <td>{{.LastCheck.Format "15:04:05 MST"}}</td>
                </tr>
                {{end}}
            </tbody>
        </table>
        {{end}}
        <button class="refresh-button" onclick="location.reload()">Refresh Now</button>
    </div>
</body>
</html>
`

// getSystemInfo gets IP address and ping latency
func getSystemInfo() SystemInfo {
	// Get IP address
	resp, err := http.Get("https://api.ipify.org?format=text")
	var ipAddress string
	if err == nil {
		ip, _ := io.ReadAll(resp.Body)
		ipAddress = string(ip)
		resp.Body.Close()
	} else {
		ipAddress = "Unable to determine"
	}

	// Get ping latency
	cmd := exec.Command("ping", "-c", "1", "www.google.com")
	output, err := cmd.Output()
	var pingLatency string
	if err == nil {
		outputStr := string(output)
		if strings.Contains(outputStr, "time=") {
			start := strings.Index(outputStr, "time=") + 5
			end := strings.Index(outputStr[start:], " ") + start
			pingLatency = outputStr[start:end]
		} else {
			pingLatency = "Unable to determine"
		}
	} else {
		pingLatency = "Unable to determine"
	}

	return SystemInfo{
		IPAddress:   ipAddress,
		PingLatency: pingLatency,
		LastCheck:   time.Now(),
	}
}

// checkGoogleStatus checks Google Workspace status
func checkGoogleStatus(services *[]Service) error {
	resp, err := http.Get("https://www.google.com/appsstatus/dashboard/")
	if err != nil {
		return fmt.Errorf("error checking Google status: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading Google response: %v", err)
	}

	content := string(body)
	// Check for "No incidents" text
	if strings.Contains(content, "No incidents") {
		*services = append(*services, Service{
			Name:      "Google Workspace",
			Provider:  "Google",
			Status:    "UP",
			LastCheck: time.Now(),
			Details:   "All services operational",
			Region:    "Global",
		})
	} else {
		*services = append(*services, Service{
			Name:      "Google Workspace",
			Provider:  "Google",
			Status:    "DOWN",
			LastCheck: time.Now(),
			Details:   "Service issues detected",
			Region:    "Global",
		})
	}
	return nil
}

// checkOCStatus checks Oracle Cloud status
func checkOCStatus(services *[]Service) error {
	resp, err := http.Get("https://ocistatus.oraclecloud.com/#/")
	if err != nil {
		return fmt.Errorf("error checking OCI status: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading OCI response: %v", err)
	}

	content := string(body)
	if strings.Contains(content, "Service Disruption") || strings.Contains(content, "Service Outage") {
		*services = append(*services, Service{
			Name:      "Oracle Cloud Infrastructure",
			Provider:  "Oracle Cloud",
			Status:    "DOWN",
			LastCheck: time.Now(),
			Details:   "Service disruption or outage detected",
			Region:    "Global",
		})
	} else {
		*services = append(*services, Service{
			Name:      "Oracle Cloud Infrastructure",
			Provider:  "Oracle Cloud",
			Status:    "UP",
			LastCheck: time.Now(),
			Details:   "All services operational",
			Region:    "Global",
		})
	}
	return nil
}

// checkAzureStatus checks Azure status
func checkAzureStatus(services *[]Service) error {
	resp, err := http.Get("https://status.azure.com/en-us/status")
	if err != nil {
		return fmt.Errorf("error checking Azure status: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading Azure response: %v", err)
	}

	content := string(body)
	if strings.Contains(content, "Service degradation") || strings.Contains(content, "Service disruption") {
		*services = append(*services, Service{
			Name:      "Azure Services",
			Provider:  "Azure",
			Status:    "DOWN",
			LastCheck: time.Now(),
			Details:   "Service degradation or disruption detected",
			Region:    "Global",
		})
	} else {
		*services = append(*services, Service{
			Name:      "Azure Services",
			Provider:  "Azure",
			Status:    "UP",
			LastCheck: time.Now(),
			Details:   "All services operational",
			Region:    "Global",
		})
	}
	return nil
}

func generateReport() StatusReport {
	var wg sync.WaitGroup
	var services []Service

	// Get system information
	systemInfo := getSystemInfo()

	// Check all providers concurrently
	wg.Add(3)
	go func() {
		defer wg.Done()
		if err := checkGoogleStatus(&services); err != nil {
			log.Printf("Error checking Google: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := checkOCStatus(&services); err != nil {
			log.Printf("Error checking OCI: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := checkAzureStatus(&services); err != nil {
			log.Printf("Error checking Azure: %v", err)
		}
	}()

	wg.Wait()

	// Organize services by provider
	providers := make(map[string][]Service)
	for _, service := range services {
		providers[service.Provider] = append(providers[service.Provider], service)
	}

	var reportProviders []Provider
	for name, services := range providers {
		var url string
		switch name {
		case "Google":
			url = "https://www.google.com/appsstatus/dashboard/"
		case "Oracle Cloud":
			url = "https://ocistatus.oraclecloud.com/#/"
		case "Azure":
			url = "https://status.azure.com/en-us/status"
		}
		reportProviders = append(reportProviders, Provider{
			Name:     name,
			URL:      url,
			Services: services,
		})
	}

	// Set timezone to Pacific
	loc, _ := time.LoadLocation("America/Los_Angeles")
	return StatusReport{
		Timestamp:  time.Now().In(loc),
		SystemInfo: systemInfo,
		Providers:  reportProviders,
	}
}

func main() {
	// Create templates directory if it doesn't exist
	if err := os.MkdirAll("templates", 0755); err != nil {
		log.Fatal(err)
	}

	// Create the HTML template file
	tmpl, err := template.New("status").Parse(htmlTemplate)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting Hyperscaler Status Monitor...")
	fmt.Println("Generating reports every 5 minutes...")
	fmt.Println("Press Ctrl+C to stop")

	for {
		// Generate report
		report := generateReport()

		// Create output file
		file, err := os.Create("templates/status.html")
		if err != nil {
			log.Printf("Error creating file: %v", err)
			continue
		}

		// Execute template
		if err := tmpl.Execute(file, report); err != nil {
			log.Printf("Error executing template: %v", err)
			file.Close()
			continue
		}

		file.Close()
		fmt.Printf("Report generated at %s\n", report.Timestamp.Format(time.RFC3339))

		// Wait for 5 minutes
		time.Sleep(5 * time.Minute)
	}
}
