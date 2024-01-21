package main

import (
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/process"
)

type SSHSessionConfig struct {
	LoginUser string `json:"loginUser"`
	LoginIP   string `json:"loginIP"`
}
type sshSession struct {
	Username  string `json:"username"`
	PID       int32  `json:"PID"`
	Terminal  string `json:"terminal"`
	Host      string `json:"host"`
	LoginTime string `json:"loginTime"`
}

// https://github.com/1Panel-dev/1Panel/blob/732080b0bf1bc340e9d4bbf87734976545878879/backend/utils/websocket/process_data.go#L306C7-L306C7
func getSSHSessions(config SSHSessionConfig) (res []sshSession, err error) {
	var (
		result    []sshSession
		users     []host.UserStat
		processes []*process.Process
	)
	processes, err = process.Processes()
	if err != nil {
		return
	}
	users, err = host.Users()
	if err != nil {
		return
	}
	for _, proc := range processes {
		name, _ := proc.Name()
		if name != "sshd" || proc.Pid == 0 {
			continue
		}
		connections, _ := proc.Connections()
		for _, conn := range connections {
			for _, user := range users {
				if user.Host == "" {
					continue
				}
				if conn.Raddr.IP == user.Host {
					if config.LoginUser != "" && !strings.Contains(user.User, config.LoginUser) {
						continue
					}
					if config.LoginIP != "" && !strings.Contains(user.Host, config.LoginIP) {
						continue
					}
					if terminal, err := proc.Cmdline(); err == nil {
						if strings.Contains(terminal, user.Terminal) {
							session := sshSession{
								Username: user.User,
								Host:     user.Host,
								Terminal: user.Terminal,
								PID:      proc.Pid,
							}
							t := time.Unix(int64(user.Started), 0)
							session.LoginTime = t.Format("2006-1-2 15:04:05")
							result = append(result, session)
						}
					}
				}
			}
		}
	}
	res = result
	// res, err = json.Marshal(result)
	return
}
