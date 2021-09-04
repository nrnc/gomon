package dbrepo

import (
	"context"
	"time"

	"github.com/nchukkaio/gomon/internal/models"
)

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

	return newID, nil
}

func (m *postgresDBRepo) GetHostByID(id int) (models.Host, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
				select host_name,canonical_name,url,ip,ipv6,location,os,active,created_at,updated_at
				from hosts
				where id = $1
	`
	var host models.Host
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
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
