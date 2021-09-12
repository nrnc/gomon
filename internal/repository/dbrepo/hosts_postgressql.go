package dbrepo

import (
	"context"
	"time"

	"github.com/nchukkaio/gomon/internal/models"
)

func (m *postgresDBRepo) AllHosts() ([]models.Host, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
				select id,host_name,canonical_name,url,ip,ipv6,location,os,active,created_at,updated_at
				from hosts
	`
	var hosts []models.Host
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return hosts, err
	}
	defer rows.Close()

	for rows.Next() {
		var host models.Host
		rows.Scan(
			&host.ID,
			&host.HostName,
			&host.CanonicalName,
			&host.URL,
			&host.IP,
			&host.IPV6,
			&host.Location,
			&host.OS,
			&host.Active,
			&host.CreatedAt,
			&host.UpdatedAt,
		)
		hosts = append(hosts, host)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return hosts, nil
}

func (m *postgresDBRepo) InsertHost(host models.Host) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
				insert into hosts (host_name,canonical_name,url,ip,ipv6,location,os,active,created_at,updated_at)
				values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
				returning id
	`
	var newID int
	err := m.DB.QueryRowContext(ctx, query,
		host.HostName,
		host.CanonicalName,
		host.URL,
		host.IP,
		host.IPV6,
		host.Location,
		host.OS,
		host.Active,
		host.CreatedAt,
		host.UpdatedAt,
	).Scan(&newID)
	if err != nil {
		return newID, err
	}
	// add host service and set to inactive
	stmt := `
				insert into host_services(host_id,service_id,active,schedule_number,schedule_unit,
				status,created_at,updated_at) values($1,1,0,3,'m','pending',$2,$3)
	`
	_, err = m.DB.ExecContext(ctx, stmt, newID, time.Now(), time.Now())
	if err != nil {
		return newID, err
	}
	return newID, nil
}

func (m *postgresDBRepo) GetHostByID(id int) (models.Host, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
				select id,host_name,canonical_name,url,ip,ipv6,location,os,active,created_at,updated_at
				from hosts
				where id = $1
	`
	var host models.Host
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&host.ID,
		&host.HostName,
		&host.CanonicalName,
		&host.URL,
		&host.IP,
		&host.IPV6,
		&host.Location,
		&host.OS,
		&host.Active,
		&host.CreatedAt,
		&host.UpdatedAt,
	)

	if err != nil {
		return host, err
	}

	return host, nil
}

func (m *postgresDBRepo) UpdateHost(host models.Host) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
			update hosts
			set
			host_name = $2,
			canonical_name = $3,
			url = $4,
			ip = $5,
			ipv6 = $6,
			location = $7,
			os = $8,
			active = $9,
			created_at = $10,
			updated_at = $11
			where id = $1
	`

	_, err := m.DB.QueryContext(ctx, query,
		host.ID,
		host.HostName,
		host.CanonicalName,
		host.URL,
		host.IP,
		host.IPV6,
		host.Location,
		host.OS,
		host.Active,
		host.CreatedAt,
		host.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}
