package config

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

func NewRedisConfig() *RedisConfig {
	return &RedisConfig{
		Host:     Redis.Host,
		Port:     Redis.Port,
		Password: Redis.Password,
	}
}

func (c *RedisConfig) GetURL() string {
	if c.Password != "" {
		return "redis://:" + c.Password + "@" + c.Host + ":" + c.Port
	}
	return "redis://" + c.Host + ":" + c.Port
}
