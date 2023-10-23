package configs

import "fmt"

type DatabaseConfig struct {
	Host string
	Port int
	User string
	Pass string
	Name string
}

func (c *DatabaseConfig) MySqlConnectStr() string {
	format := "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	return fmt.Sprintf(format, c.User, c.Pass, c.Host, c.Port, c.Name)
}

func (c *DatabaseConfig) PgSqlConnectStr() string {
	format := "host=%s user=%s password=%s dbname=%s port=%d"
	return fmt.Sprintf(format, c.Host, c.User, c.Pass, c.Name, c.Port)
}

func NewDatabaseConfig(env Env) *DatabaseConfig {
	return &DatabaseConfig{
		Host: env["DB_HOST"].(string),
		Port: env["DB_PORT"].(int),
		User: env["DB_USER"].(string),
		Pass: env["DB_PASS"].(string),
		Name: env["DB_NAME"].(string),
	}
}
