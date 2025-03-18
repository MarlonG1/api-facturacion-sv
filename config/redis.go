package config

import "github.com/MarlonG1/api-facturacion-sv/config/env"

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

func NewRedisConfig() *RedisConfig {
	return &RedisConfig{
		Host:     env.Redis.Host,
		Port:     env.Redis.Port,
		Password: env.Redis.Password,
	}
}

func (c *RedisConfig) GetURL() string {
	if c.Password != "" {
		return "redis://:" + c.Password + "@" + c.Host + ":" + c.Port
	}
	return "redis://" + c.Host + ":" + c.Port
}
