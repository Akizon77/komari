// report.go
package database

import (
	"fmt"

	"gorm.io/gorm"
)

// Report 表示客户端报告数据
type Report struct {
	UUID        string            `json:"uuid"`
	Token       string            `json:"token"`
	CPU         CPUReport         `json:"cpu"`
	OS          string            `json:"os"`
	Ram         RamReport         `json:"ram"`
	Swap        RamReport         `json:"swap"`
	Load        LoadReport        `json:"load"`
	Disk        DiskReport        `json:"disk"`
	Network     NetworkReport     `json:"network"`
	Connections ConnectionsReport `json:"connections"`
	Uptime      int64             `json:"uptime"`
	Process     int               `json:"process"`
	IPAddress   IPAddressReport   `json:"ip"`
	Message     string            `json:"message"`
}

type IPAddressReport struct {
	IPv4 string `json:"ipv4"`
	IPv6 string `json:"ipv6"`
}

type CPUReport struct {
	Name  string  `json:"name"`
	Cores int     `json:"cores"`
	Arch  string  `json:"arch"`
	Usage float64 `json:"usage"`
}

type GPUReport struct {
	Name  string  `json:"name"`
	Usage float64 `json:"usage"`
}

type RamReport struct {
	Total int64 `json:"total"`
	Used  int64 `json:"used"`
}

type LoadReport struct {
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}

type DiskReport struct {
	Total int64 `json:"total"`
	Used  int64 `json:"used"`
}

type NetworkReport struct {
	Up        int64 `json:"up"`
	Down      int64 `json:"down"`
	TotalUp   int64 `json:"totalUp"`
	TotalDown int64 `json:"totalDown"`
}

type ConnectionsReport struct {
	TCP int `json:"tcp"`
	UDP int `json:"udp"`
}

// SaveReport 保存报告数据
func SaveReport(data map[string]interface{}) (err error) {
	report := ParseReport(data)
	token := report.Token
	clientUUID, err := GetClientUUIDByToken(token)
	if err != nil {
		return err
	}
	return SaveClientReport(clientUUID, report)
}

// SaveClientReport 保存客户端报告到 History 表
func SaveClientReport(clientUUID string, report Report) (err error) {
	db := GetDBInstance()

	history := History{
		Client:         clientUUID,
		CPU:            report.CPU.Usage,
		GPU:            0, // Report 未提供 GPU Usage，设为 0（与原 nil 行为类似）
		RAM:            report.Ram.Used,
		RAMTotal:       report.Ram.Total,
		SWAP:           report.Swap.Used,
		SWAPTotal:      report.Swap.Total,
		LOAD:           report.Load.Load1, // 使用 Load1 作为主要负载指标
		TEMP:           0,                 // Report 未提供 TEMP，设为 0（与原 nil 行为类似）
		DISK:           report.Disk.Used,
		DISKTotal:      report.Disk.Total,
		NETIn:          report.Network.Down,
		NETOut:         report.Network.Up,
		NETTotalUp:     report.Network.TotalUp,
		NETTotalDown:   report.Network.TotalDown,
		PROCESS:        report.Process,
		Connections:    report.Connections.TCP,
		ConnectionsUDP: report.Connections.UDP,
	}

	// 保存到 ClientsInfo 表（CPU 和 OS 相关信息）
	clientInfo := ClientInfo{
		ClientUUID: clientUUID,
		CPUNAME:    report.CPU.Name,
		CPUARCH:    report.CPU.Arch,
		CPUCORES:   report.CPU.Cores,
		OS:         report.OS,
		GPUNAME:    "", // Report 未提供 GPU Name，留空
	}

	// 使用事务确保 History 和 ClientsInfo 一致性
	err = db.Transaction(func(tx *gorm.DB) error {
		// 保存 History
		if err := tx.Create(&history).Error; err != nil {
			return fmt.Errorf("failed to save history: %v", err)
		}

		// 更新或插入 ClientsInfo（使用 Upsert 逻辑）
		if err := tx.Save(&clientInfo).Error; err != nil {
			return fmt.Errorf("failed to save client info: %v", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// ParseReport 解析报告数据
func ParseReport(data map[string]interface{}) (report Report) {
	report = Report{
		Token: data["token"].(string),
		CPU: CPUReport{
			Name:  getString(data, "cpu.name"),
			Cores: getInt(data, "cpu.cores"),
			Arch:  getString(data, "cpu.arch"),
			Usage: getFloat(data, "cpu.usage"),
		},
		OS: getString(data, "os"),
		Ram: RamReport{
			Total: getInt64(data, "ram.total"),
			Used:  getInt64(data, "ram.used"),
		},
		Swap: RamReport{
			Total: getInt64(data, "swap.total"),
			Used:  getInt64(data, "swap.used"),
		},
		Load: LoadReport{
			Load1:  getFloat(data, "load.load1"),
			Load5:  getFloat(data, "load.load5"),
			Load15: getFloat(data, "load.load15"),
		},
		Disk: DiskReport{
			Total: getInt64(data, "disk.total"),
			Used:  getInt64(data, "disk.used"),
		},
		Network: NetworkReport{
			Up:        getInt64(data, "network.up"),
			Down:      getInt64(data, "network.down"),
			TotalUp:   getInt64(data, "network.totalUp"),
			TotalDown: getInt64(data, "network.totalDown"),
		},
		Connections: ConnectionsReport{
			TCP: getInt(data, "connections.tcp"),
			UDP: getInt(data, "connections.udp"),
		},
		Uptime:  getInt64(data, "uptime"),
		Process: getInt(data, "process"),
		Message: getString(data, "message"),
	}
	return
}

// getString 从 map 中获取字符串
func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// getInt 从 map 中获取整数
func getInt(data map[string]interface{}, key string) int {
	if val, ok := data[key]; ok {
		if num, ok := val.(float64); ok {
			return int(num)
		}
	}
	return 0
}

// getInt64 从 map 中获取 int64
func getInt64(data map[string]interface{}, key string) int64 {
	if val, ok := data[key]; ok {
		if num, ok := val.(float64); ok {
			return int64(num)
		}
	}
	return 0
}

// getFloat 从 map 中获取 float64
func getFloat(data map[string]interface{}, key string) float64 {
	if val, ok := data[key]; ok {
		if num, ok := val.(float64); ok {
			return num
		}
	}
	return 0.0
}
