//go:build windows
// +build windows

package main

import "github.com/shirou/gopsutil/v3/winservices"

type WinService struct {
	Services []winservices.Service `json:"winservices"`
}

func Winservice() any {
	// if runtime.GOOS == "windows" {
	winservicesInfo := WinService{}
	winservicesInfo.Services, _ = winservices.ListServices()
	return winservicesInfo
	// data = winservicesInfo
	// }
}
