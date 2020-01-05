package conf

import (
	"os"
	"strconv"
	"time"
)

// Config vernemq-auth configurations
type Config struct {
	Port         int
	LogLevel     string
	Redis        RedisConfig
	CacheControl CacheControl
	AuthConfig   AuthConfig
}

// AuthConfig auth service configurations
type AuthConfig struct {
	LruSize int
}

// RedisConfig redis configurations
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// CacheControl vernemq cache control header configurations
type CacheControl struct {
	Register  time.Duration
	Publish   time.Duration
	Subscribe time.Duration
}

// ParseConfig returns config
func ParseConfig() Config {
	return Config{
		Port:     intWithDefault("PORT", 9595),
		LogLevel: stringWithDefault("LOG_LEVEL", "error"),
		AuthConfig: AuthConfig{
			LruSize: intWithDefault(os.Getenv("AUTH_LRU_SIZE"), 2000),
		},
		Redis: RedisConfig{
			Addr:     stringWithDefault(os.Getenv("REDIS_ADDR"), "localhost:6379"),
			Password: stringWithDefault(os.Getenv("REDIS_PASSWORD"), ""),
			DB:       intWithDefault(os.Getenv("REDIS_DB"), 0),
		},
		CacheControl: CacheControl{
			Register:  durationWithDefault(os.Getenv("CC_REGISTER"), time.Minute*10),
			Publish:   durationWithDefault(os.Getenv("CC_PUBLISH"), time.Minute*10),
			Subscribe: durationWithDefault(os.Getenv("CC_SUBSCRIBE"), time.Minute*10),
		},
	}
}

func intWithDefault(value string, def int) int {
	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return def
	}

	return int(v)
}

func stringWithDefault(value string, def string) string {
	if len(value) == 0 {
		return def
	}

	return value
}

func durationWithDefault(value string, def time.Duration) time.Duration {
	duration, err := time.ParseDuration(value)
	if err != nil {
		return def
	}

	return duration
}
