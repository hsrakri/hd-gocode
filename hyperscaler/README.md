# Hyperscaler Status Monitor

A real-time status monitoring tool for major cloud providers (Google, Oracle Cloud, and Azure) with a beautiful web interface.

## Features

- **Real-time Status Monitoring**: Continuously monitors the status of major cloud providers
- **System Information**: Displays your system's IP address and network latency
- **Simple Web Interface**: Modern, responsive design with real-time updates
- **Auto-refresh**: Automatically generates new reports every 5 minutes
- **Manual Refresh**: Option to manually refresh the status at any time
- **Detailed Reports**: Comprehensive status information for each service
- **Concurrent Checks**: Efficiently checks multiple providers simultaneously

## Status Indicators

- **UP**: Service is operational and functioning normally
- **DOWN**: Service is experiencing issues or disruptions

## System Information Displayed

- IP Address: Your system's public IP address
- Ping Latency: Network latency to Google's servers
- Last Network Check: Timestamp of the most recent network check

## Cloud Providers Monitored

1. **Google Workspace**
   - Status: UP/DOWN
   - Region: Global
   - Details: Service operational status

2. **Oracle Cloud Infrastructure**
   - Status: UP/DOWN
   - Region: Global
   - Details: Service operational status

3. **Azure Services**
   - Status: UP/DOWN
   - Region: Global
   - Details: Service operational status

## How to Use

1. Start the monitor:
   ```bash
   go run hyperscaler_monitor.go
   ```

2. Open your web browser and navigate to:
   ```
   http://localhost:8080
   ```

3. The dashboard will show:
   - Current system information
   - Status of all monitored cloud providers
   - Last update timestamp
   - Refresh button for manual updates

## Technical Details

- Written in Go
- Uses concurrent goroutines for efficient status checking
- Implements thread-safe report generation
- Pacific timezone for timestamps
- Responsive web interface with modern CSS

## Dependencies

- Go standard library
- No external dependencies required

## Notes

- The monitor checks provider status pages for keywords indicating service status
- Network latency is measured using ping to Google's servers
- IP address is determined using ipify.org API
- All timestamps are displayed in Pacific timezone 