package dbrepo

import (
	"context"
	"log"
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
		var hs []models.HostService
		serviceQuery := `
					select 
						hs.id,hs.host_id,hs.service_id,hs.active,hs.schedule_number,hs.schedule_unit,hs.last_check,hs.created_at,hs.updated_at,hs.status,
						s.id,s.service_name,s.active,s.icon,s.created_at,s.updated_at
					from
						host_services hs
					left join services s on (s.id=hs.service_id)
					where
						host_id = $1
		`
		serviceRows, err := m.DB.QueryContext(ctx, serviceQuery, host.ID)

		if err != nil {
			log.Println(err)
			return hosts, err
		}
		for serviceRows.Next() {
			var h models.HostService
			err := serviceRows.Scan(
				&h.ID,
				&h.HostID,
				&h.ServiceID,
				&h.Active,
				&h.ScheduleNumber,
				&h.ScheduleUnit,
				&h.LastCheck,
				&h.CreatedAt,
				&h.UpdatedAt,
				&h.Status,
				&h.Service.ID,
				&h.Service.ServiceName,
				&h.Service.Active,
				&h.Service.Icon,
				&h.Service.CreatedAt,
				&h.Service.UpdatedAt,
			)

			if err != nil {
				return hosts, err
			}
			hs = append(hs, h)
		}
		host.HostServices = hs
		serviceRows.Close()
		hosts = append(hosts, host)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return hosts, nil
}

func (m *postgresDBRepo) GetAllServicesStatusCount() (int, int, int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
				select
					(select count(id) from host_services where active = 1 and status = 'pending') as pending,
					(select count(id) from host_services where active = 1 and status = 'healthy') as healthy,
					(select count(id) from host_services where active = 1 and status = 'warning') as warning,
					(select count(id) from host_services where active = 1 and status = 'problem') as problem
	`
	var pending, healthy, warning, problem int
	row := m.DB.QueryRowContext(ctx, stmt)
	err := row.Scan(&pending, &healthy, &warning, &problem)

	if err != nil {
		return 0, 0, 0, 0, err
	}

	return pending, healthy, warning, problem, nil

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

	var hs []models.HostService
	query = `
				select 
					hs.id,hs.host_id,hs.service_id,hs.active,hs.schedule_number,hs.schedule_unit,hs.last_check,hs.created_at,hs.updated_at,hs.status,
					s.id,s.service_name,s.active,s.icon,s.created_at,s.updated_at
				from
					host_services hs
				left join services s on (s.id=hs.service_id)
				where
					host_id = $1
	`
	rows, err := m.DB.QueryContext(ctx, query, host.ID)

	if err != nil {
		log.Println(err)
		return host, err
	}
	for rows.Next() {
		var h models.HostService
		err := rows.Scan(
			&h.ID,
			&h.HostID,
			&h.ServiceID,
			&h.Active,
			&h.ScheduleNumber,
			&h.ScheduleUnit,
			&h.LastCheck,
			&h.CreatedAt,
			&h.UpdatedAt,
			&h.Status,
			&h.Service.ID,
			&h.Service.ServiceName,
			&h.Service.Active,
			&h.Service.Icon,
			&h.Service.CreatedAt,
			&h.Service.UpdatedAt,
		)

		if err != nil {
			return host, err
		}
		hs = append(hs, h)

	}
	host.HostServices = hs
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

func (m *postgresDBRepo) UpdateHostServiceStatus(serviceID, hostID, active int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
				update host_services set active = $1 where service_id = $2 and host_id = $3
	`

	_, err := m.DB.ExecContext(ctx, stmt, active, serviceID, hostID)

	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) GetServicesByStatus(status string) ([]models.HostService, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
				select
				hs.id,hs.host_id,hs.service_id,hs.active,hs.schedule_number,hs.schedule_unit,hs.last_check,hs.created_at,hs.updated_at,hs.status,
				h.host_name,s.service_name
				from 
					host_services hs
					left join hosts h on (h.id = hs.host_id)
					left join services s on (hs.service_id = s.id)
				where
					status = $1

	`
	var hostServices []models.HostService
	rows, err := m.DB.QueryContext(ctx, query, status)
	if err != nil {
		return hostServices, err
	}
	defer rows.Close()
	for rows.Next() {
		var hs models.HostService
		err := rows.Scan(
			&hs.ID,
			&hs.HostID,
			&hs.ServiceID,
			&hs.Active,
			&hs.ScheduleNumber,
			&hs.ScheduleUnit,
			&hs.LastCheck,
			&hs.CreatedAt,
			&hs.UpdatedAt,
			&hs.Status,
			&hs.HostName,
			&hs.Service.ServiceName,
		)
		if err != nil {
			return nil, err
		}
		hostServices = append(hostServices, hs)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return hostServices, nil
}
