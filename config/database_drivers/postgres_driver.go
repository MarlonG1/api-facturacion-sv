package database_drivers

import (
	"fmt"
	"github.com/MarlonG1/api-facturacion-sv/config/env"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDriver struct{}

func NewPostgresDriver() *PostgresDriver {
	return &PostgresDriver{}
}

func (p *PostgresDriver) GetDSN() gorm.Dialector {
	return postgres.Open(p.GetStringConnection())
}

func (p *PostgresDriver) GetStringConnection() string {
	return fmt.Sprintf("host=%s users=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/El_Salvador  options='-c client_encoding=%s'",
		env.Database.User,
		env.Database.Password,
		env.Database.Host,
		env.Database.Port,
		env.Database.Name,
		env.Database.Charset,
	)
}

func (p *PostgresDriver) GetHost() string {
	return env.Database.Port
}

func (p *PostgresDriver) GetDriverName() string {
	return "PostgreSQL"
}
