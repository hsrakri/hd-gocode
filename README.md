# Website Status Monitor

A Go program that monitors the status of major websites using multiple status checking services. This tool provides real-time status updates for popular websites by checking both Downdetector and IsItDownRightNow services.

## Features

- Monitors multiple major websites simultaneously
- Checks status from two different monitoring services:
  - Downdetector.com
  - IsItDownRightNow.com
- Color-coded status display (green for UP, red for DOWN)
- Detailed status information for each website
- Automatic updates every 5 minutes
- Real-time console output with timestamps

## Currently Monitored Websites

- Uber (uber.com)
- Netflix (netflix.com)
- Amazon (amazon.com)
- Google (google.com)
- Facebook (facebook.com)

## Prerequisites

- Go 1.16 or higher installed on your system
- Internet connection to check website statuses

## Installation

1. Clone the repository:
```bash
git clone https://github.com/hsrakri/hd-gocode.git
cd hd-gocode
```

2. Run the program:
```bash
go run website_monitor.go
```

## Usage

The program will start monitoring the websites immediately after running. You'll see output like this:

```
Starting Website Status Monitor...
Press Ctrl+C to stop

Checking sites at 2024-04-06T12:00:00Z

==================================================
Website: Uber
Status: UP
Details: No issues reported on Downdetector
Last Check: 2024-04-06T12:00:00Z
URL: https://www.uber.com
==================================================
```

- The program runs continuously until you stop it with Ctrl+C
- Status updates occur every 5 minutes
- Each website's status is checked against both monitoring services
- Color coding helps quickly identify status:
  - Green: Website is UP
  - Red: Website is DOWN

## Program Structure

- `Site` struct: Represents a website to monitor
  - Name: Website name
  - URL: Website URL
  - Status: Current status (UP/DOWN)
  - LastCheck: Timestamp of last check
  - Details: Detailed status information

- Main functions:
  - `checkDowndetector`: Checks status on downdetector.com
  - `checkIsItDown`: Checks status on isitdownrightnow.com
  - `printStatus`: Displays formatted status information

## Adding New Websites

To add a new website to monitor, add a new entry to the `sites` slice in the `main` function:

```go
sites := []Site{
    {Name: "NewWebsite", URL: "https://www.newwebsite.com"},
    // ... existing sites ...
}
```

## Error Handling

The program includes error handling for:
- Network connection issues
- Invalid responses from monitoring services
- Failed status checks

## Contributing

Feel free to contribute to this project by:
1. Adding more websites to monitor
2. Implementing additional monitoring services
3. Adding new features or improvements
4. Fixing bugs

## License

This project is open source and available under the MIT License.

## Author

Created by Haarith D 