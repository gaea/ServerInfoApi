CREATE DATABASE servers;

CREATE USER IF NOT EXISTS gaea;

GRANT ALL ON DATABASE servers TO gaea;

USE servers;

CREATE TABLE IF NOT EXISTS host_infos (
	id INT8 NOT NULL DEFAULT unique_rowid(),
	host STRING NULL,
	servers_changed BOOL NULL,
	ssl_grade STRING NULL,
	previous_ssl_grade STRING NULL,
	logo STRING NULL,
	title STRING NULL,
	is_down BOOL NULL,
	last_checked TIMESTAMPTZ NULL,
	CONSTRAINT "primary" PRIMARY KEY (id ASC),
	FAMILY "primary" (id, host, servers_changed, ssl_grade, previous_ssl_grade, logo, title, is_down, last_checked)
);

CREATE TABLE IF NOT EXISTS server_infos (
	id INT8 NOT NULL DEFAULT unique_rowid(),
	address STRING NULL,
	server_name STRING NULL,
	ssl_grade STRING NULL,
	country STRING NULL,
	owner STRING NULL,
	host_info_id INT8 NOT NULL,
	CONSTRAINT "primary" PRIMARY KEY (id ASC),
	FAMILY "primary" (id, address, server_name, ssl_grade, country, owner, host_info_id)
);

CREATE TABLE IF NOT EXISTS histories (
	id INT8 NOT NULL DEFAULT unique_rowid(),
	host STRING NULL,
	date TIMESTAMPTZ NULL,
	host_info_id INT8 NULL,
	CONSTRAINT "primary" PRIMARY KEY (id ASC),
	FAMILY "primary" (id, host, date, host_info_id)
);
