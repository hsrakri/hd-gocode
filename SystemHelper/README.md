# System Helper

A web-based system monitoring and network diagnostics tool that provides real-time system information and network troubleshooting capabilities. Supports both macOS and Linux operating systems.

## Features

### System Monitoring
- Real-time CPU and Memory usage with cross-platform support
- Top processes with CPU and Memory consumption
- Open network ports monitoring (TCP/UDP)
- Process details and open files inspection
- Auto-refresh every 5 seconds

### Network Diagnostics
- Ping tests for IPv4/IPv6 addresses and domains
- MTR (My TraceRoute) for network path analysis
- DNS lookup and resolution information
- Real-time network statistics

### Cross-Platform Support
- macOS: Uses native commands (ps, vm_stat, netstat)
- Linux: Uses standard Linux tools (top, free, netstat)
- Adaptive parsing for different output formats
- Consistent UI across platforms

## Requirements

- Go 1.16 or later
- Required system tools:
  - macOS:
    - `ps` (built-in)
    - `vm_stat` (built-in)
    - `netstat` (built-in)
    - `ping` (built-in)
    - `mtr` (optional, for network diagnostics)
    - `dig` (for DNS lookups)
    - `lsof` (for process details)
  - Linux:
    - `top`
    - `free`
    - `netstat`
    - `ps`
    - `ping`
    - `mtr` (optional)
    - `dig`
    - `lsof`

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/hd-gocode.git
cd hd-gocode/SystemHelper
```

2. Run the application:
```bash
go run systemhelper.go
```

3. Open your browser and navigate to:
```
http://localhost:8081
```

## Usage

### System Monitoring
- The dashboard automatically updates system statistics every 5 seconds
- View real-time CPU and memory usage graphs
- Monitor top processes sorted by CPU usage
- View active network ports and their states
- Click on process IDs to view detailed information

### Network Diagnostics
1. Enter an IP address or domain name in the input field
2. Click "Run Diagnostics" to perform:
   - Ping test (IPv4/IPv6)
   - MTR analysis (if available)
   - DNS lookup
3. View results in real-time in the respective tabs

## Security Note

This application requires system-level access to run various diagnostic commands. Make sure to:
- Run with appropriate permissions
- Restrict access to trusted users only
- Use in a controlled environment
- Be cautious when exposing port 8081 to external networks

## Contributing

Contributions are welcome! Please feel free to submit pull requests with improvements or bug fixes.

## License

This project is licensed under the MIT License - see the LICENSE file for details. 