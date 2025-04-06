# Website Status Monitor

A Go program that monitors website availability and provides real-time status updates.

## Features

- **Real-time Monitoring**: Continuously checks website availability
- **Status Reporting**: Provides detailed status information for each monitored website
- **Configurable Check Interval**: Adjustable monitoring frequency
- **Concurrent Checks**: Efficiently monitors multiple websites simultaneously
- **Detailed Status Information**: Includes response time and status codes

## How to Use

1. Start the monitor:
   ```bash
   go run website_monitor.go
   ```

2. The program will:
   - Check website availability at regular intervals
   - Display status updates in real-time
   - Show response times and status codes
   - Indicate if websites are up or down

## Status Indicators

- **UP**: Website is accessible and responding
- **DOWN**: Website is not accessible or not responding
- **Response Time**: Time taken to receive a response
- **Status Code**: HTTP status code returned by the server

## Technical Details

- Written in Go
- Uses concurrent goroutines for efficient monitoring
- Implements timeout handling for slow responses
- Provides detailed error reporting

## Dependencies

- Go standard library
- No external dependencies required

## Notes

- The monitor uses HTTP GET requests to check website availability
- Timeout is set to 10 seconds for each request
- Status codes 200-299 are considered successful
- All other status codes indicate potential issues

## Example Output

```
Website Status Monitor
=====================
Checking websites every 30 seconds...
Press Ctrl+C to stop

[2024-04-05 20:15:00] Checking https://example.com
[2024-04-05 20:15:01] https://example.com is UP (200 OK) - Response time: 1.2s
[2024-04-05 20:15:00] Checking https://test.com
[2024-04-05 20:15:02] https://test.com is DOWN - Connection timeout
```

## Error Handling

The program handles various types of errors:
- Connection timeouts
- DNS resolution failures
- Invalid URLs
- Network errors
- HTTP errors

Each error is logged with a timestamp and detailed information for troubleshooting. 