package main

import (
	"strings"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

type NetConfig struct {
	Port        uint32 `json:"port"`
	ProcessName string `json:"processName"`
	ProcessID   int32  `json:"processID"`
}

var netTypes = [...]string{"tcp", "udp"}

// https://github.com/1Panel-dev/1Panel/blob/732080b0bf1bc340e9d4bbf87734976545878879/backend/utils/websocket/process_data.go#L361C14-L361C14
func getNetConnections(config NetConfig) (res []processConnect, err error) {
	var (
		result []processConnect
		proc   *process.Process
	)
	for _, netType := range netTypes {
		connections, _ := net.Connections(netType)
		if err == nil {
			for _, conn := range connections {
				if config.ProcessID > 0 && config.ProcessID != conn.Pid {
					continue
				}
				proc, err = process.NewProcess(conn.Pid)
				if err == nil {
					name, _ := proc.Name()
					if name != "" && config.ProcessName != "" && !strings.Contains(name, config.ProcessName) {
						continue
					}
					if config.Port > 0 && config.Port != conn.Laddr.Port && config.Port != conn.Raddr.Port {
						continue
					}
					result = append(result, processConnect{
						Type:   netType,
						Status: conn.Status,
						Laddr:  conn.Laddr,
						Raddr:  conn.Raddr,
						PID:    conn.Pid,
						Name:   name,
					})
				}

			}
		}
	}
	res = result
	// res, err = json.Marshal(result)
	return
}
