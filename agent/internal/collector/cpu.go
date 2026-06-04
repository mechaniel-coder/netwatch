package collector

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

// Snapshot holds one reporting interval's worth of system metrics.
type Snapshot struct {
	CPUPercent     float64 `json:"cpu_percent"`
	MemUsedBytes   uint64  `json:"mem_used_bytes"`
	MemTotalBytes  uint64  `json:"mem_total_bytes"`
	DiskUsedBytes  uint64  `json:"disk_used_bytes"`
	DiskTotalBytes uint64  `json:"disk_total_bytes"`
	NetRxBytes     uint64  `json:"net_rx_bytes"`
	NetTxBytes     uint64  `json:"net_tx_bytes"`
}

// Collect gathers all metrics in a single pass and returns a Snapshot.
func Collect() (*Snapshot, error) {
	cpuPct, err := collectCPU()
	if err != nil {
		return nil, fmt.Errorf("cpu: %w", err)
	}

	memUsed, memTotal, err := collectMemory()
	if err != nil {
		return nil, fmt.Errorf("memory: %w", err)
	}

	diskUsed, diskTotal, err := collectDisk("/")
	if err != nil {
		return nil, fmt.Errorf("disk: %w", err)
	}

	netRx, netTx, err := collectNetwork()
	if err != nil {
		return nil, fmt.Errorf("network: %w", err)
	}

	return &Snapshot{
		CPUPercent:     cpuPct,
		MemUsedBytes:   memUsed,
		MemTotalBytes:  memTotal,
		DiskUsedBytes:  diskUsed,
		DiskTotalBytes: diskTotal,
		NetRxBytes:     netRx,
		NetTxBytes:     netTx,
	}, nil
}

func collectCPU() (float64, error) {
	pcts, err := cpu.Percent(0, false)
	if err != nil || len(pcts) == 0 {
		return 0, err
	}
	return pcts[0], nil
}

func collectMemory() (used, total uint64, err error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return 0, 0, err
	}
	return v.Used, v.Total, nil
}

func collectDisk(path string) (used, total uint64, err error) {
	u, err := disk.Usage(path)
	if err != nil {
		return 0, 0, err
	}
	return u.Used, u.Total, nil
}

func collectNetwork() (rx, tx uint64, err error) {
	counters, err := net.IOCounters(false)
	if err != nil || len(counters) == 0 {
		return 0, 0, err
	}
	return counters[0].BytesRecv, counters[0].BytesSent, nil
}
