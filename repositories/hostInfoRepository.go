package repositories

import (
	"fmt"
	"log"

	"../configs"
	"../models"
	_ "github.com/lib/pq"
)

func HostInfoRepository() *Repository {
	return &Repository{db: configs.DatabaseSetup()}
}

func (rep *Repository) CreateHostInfo(hostInfo *models.HostInfo) error {
	var hostInfoId int

	qryString := fmt.Sprintf(
		`INSERT INTO host_infos 
			(host, servers_changed, ssl_grade, previous_ssl_grade, logo, title, is_down, last_checked) 
		VALUES 
			('%s', %t, '%s', '%s', '%s', '%s', %t, '%s') 
		RETURNING 
			id`,
		hostInfo.Host, hostInfo.ServersChanged, hostInfo.SslGrade, hostInfo.SslGrade, hostInfo.Logo, hostInfo.Title, hostInfo.IsDown, "NOW()")

	err := rep.db.QueryRow(qryString).Scan(&hostInfoId)

	if err != nil {
		return err
	} else {
		hostInfo.Id = hostInfoId
		hostInfo.PreviousSslGrade = hostInfo.SslGrade

		for i := 0; i < len(hostInfo.Servers); i++ {
			hostInfo.Servers[i].HostInfoId = hostInfo.Id
			err = rep.CreateServerInfo(&hostInfo.Servers[i])

			if err != nil {
				return err
			}
		}

		return nil
	}
}

func (rep *Repository) CheckHostInfo(dbHostInfo models.HostInfo, updatedHostInfo *models.HostInfo) error {
	updatedServers := false
	updatedHostInfo.Id = dbHostInfo.Id

	lenUpdatedServers := len(updatedHostInfo.Servers)
	lenDBServers := len(dbHostInfo.Servers)
	dbServersFound := make([]string, 0)

	// check if a new server was added or removed of the related domain
	updatedServers = lenUpdatedServers != lenDBServers

	for i := 0; i < lenUpdatedServers; i++ {
		updatedServerFound := false

		for j := 0; j < lenDBServers; j++ {
			if updatedHostInfo.Servers[i].Address == dbHostInfo.Servers[j].Address {
				dbServersFound = append(dbServersFound, dbHostInfo.Servers[j].Address)
				updatedServerFound = true
				serverUpdated := dbHostInfo.Servers[j].CheckUpdate(updatedHostInfo.Servers[i])

				if serverUpdated {
					updatedServers = true

					updatedHostInfo.Servers[i].Id = dbHostInfo.Servers[j].Id
					rep.UpdateServerInfo(&updatedHostInfo.Servers[i])
					break
				}
			}
		}

		if !updatedServerFound {
			rep.CreateServerInfo(&updatedHostInfo.Servers[i])
		}
	}

	for i := 0; i < lenDBServers; i++ {
		deleteServer := true
		for j := 0; j < len(dbServersFound); j++ {

			if dbHostInfo.Servers[i].Address == dbServersFound[j] {
				deleteServer = false

				break
			}
		}

		if deleteServer {
			rep.DeleteServerInfo(dbHostInfo.Servers[i].Id)
		}
	}

	if updatedServers {
		updatedHostInfo.ServersChanged = updatedServers
		updatedHostInfo.SetSslGrade()
		updatedHostInfo.PreviousSslGrade = dbHostInfo.SslGrade
		rep.UpdateHostInfo(updatedHostInfo)
	}

	return nil
}

func (rep *Repository) UpdateHostInfo(hostInfo *models.HostInfo) error {
	qryString := fmt.Sprintf(`
								UPDATE 
									host_infos
								SET
									servers_changed = %t, ssl_grade = '%s', previous_ssl_grade = '%s', logo = '%s', title = '%s', is_down = %t, last_checked = NOW()
								WHERE 
									id = %d`, hostInfo.ServersChanged, hostInfo.SslGrade, hostInfo.PreviousSslGrade, hostInfo.Logo, hostInfo.Title, hostInfo.IsDown, hostInfo.Id)

	_, err := rep.db.Exec(qryString)

	return err
}

func (rep *Repository) DetailHostInfo(host string) (models.HostInfo, error) {
	var hostInfo models.HostInfo

	qryString := fmt.Sprintf(`SELECT 
									id, host, servers_changed, ssl_grade, previous_ssl_grade, logo, title, is_down, last_checked 
								FROM 
									host_infos 
								WHERE 
									host = '%s'`, host)

	log.Println(qryString)

	rows, err := rep.db.Query(qryString)

	if err != nil {
		return hostInfo, err
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&hostInfo.Id, &hostInfo.Host, &hostInfo.ServersChanged,
			&hostInfo.SslGrade, &hostInfo.PreviousSslGrade, &hostInfo.Logo, &hostInfo.Title,
			&hostInfo.IsDown, &hostInfo.LastChecked); err != nil {
			return hostInfo, err
		} else {
			var serverInfos, err = rep.ListServerInfos(hostInfo.Id)

			if err != nil {
				log.Println(err)
				return hostInfo, err
			}

			hostInfo.Servers = serverInfos
		}
	}

	return hostInfo, nil
}

func (rep *Repository) ListHostInfo() ([]models.HostInfo, error) {
	hostInfos := []models.HostInfo{}

	qryString := `SELECT 
									id, host, servers_changed, ssl_grade, previous_ssl_grade, logo, title, is_down, last_checked  
								FROM 
									host_infos 
								ORDER BY 
									last_checked DESC`

	log.Println(qryString)

	rows, err := rep.db.Query(qryString)

	if err != nil {
		log.Println(err)
		return hostInfos, err
	}

	defer rows.Close()

	for rows.Next() {
		var hostInfo models.HostInfo
		if err := rows.Scan(&hostInfo.Id, &hostInfo.Host, &hostInfo.ServersChanged,
			&hostInfo.SslGrade, &hostInfo.PreviousSslGrade, &hostInfo.Logo, &hostInfo.Title,
			&hostInfo.IsDown, &hostInfo.LastChecked); err != nil {
			return hostInfos, err
		} else {
			var serverInfos, err = rep.ListServerInfos(hostInfo.Id)

			if err != nil {
				log.Println(err)
				return hostInfos, err
			}

			hostInfo.Servers = serverInfos
		}

		hostInfos = append(hostInfos, hostInfo)
	}

	return hostInfos, nil
}
