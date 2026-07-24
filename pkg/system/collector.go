package system

import (
	"bufio"
	"context"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// SystemMetrics 系统指标数据
type SystemMetrics struct {
	CPU     CPUMetrics     `json:"cpu"`
	Memory  MemoryMetrics  `json:"memory"`
	Disk    DiskMetrics    `json:"disk"`
	Network NetworkMetrics `json:"network"`
}

// CPUMetrics CPU 指标
type CPUMetrics struct {
	UsagePercent float64 `json:"usage_percent"`
	UserPercent  float64 `json:"user_percent"`
	SystemPercent float64 `json:"system_percent"`
	CoreCount    int     `json:"core_count"`
}

// MemoryMetrics 内存指标
type MemoryMetrics struct {
	Total        uint64  `json:"total"`
	Used         uint64  `json:"used"`
	Free         uint64  `json:"free"`
	Available    uint64  `json:"available"`
	UsagePercent float64 `json:"usage_percent"`
	Cached       uint64  `json:"cached"`
}

// DiskMetrics 磁盘指标
type DiskMetrics struct {
	Total        uint64  `json:"total"`
	Used         uint64  `json:"used"`
	Free         uint64  `json:"free"`
	UsagePercent float64 `json:"usage_percent"`
}

// NetworkMetrics 网络指标
type NetworkMetrics struct {
	BytesRecv uint64            `json:"bytes_recv"`
	BytesSent uint64            `json:"bytes_sent"`
	Interfaces map[string]NetworkInterface `json:"interfaces"`
}

// NetworkInterface 单个网络接口指标
type NetworkInterface struct {
	Name      string `json:"name"`
	BytesRecv uint64 `json:"bytes_recv"`
	BytesSent uint64 `json:"bytes_sent"`
}

// Collector 系统信息采集器
type Collector struct {
	metrics    SystemMetrics
	mu         sync.RWMutex
	prevCPU    []uint64
	prevTime   time.Time
}

// NewCollector 创建新的采集器实例
func NewCollector() *Collector {
	return &Collector{}
}

// Collect 采集一次系统指标
func (c *Collector) Collect() SystemMetrics {
	metrics := SystemMetrics{}
	metrics.CPU = c.collectCPU()
	metrics.Memory = c.collectMemory()
	metrics.Disk = c.collectDisk()
	metrics.Network = c.collectNetwork()

	c.mu.Lock()
	c.metrics = metrics
	c.mu.Unlock()

	return metrics
}

// collectCPU 采集 CPU 信息
func (c *Collector) collectCPU() CPUMetrics {
	metrics := CPUMetrics{}

	// 获取 CPU 核心数
	if data, err := os.ReadFile("/proc/cpuinfo"); err == nil {
		scanner := bufio.NewScanner(strings.NewReader(string(data)))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "processor") {
				metrics.CoreCount++
			}
		}
	}

	// 从 /proc/stat 读取 CPU 使用率
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return metrics
	}

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "cpu ") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 8 {
			continue
		}

		var values []uint64
		for i := 1; i < 8; i++ {
			v, _ := strconv.ParseUint(fields[i], 10, 64)
			values = append(values, v)
		}

		// 计算 CPU 使用率
		if c.prevCPU != nil && len(c.prevCPU) == 7 {
			prevIdle := c.prevCPU[3] + c.prevCPU[4]
			idle := values[3] + values[4]

			var prevTotal, total uint64
			for _, v := range c.prevCPU {
				prevTotal += v
			}
			for _, v := range values {
				total += v
			}

			if total != prevTotal {
				diffIdle := idle - prevIdle
				diffTotal := total - prevTotal
				usage := float64(diffTotal-diffIdle) / float64(diffTotal) * 100
				metrics.UsagePercent = usage

				userDiff := values[0] - c.prevCPU[0]
				sysDiff := values[2] - c.prevCPU[2]
				metrics.UserPercent = float64(userDiff) / float64(diffTotal) * 100
				metrics.SystemPercent = float64(sysDiff) / float64(diffTotal) * 100
			}
		}

		c.prevCPU = values
		c.prevTime = time.Now()
		break
	}

	return metrics
}

// collectMemory 采集内存信息
func (c *Collector) collectMemory() MemoryMetrics {
	metrics := MemoryMetrics{}

	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return metrics
	}

	var total, free, available, cached uint64
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		value, err := strconv.ParseUint(parts[1], 10, 64)
		if err != nil {
			continue
		}

		// /proc/meminfo 以 KB 为单位，转换为字节
		valueKB := value * 1024

		switch parts[0] {
		case "MemTotal:":
			total = valueKB
		case "MemFree:":
			free = valueKB
		case "MemAvailable:":
			available = valueKB
		case "Cached:":
			cached = valueKB
		}
	}

	metrics.Total = total
	metrics.Free = free
	metrics.Cached = cached

	if available > 0 {
		metrics.Available = available
	} else {
		metrics.Available = free + cached
	}

	metrics.Used = total - metrics.Available
	if total > 0 {
		metrics.UsagePercent = float64(metrics.Used) / float64(total) * 100
	}

	return metrics
}

// collectDisk 采集磁盘信息
func (c *Collector) collectDisk() DiskMetrics {
	metrics := DiskMetrics{}

	// 从 /proc/mounts 或直接使用 statfs
	fs, err := os.Statfs("/")
	if err != nil {
		return metrics
	}

	metrics.Total = fs.Blocks * uint64(fs.Bsize)
	metrics.Free = fs.Bfree * uint64(fs.Bsize)
	metrics.Used = metrics.Total - metrics.Free

	if metrics.Total > 0 {
		metrics.UsagePercent = float64(metrics.Used) / float64(metrics.Total) * 100
	}

	return metrics
}

// collectNetwork 采集网络信息
func (c *Collector) collectNetwork() NetworkMetrics {
	metrics := NetworkMetrics{
		Interfaces: make(map[string]NetworkInterface),
	}

	data, err := os.ReadFile("/proc/net/dev")
	if err != nil {
		return metrics
	}

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := scanner.Text()

		// 跳过头部行
		if strings.HasPrefix(line, "Inter-") || strings.HasPrefix(line, " face") || strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		name := strings.TrimSpace(parts[0])
		if name == "lo" {
			continue // 跳过回环接口
		}

		fields := strings.Fields(parts[1])
		if len(fields) < 10 {
			continue
		}

		bytesRecv, _ := strconv.ParseUint(fields[0], 10, 64)
		bytesSent, _ := strconv.ParseUint(fields[8], 10, 64)

		metrics.BytesRecv += bytesRecv
		metrics.BytesSent += bytesSent
		metrics.Interfaces[name] = NetworkInterface{
			Name:      name,
			BytesRecv: bytesRecv,
			BytesSent: bytesSent,
		}
	}

	return metrics
}

// GetMetrics 获取最新的系统指标（线程安全）
func (c *Collector) GetMetrics() SystemMetrics {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.metrics
}

// StartPeriodicCollect 启动周期性采集
func (c *Collector) StartPeriodicCollect(ctx context.Context, interval time.Duration) {
	// 先采集一次
	c.Collect()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.Collect()
		}
	}
}
