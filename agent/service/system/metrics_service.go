package system

import (
	"fmt"
	"os"
	"os/user"
	"sort"
	"strconv"
	"strings"
	"syscall"
)

// CPUInfo holds CPU usage information.
type CPUInfo struct {
	ModelName    string  `json:"model_name"`
	Cores        int     `json:"cores"`
	Threads      int     `json:"threads"`
	Usage        float64 `json:"usage"`
	Frequency    string  `json:"frequency"`
	Temperature  float64 `json:"temperature"`
}

// MemoryInfo holds memory usage information.
type MemoryInfo struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Available   uint64  `json:"available"`
	Usage       float64 `json:"usage"`
	SwapTotal   uint64  `json:"swap_total"`
	SwapUsed    uint64  `json:"swap_used"`
	SwapFree    uint64  `json:"swap_free"`
}

// DiskInfo holds disk usage information.
type DiskInfo struct {
	Device      string  `json:"device"`
	MountPoint  string  `json:"mount_point"`
	FsType      string  `json:"fs_type"`
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Available   uint64  `json:"available"`
	Usage       float64 `json:"usage"`
	InodesTotal uint64  `json:"inodes_total"`
	InodesUsed  uint64  `json:"inodes_used"`
}

// NetworkInfo holds network interface information.
type NetworkInfo struct {
	Name      string  `json:"name"`
	RxBytes   uint64  `json:"rx_bytes"`
	TxBytes   uint64  `json:"tx_bytes"`
	RxSpeed   float64 `json:"rx_speed"`
	TxSpeed   float64 `json:"tx_speed"`
}

// ProcessInfo holds process information.
type ProcessInfo struct {
	PID     int     `json:"pid"`
	Name    string  `json:"name"`
	User    string  `json:"user"`
	CPU     float64 `json:"cpu"`
	Memory  float64 `json:"memory"`
	Status  string  `json:"status"`
	Command string  `json:"command"`
}

// SystemOverview holds a summary of system information.
type SystemOverview struct {
	Hostname   string     `json:"hostname"`
	OS         string     `json:"os"`
	Kernel     string     `json:"kernel"`
	Uptime     uint64     `json:"uptime"`
	Arch       string     `json:"arch"`
	CPU        *CPUInfo   `json:"cpu"`
	Memory     *MemoryInfo `json:"memory"`
	Disk       []DiskInfo  `json:"disk"`
	Network    []NetworkInfo `json:"network"`
}

// MetricsService provides system metrics collection operations.
type MetricsService struct{}

// NewMetricsService creates a new MetricsService.
func NewMetricsService() *MetricsService {
	return &MetricsService{}
}

// GetCPU returns CPU usage information.
func (s *MetricsService) GetCPU() (*CPUInfo, error) {
	info := &CPUInfo{}

	// Read CPU info
	data, err := os.ReadFile("/proc/cpuinfo")
	if err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "model name") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					info.ModelName = strings.TrimSpace(parts[1])
				}
				break
			}
		}
		for _, line := range lines {
			if strings.HasPrefix(line, "processor") {
				info.Threads++
			}
		}
	}

	// Read CPU frequency
	freqData, err := os.ReadFile("/sys/devices/system/cpu/cpu0/cpufreq/cpuinfo_max_freq")
	if err == nil {
		freq, _ := strconv.Atoi(strings.TrimSpace(string(freqData)))
		info.Frequency = fmt.Sprintf("%.2f GHz", float64(freq)/1e6)
	}

	// Read CPU usage from /proc/stat
	cpuUsage, err := s.readCPUUsage()
	if err == nil {
		info.Usage = cpuUsage
	}

	return info, nil
}

// GetMemory returns memory usage information.
func (s *MetricsService) GetMemory() (*MemoryInfo, error) {
	info := &MemoryInfo{}

	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return nil, fmt.Errorf("failed to read meminfo: %w", err)
	}

	memValues := parseProcFile(string(data))
	if total, ok := memValues["MemTotal"]; ok {
		info.Total = total
	}
	if available, ok := memValues["MemAvailable"]; ok {
		info.Available = available
		info.Used = info.Total - info.Available
	}
	if swapTotal, ok := memValues["SwapTotal"]; ok {
		info.SwapTotal = swapTotal
	}
	if swapFree, ok := memValues["SwapFree"]; ok {
		info.SwapFree = swapFree
		info.SwapUsed = info.SwapTotal - info.SwapFree
	}

	if info.Total > 0 {
		info.Usage = float64(info.Used) / float64(info.Total) * 100.0
	}

	return info, nil
}

// GetDisk returns disk usage information for all mounted filesystems.
func (s *MetricsService) GetDisk() ([]DiskInfo, error) {
	var disks []DiskInfo

	mountsData, err := os.ReadFile("/proc/mounts")
	if err != nil {
		return nil, fmt.Errorf("failed to read mounts: %w", err)
	}

	lines := strings.Split(string(mountsData), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		device := fields[0]
		mountPoint := fields[1]
		fsType := fields[2]

		// Skip virtual filesystems
		if strings.HasPrefix(device, "none") ||
			strings.HasPrefix(device, "tmpfs") ||
			strings.HasPrefix(device, "cgroup") ||
			strings.HasPrefix(device, "proc") ||
			strings.HasPrefix(device, "sys") ||
			strings.HasPrefix(device, "dev") ||
			mountPoint == "/proc" ||
			mountPoint == "/sys" ||
			mountPoint == "/dev" ||
			mountPoint == "/run" {
			continue
		}

		fs := syscall.Statfs_t{}
		if err := syscall.Statfs(mountPoint, &fs); err != nil {
			continue
		}

		total := fs.Blocks * uint64(fs.Bsize)
		available := fs.Bavail * uint64(fs.Bsize)
		used := total - available

		var usage float64
		if total > 0 {
			usage = float64(used) / float64(total) * 100.0
		}

		disks = append(disks, DiskInfo{
			Device:      device,
			MountPoint:  mountPoint,
			FsType:      fsType,
			Total:       total,
			Used:        used,
			Available:   available,
			Usage:       usage,
			InodesTotal: fs.Files,
			InodesUsed:  fs.Files - fs.Ffree,
		})
	}

	return disks, nil
}

// GetNetwork returns network interface statistics.
func (s *MetricsService) GetNetwork() ([]NetworkInfo, error) {
	var networks []NetworkInfo

	entries, err := os.ReadDir("/sys/class/net")
	if err != nil {
		return nil, fmt.Errorf("failed to read network interfaces: %w", err)
	}

	for _, entry := range entries {
		name := entry.Name()
		if name == "lo" {
			continue
		}

		info := NetworkInfo{Name: name}

		// Read RX/TX bytes
		rxBytes, err := readSysFile(fmt.Sprintf("/sys/class/net/%s/statistics/rx_bytes", name))
		if err == nil {
			info.RxBytes = rxBytes
		}
		txBytes, err := readSysFile(fmt.Sprintf("/sys/class/net/%s/statistics/tx_bytes", name))
		if err == nil {
			info.TxBytes = txBytes
		}

		networks = append(networks, info)
	}

	return networks, nil
}

// GetOverview returns a comprehensive system overview.
func (s *MetricsService) GetOverview() (*SystemOverview, error) {
	overview := &SystemOverview{}

	hostname, err := os.Hostname()
	if err == nil {
		overview.Hostname = hostname
	}

	uname := syscall.Utsname{}
	if err := syscall.Uname(&uname); err == nil {
		overview.Kernel = charsToString(uname.Release[:])
		overview.Arch = charsToString(uname.Machine[:])
	}

	osInfo := getOSInfo()
	overview.OS = osInfo

	// Uptime
	uptimeData, err := os.ReadFile("/proc/uptime")
	if err == nil {
		parts := strings.Fields(string(uptimeData))
		if len(parts) > 0 {
			uptimeSec, _ := strconv.ParseFloat(parts[0], 64)
			overview.Uptime = uint64(uptimeSec)
		}
	}

	cpu, _ := s.GetCPU()
	overview.CPU = cpu

	mem, _ := s.GetMemory()
	overview.Memory = mem

	disk, _ := s.GetDisk()
	overview.Disk = disk

	network, _ := s.GetNetwork()
	overview.Network = network

	return overview, nil
}

// GetProcesses returns a paginated list of processes.
func (s *MetricsService) GetProcesses(page, size int) ([]ProcessInfo, int64, error) {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read /proc: %w", err)
	}

	var processes []ProcessInfo
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}

		proc, err := s.getProcessInfo(pid)
		if err != nil {
			continue
		}
		processes = append(processes, *proc)
	}

	// Sort by CPU usage descending
	sort.Slice(processes, func(i, j int) bool {
		return processes[i].CPU > processes[j].CPU
	})

	total := int64(len(processes))
	start := (page - 1) * size
	if start >= len(processes) {
		return nil, total, nil
	}
	end := start + size
	if end > len(processes) {
		end = len(processes)
	}

	return processes[start:end], total, nil
}

// getProcessInfo reads process information from /proc.
func (s *MetricsService) getProcessInfo(pid int) (*ProcessInfo, error) {
	info := &ProcessInfo{PID: pid}

	// Read /proc/[pid]/comm for name
	commData, err := os.ReadFile(fmt.Sprintf("/proc/%d/comm", pid))
	if err == nil {
		info.Name = strings.TrimSpace(string(commData))
	}

	// Read /proc/[pid]/cmdline for command
	cmdlineData, err := os.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err == nil {
		info.Command = strings.ReplaceAll(string(cmdlineData), "\x00", " ")
	}

	// Read /proc/[pid]/stat for user, cpu, memory
	statData, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err == nil {
		fields := strings.Fields(string(statData))
		if len(fields) > 3 {
			uid, _ := strconv.Atoi(fields[1])
			info.User = lookupUsername(uid)
		}
		if len(fields) > 17 {
			utime, _ := strconv.ParseFloat(fields[13], 64)
			stime, _ := strconv.ParseFloat(fields[14], 64)
			totalTime := utime + stime
			info.CPU = totalTime * 100
		}
	}

	// Read /proc/[pid]/status for memory
	statusData, err := os.ReadFile(fmt.Sprintf("/proc/%d/status", pid))
	if err == nil {
		lines := strings.Split(string(statusData), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "VmRSS:") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					rss, _ := strconv.ParseFloat(fields[1], 64)
					info.Memory = rss * 100
				}
			}
			if strings.HasPrefix(line, "State:") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					info.Status = fields[1]
				}
			}
		}
	}

	return info, nil
}

// readCPUUsage calculates CPU usage from /proc/stat.
func (s *MetricsService) readCPUUsage() (float64, error) {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return 0, err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if !strings.HasPrefix(line, "cpu ") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 5 {
			return 0, nil
		}

		var total, idle float64
		for i := 1; i < len(fields); i++ {
			val, _ := strconv.ParseFloat(fields[i], 64)
			total += val
			if i == 4 {
				idle = val
			}
		}

		if total > 0 {
			return (1 - idle/total) * 100, nil
		}
	}
	return 0, nil
}

// parseProcFile parses a key-value file like /proc/meminfo.
func parseProcFile(data string) map[string]uint64 {
	result := make(map[string]uint64)
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(strings.TrimSuffix(parts[1], " kB"))
		val, _ := strconv.ParseUint(value, 10, 64)
		result[key] = val * 1024 // Convert kB to bytes
	}
	return result
}

// readSysFile reads a numeric value from a sysfs file.
func readSysFile(path string) (uint64, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64)
}

// lookupUsername returns a username for a UID.
func lookupUsername(uid int) string {
	u, err := user.LookupId(strconv.Itoa(uid))
	if err != nil {
		return strconv.Itoa(uid)
	}
	return u.Username
}

// charsToString converts a C char array to a Go string.
func charsToString(ca []int8) string {
	s := make([]byte, len(ca))
	var l int
	for ; l < len(ca); l++ {
		if ca[l] == 0 {
			break
		}
		s[l] = byte(ca[l])
	}
	return string(s[:l])
}

// getOSInfo returns a string describing the operating system.
func getOSInfo() string {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return "Linux"
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			value := strings.TrimPrefix(line, "PRETTY_NAME=")
			return strings.Trim(value, "\"")
		}
	}
	return "Linux"
}
