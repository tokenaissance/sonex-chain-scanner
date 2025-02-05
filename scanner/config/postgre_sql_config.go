package config

import "fmt"

type PostgreSQLConfig struct {
	Host     string `json:","`
	Port     int64  `json:","`
	User     string `json:","`
	Password string `json:","`
	DbName   string `json:","`
}

func (c *PostgreSQLConfig) GetConnectString() string {
	return fmt.Sprintf(`postgres://%s:%s@%s:%d/%s?sslmode=disable`, c.User, c.Password, c.Host, c.Port, c.DbName)
}
