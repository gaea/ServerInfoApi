package repositories

import (
	"fmt"
	"log"

	"../configs"
	"../models"
	_ "github.com/lib/pq"
)

func ServerInfoRepository() *Repository {
	return &Repository{db: configs.DatabaseSetup()}
}

func (rep *Repository) CreateServerInfo(serverInfo *models.ServerInfo) error {
	var serverInfoId int

	qryString := fmt.Sprintf(
		`INSERT INTO server_infos 
			(address, server_name, ssl_grade, country, owner, host_info_id) 
		VALUES 
			('%s', '%s', '%s', '%s', '%s', %d) 
		RETURNING 
			id`,
		serverInfo.Address, serverInfo.ServerName, serverInfo.SslGrade, serverInfo.Country, serverInfo.Owner, serverInfo.HostInfoId)

	err := rep.db.QueryRow(qryString).Scan(&serverInfoId)

	if err != nil {
		return err
	}

	serverInfo.Id = serverInfoId

	return nil
}

func (rep *Repository) UpdateServerInfo(serverInfo *models.ServerInfo) error {
	qryString := fmt.Sprintf(`
								UPDATE 
									server_infos
								SET
									server_name = '%s', ssl_grade = '%s', country = '%s', owner = '%s'
								WHERE 
									id = %d`, serverInfo.ServerName, serverInfo.SslGrade, serverInfo.Country, serverInfo.Owner, serverInfo.Id)

	_, err := rep.db.Exec(qryString)

	return err
}

func (rep *Repository) ListServerInfos(hostInfoId int) ([]models.ServerInfo, error) {
	serverInfos := []models.ServerInfo{}

	qryString := fmt.Sprintf(`SELECT 
									id, address, server_name, ssl_grade, country, owner, host_info_id 
								FROM 
									server_infos 
								WHERE 
									host_info_id = %d`, hostInfoId)

	log.Println(qryString)

	rows, err := rep.db.Query(qryString)

	if err != nil {
		log.Println(err)
		return serverInfos, err
	}

	defer rows.Close()

	for rows.Next() {
		var serverInfo models.ServerInfo
		if err := rows.Scan(&serverInfo.Id, &serverInfo.Address, &serverInfo.ServerName,
			&serverInfo.SslGrade, &serverInfo.Country, &serverInfo.Owner, &serverInfo.HostInfoId); err != nil {
			log.Fatal(err)
		}
		serverInfos = append(serverInfos, serverInfo)
	}

	return serverInfos, nil
}

func (rep *Repository) DeleteServerInfo(id int) error {
	qryString := fmt.Sprintf("DELETE from server_infos WHERE id = %d", id)

	_, err := rep.db.Exec(qryString)
	return err
}
