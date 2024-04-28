package config

import (
	"os"
	"strconv"
	"sync"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

var (
	c        *Config
	syncOnce sync.Once
)

type Config struct {
	Service struct {
		Name        string
		Port        int
		Timeout     time.Duration
		Timezone    *time.Location
		Secret      string
		Environment string
		Debug       bool
	}
	Postgres struct {
		Host            string
		Port            string
		Username        string
		Password        string
		Database        string
		Schema          string
		PoolMaxConns    int
		PoolMinConns    int
		MaxConnLifetime int
		MaxConnIdleTime int
		SSLMode         string
	}
	GCP struct {
		ProjectID      string
		ServiceAccount []byte
	}
}

func load() *Config {
	cfg := new(Config)
	cfg.service()
	cfg.postgre()
	cfg.gcp()

	return cfg
}

func Get() *Config {
	syncOnce.Do(func() {
		c = load()
	})

	return c
}

func (cfg *Config) service() {
	cfg.Service.Name = os.Getenv("SERVICE_NAME")
	cfg.Service.Port, _ = strconv.Atoi(os.Getenv("SERVICE_PORT"))

	timeoutSec, _ := strconv.Atoi(os.Getenv("SERVICE_TIMEOUT"))
	cfg.Service.Timeout = time.Duration(timeoutSec) * time.Second

	location, _ := time.LoadLocation(os.Getenv("SERVICE_TIMEZONE"))
	cfg.Service.Timezone = location

	cfg.Service.Secret = os.Getenv("SERVICE_ENVIRONMENT")

	cfg.Service.Debug, _ = strconv.ParseBool(os.Getenv("SERVICE_DEBUG"))

	cfg.Service.Secret = os.Getenv("SERVICE_SECRET")
}

func (cfg *Config) postgre() {
	cfg.Postgres.Host = os.Getenv("POSTGRES_HOST")
	cfg.Postgres.Port = os.Getenv("POSTGRES_PORT")
	cfg.Postgres.Username = os.Getenv("POSTGRES_USERNAME")
	cfg.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.Postgres.Database = os.Getenv("POSTGRES_DB")
	cfg.Postgres.Schema = os.Getenv("POSTGRES_SCHEMA")
	cfg.Postgres.PoolMaxConns, _ = strconv.Atoi(os.Getenv("POSTGRES_POOL_MAX_CONNS"))
	cfg.Postgres.PoolMinConns, _ = strconv.Atoi(os.Getenv("POSTGRES_POOL_MIN_CONNS"))
	cfg.Postgres.MaxConnLifetime, _ = strconv.Atoi(os.Getenv("POSTGRES_MAX_CONN_LIFETIME"))
	cfg.Postgres.MaxConnIdleTime, _ = strconv.Atoi(os.Getenv("POSTGRES_MAX_CONN_IDLE_TIME"))
	cfg.Postgres.SSLMode = os.Getenv("POSTGRES_SSLMODE")
}

func (cfg *Config) gcp() {
	cfg.GCP.ServiceAccount = []byte(os.Getenv("GCP_SERVICE_ACCOUNT"))
	cfg.GCP.ProjectID = os.Getenv("GCP_PROJECT_ID")
}
