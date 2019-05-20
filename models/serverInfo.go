package models

import (
	"fmt"
	"time"
)

var sslGrades = map[string]int{
	"A+": 7,
	"A":  6,
	"B":  5,
	"C":  4,
	"D":  3,
	"E":  2,
	"F":  1,
}

type HostInfo struct {
	Id               int          `json:"-" gorm:"primary_key"`
	Host             string       `json:"host"`
	ServersChanged   bool         `json:"servers_changed"`
	SslGrade         string       `json:"ssl_grade"`
	PreviousSslGrade string       `json:"previous_ssl_grade"`
	Logo             string       `json:"logo"`
	Title            string       `json:"title"`
	IsDown           bool         `json:"is_down"`
	Servers          []ServerInfo `json:"servers" gorm:"many2many:hostinfo_servers"`
	LastChecked      time.Time    `json:"-"`
}

type ServerInfo struct {
	Id         int    `json:"-" gorm:"primary_key"`
	Address    string `json:"address"`
	ServerName string `json:"server_name"`
	SslGrade   string `json:"ssl_grade"`
	Country    string `json:"country"`
	Owner      string `json:"owner"`
}

type History struct {
	Id         int       `json:"-" gorm:"primary_key"`
	Host       string    `json:"host"`
	Date       time.Time `json:"date"`
	HostInfo   HostInfo  `json:"info" gorm:"foreignkey:HostInfoId"`
	HostInfoId int       `json:"-"`
}

func (hostInfo HostInfo) PrintInfo() {
	fmt.Println(hostInfo.Host)
	fmt.Println(hostInfo.Logo)
	fmt.Println(hostInfo.Title)
	fmt.Println(hostInfo.IsDown)

	for _, server := range hostInfo.Servers {
		fmt.Println(server.Address)
		fmt.Println(server.ServerName)
		fmt.Println(server.Country)
		fmt.Println(server.Owner)
	}
}

func (hostInfo *HostInfo) SetStatus(status string) {
	if status == "READY" {
		hostInfo.IsDown = false
	} else {
		hostInfo.IsDown = true
	}
}

func (hostInfo *HostInfo) SetSslGrade() {
	numServers := len(hostInfo.Servers)
	tempGradeValue := 7
	tempGrade := "A+"

	for i := 0; i < numServers; i++ {
		if sslGrades[hostInfo.Servers[i].SslGrade] <= tempGradeValue {
			tempGradeValue = sslGrades[hostInfo.Servers[i].SslGrade]
			tempGrade = hostInfo.Servers[i].SslGrade
		}
	}

	hostInfo.SslGrade = tempGrade
}
