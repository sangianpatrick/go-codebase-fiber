package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/sangianpatrick/go-codebase-fiber/config"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

var (
	db       *sql.DB
	syncOnce sync.Once
)

func constructPgx() *sql.DB {
	c := config.Get()

	purl := &url.URL{}
	purl.Scheme = "postgres"
	purl.User = url.UserPassword(c.Postgres.Username, c.Postgres.Password)
	purl.Host = fmt.Sprintf("%s:%s", c.Postgres.Host, c.Postgres.Port)
	purl.Path = c.Postgres.Database

	qs := purl.Query()
	qs.Set("search_path", c.Postgres.Schema)
	qs.Set("sslmode", c.Postgres.SSLMode)
	qs.Set("pool_min_conns", fmt.Sprint(c.Postgres.PoolMinConns))
	qs.Set("pool_max_conns", fmt.Sprint(c.Postgres.PoolMaxConns))
	qs.Set("pool_max_conn_lifetime", (time.Duration(c.Postgres.MaxConnLifetime) * time.Second).String())
	qs.Set("pool_max_conn_idle_time", (time.Duration(c.Postgres.MaxConnIdleTime) * time.Second).String())
	purl.RawQuery = qs.Encode()

	pgxConfig, err := pgxpool.ParseConfig(purl.String())
	if err != nil {
		log.Println(err)
		pgxConfig = &pgxpool.Config{}
	}

	pgxConfig.ConnConfig.Tracer = otelpgx.NewTracer(
		otelpgx.WithAttributes(
			semconv.DBSystemPostgreSQL,
		),
		otelpgx.WithIncludeQueryParameters(),
	)

	pgxConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxUUID.Register(conn.TypeMap())
		return nil
	}

	pgxConnPool, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		log.Println(err)
		pgxConnPool = &pgxpool.Pool{}
	}

	return stdlib.OpenDBFromPool(pgxConnPool)
}

func GetDatabase() *sql.DB {
	syncOnce.Do(func() {
		db = constructPgx()
	})

	return db
}
