package database_drivers

import (
	"fmt"
	"github.com/MarlonG1/api-facturacion-sv/config/env"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlDriver struct{}

func NewMysqlDriver() *MysqlDriver {
	return &MysqlDriver{}
}

func (m *MysqlDriver) GetDSN() gorm.Dialector {
	return mysql.Open(m.GetStringConnection())
}

func (m *MysqlDriver) GetStringConnection() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=America%%2FEl_Salvador",
		env.Database.User,
		env.Database.Password,
		env.Database.Host,
		env.Database.Port,
		env.Database.Name,
		env.Database.Charset,
	)
}

func (m *MysqlDriver) GetHost() string {
	return env.Database.Port
}

func (m *MysqlDriver) GetDriverName() string {
	return "MySQL"
}
