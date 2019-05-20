package models

type SslLabsResponse struct {
	Host     string          `json:"host"`
	Port     int             `json:"port"`
	Protocol string          `json:"protocol"`
	Status   string          `json:"status"`
	Servers  []SslLabsServer `json:"endpoints"`
}

type SslLabsServer struct {
	Address    string `json:"ipAddress"`
	ServerName string `json:"serverName"`
	SslGrade   string `json:"grade"`
}
