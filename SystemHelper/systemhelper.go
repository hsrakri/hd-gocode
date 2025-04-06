package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// SystemStats represents system statistics
type SystemStats struct {
	CPUUsage     float64
	MemoryUsage  float64
	TopProcesses []ProcessInfo
	OpenPorts    []PortInfo
}

// ProcessInfo represents process information
type ProcessInfo struct {
	PID     int
	Command string
	CPU     float64
	Memory  float64
}

// PortInfo represents network port information
type PortInfo struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	State    string `json:"state"`
}

// NetworkDiagnostics represents network test results
type NetworkDiagnostics struct {
	PingResults string
	MTRResults  string
	DNSResults  string
	LastCheck   time.Time
}

var (
	systemStatsMutex sync.RWMutex
	networkMutex     sync.RWMutex
	currentStats     SystemStats
	networkResults   NetworkDiagnostics
)

func main() {
	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Handle routes
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/api/system/stats", handleSystemStats)
	http.HandleFunc("/api/network/diagnostics", handleNetworkDiagnostics)
	http.HandleFunc("/api/process/", handleProcessInfo)

	// Start background system stats collection
	go collectSystemStats()

	fmt.Println("Starting SystemHelper server on :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleSystemStats(w http.ResponseWriter, r *http.Request) {
	systemStatsMutex.RLock()
	defer systemStatsMutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentStats)
}

func handleNetworkDiagnostics(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		target := r.FormValue("target")
		if target == "" {
			http.Error(w, "Target is required", http.StatusBadRequest)
			return
		}

		results := runNetworkDiagnostics(target)
		networkMutex.Lock()
		networkResults = results
		networkMutex.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
		return
	}

	networkMutex.RLock()
	defer networkMutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(networkResults)
}

func handleProcessInfo(w http.ResponseWriter, r *http.Request) {
	pidStr := strings.TrimPrefix(r.URL.Path, "/api/process/")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		http.Error(w, "Invalid process ID", http.StatusBadRequest)
		return
	}

	info := getProcessInfo(pid)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

func collectSystemStats() {
	for {
		stats := getSystemStats()
		systemStatsMutex.Lock()
		currentStats = stats
		systemStatsMutex.Unlock()
		time.Sleep(5 * time.Second)
	}
}

func getSystemStats() SystemStats {
	var stats SystemStats

	// Get CPU usage
	if runtime.GOOS == "darwin" {
		// For macOS
		cmd := exec.Command("ps", "-A", "-o", "%cpu")
		output, _ := cmd.Output()
		stats.CPUUsage = parseCPUUsage(string(output))
	} else if runtime.GOOS == "linux" {
		cmd := exec.Command("top", "-bn1")
		output, _ := cmd.Output()
		stats.CPUUsage = parseCPUUsage(string(output))
	}

	// Get memory usage
	if runtime.GOOS == "darwin" {
		// For macOS
		cmd := exec.Command("vm_stat")
		output, _ := cmd.Output()
		stats.MemoryUsage = parseMemoryUsage(string(output))
	} else if runtime.GOOS == "linux" {
		cmd := exec.Command("free")
		output, _ := cmd.Output()
		stats.MemoryUsage = parseMemoryUsage(string(output))
	}

	// Get top processes
	if runtime.GOOS == "darwin" {
		cmd := exec.Command("ps", "-A", "-o", "pid,command,%cpu,%mem", "-r")
		output, _ := cmd.Output()
		stats.TopProcesses = parseTopProcessesMacOS(string(output))
	} else if runtime.GOOS == "linux" {
		cmd := exec.Command("ps", "-eo", "pid,comm,%cpu,%mem", "--sort=-%cpu")
		output, _ := cmd.Output()
		stats.TopProcesses = parseTopProcessesLinux(string(output))
	}

	// Get open ports
	if runtime.GOOS == "darwin" {
		cmd := exec.Command("netstat", "-anv", "-p", "tcp,udp")
		output, _ := cmd.Output()
		stats.OpenPorts = parseOpenPortsMacOS(string(output))
	} else if runtime.GOOS == "linux" {
		cmd := exec.Command("netstat", "-tuln")
		output, _ := cmd.Output()
		stats.OpenPorts = parseOpenPortsLinux(string(output))
	}

	return stats
}

func runNetworkDiagnostics(target string) NetworkDiagnostics {
	var results NetworkDiagnostics

	// Run ping
	pingCmd := exec.Command("ping", "-c", "4", target)
	pingOutput, _ := pingCmd.Output()
	results.PingResults = string(pingOutput)

	// Run MTR if available
	mtrCmd := exec.Command("mtr", "--report", target)
	mtrOutput, _ := mtrCmd.Output()
	results.MTRResults = string(mtrOutput)

	// Run DNS lookup
	dnsCmd := exec.Command("dig", target)
	dnsOutput, _ := dnsCmd.Output()
	results.DNSResults = string(dnsOutput)

	results.LastCheck = time.Now()
	return results
}

func getProcessInfo(pid int) map[string]interface{} {
	info := make(map[string]interface{})

	// Get process details using ps
	psCmd := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "pid,ppid,%cpu,%mem,command")
	psOutput, _ := psCmd.Output()
	info["ps"] = string(psOutput)

	// Get open files using lsof
	lsofCmd := exec.Command("lsof", "-p", strconv.Itoa(pid))
	lsofOutput, _ := lsofCmd.Output()
	info["lsof"] = string(lsofOutput)

	return info
}

func getTopProcesses() []ProcessInfo {
	var processes []ProcessInfo

	if runtime.GOOS == "linux" {
		cmd := exec.Command("ps", "-eo", "pid,comm,%cpu,%mem", "--sort=-%cpu")
		output, _ := cmd.Output()
		processes = parseTopProcessesLinux(string(output))
	} else {
		cmd := exec.Command("ps", "-A", "-o", "pid,command,%cpu,%mem", "-r")
		output, _ := cmd.Output()
		processes = parseTopProcessesMacOS(string(output))
	}

	return processes
}

func getOpenPorts() []PortInfo {
	var ports []PortInfo

	if runtime.GOOS == "linux" {
		cmd := exec.Command("netstat", "-tuln")
		output, _ := cmd.Output()
		ports = parseOpenPortsLinux(string(output))
	} else {
		cmd := exec.Command("netstat", "-an")
		output, _ := cmd.Output()
		ports = parseOpenPortsMacOS(string(output))
	}

	return ports
}

// Helper functions for parsing command outputs
func parseCPUUsage(output string) float64 {
	if runtime.GOOS == "darwin" {
		// For macOS, parse ps output
		lines := strings.Split(output, "\n")
		var totalCPU float64
		for _, line := range lines {
			if strings.TrimSpace(line) == "" {
				continue
			}
			// Skip header line
			if strings.Contains(line, "%CPU") {
				continue
			}
			// Parse CPU percentage
			fields := strings.Fields(line)
			if len(fields) > 0 {
				if cpu, err := strconv.ParseFloat(fields[0], 64); err == nil {
					totalCPU += cpu
				}
			}
		}
		return totalCPU
	}
	return 0.0 // Placeholder for other OS
}

func parseMemoryUsage(output string) float64 {
	if runtime.GOOS == "darwin" {
		// For macOS, parse vm_stat output
		lines := strings.Split(output, "\n")
		var pageSize uint64 = 4096 // Default page size on macOS is 4KB
		memoryStats := make(map[string]uint64)

		for _, line := range lines {
			fields := strings.Split(line, ":")
			if len(fields) != 2 {
				continue
			}

			key := strings.TrimSpace(fields[0])
			value := strings.TrimSpace(fields[1])
			if value == "" {
				continue
			}

			// Remove the dot at the end and "Pages" from beginning
			value = strings.TrimSuffix(value, ".")

			// Convert to uint64
			if num, err := strconv.ParseUint(value, 10, 64); err == nil {
				memoryStats[key] = num * pageSize
			}
		}

		// Calculate used memory
		usedMemory := memoryStats["Pages active"] +
			memoryStats["Pages inactive"] +
			memoryStats["Pages speculative"] +
			memoryStats["Pages wired down"]

		totalMemory := usedMemory + memoryStats["Pages free"]

		if totalMemory > 0 {
			return (float64(usedMemory) / float64(totalMemory)) * 100
		}
	}
	return 0.0
}

func parseTopProcessesLinux(output string) []ProcessInfo {
	// Implement Linux process parsing
	return []ProcessInfo{}
}

func parseTopProcessesMacOS(output string) []ProcessInfo {
	var processes []ProcessInfo
	lines := strings.Split(output, "\n")

	// Skip header line
	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		pid, err := strconv.Atoi(fields[0])
		if err != nil {
			continue
		}

		cpu, err := strconv.ParseFloat(fields[2], 64)
		if err != nil {
			continue
		}

		memory, err := strconv.ParseFloat(fields[3], 64)
		if err != nil {
			continue
		}

		// Combine remaining fields as command
		command := strings.Join(fields[1:], " ")

		processes = append(processes, ProcessInfo{
			PID:     pid,
			Command: command,
			CPU:     cpu,
			Memory:  memory,
		})
	}

	return processes
}

func parseOpenPortsLinux(output string) []PortInfo {
	// Implement Linux port parsing
	return []PortInfo{}
}

func parseOpenPortsMacOS(output string) []PortInfo {
	var ports []PortInfo
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if !strings.Contains(line, "LISTEN") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		proto := fields[0]
		localAddr := fields[3]

		// Extract port number
		var port int
		var err error

		if strings.Contains(localAddr, ".") {
			// IPv4
			parts := strings.Split(localAddr, ".")
			port, err = strconv.Atoi(parts[len(parts)-1])
		} else if strings.Contains(localAddr, ":") {
			// IPv6
			parts := strings.Split(localAddr, ":")
			port, err = strconv.Atoi(parts[len(parts)-1])
		}

		if err == nil && port > 0 {
			ports = append(ports, PortInfo{
				Port:     port,
				Protocol: proto,
				State:    "LISTEN",
			})
		}
	}

	return ports
}
