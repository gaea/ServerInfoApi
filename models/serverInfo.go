package models

type ServerInfo struct {
	Id         int    `json:"-"`
	Address    string `json:"address"`
	ServerName string `json:"server_name"`
	SslGrade   string `json:"ssl_grade"`
	Country    string `json:"country"`
	Owner      string `json:"owner"`
	HostInfoId int    `json:"-"`
}

func (serverInfo *ServerInfo) CheckUpdate(updatedServerInfo ServerInfo) bool {
	updated := false

	if serverInfo.SslGrade != updatedServerInfo.SslGrade ||
		serverInfo.ServerName != updatedServerInfo.ServerName ||
		serverInfo.Country != updatedServerInfo.Country ||
		serverInfo.Owner != updatedServerInfo.Owner {
		updated = true
	}

	return updated
}
