package db

import (
	"context"
	"embed"
	"fmt"

	"my-little-olap/internal/utils"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/clickhouse"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

type (
	ClickhouseDB struct {
		conn   driver.Conn
		logger *utils.Logger
	}
	Config struct {
		Host     string
		DBName   string
		User     string
		Password string
		Logger   *utils.Logger
	}
)

//go:embed migrations/*.sql
var fs embed.FS

func NewClickhouseDB(c Config) (*ClickhouseDB, error) {
	err := migrateCH(&c)
	if err != nil {
		return nil, err
	}

	conn, err := openCH(&c)
	if err != nil {
		return nil, err
	}

	ch := ClickhouseDB{conn, c.Logger}
	return &ch, err
}

func migrateCH(c *Config) error {
	d, err := iofs.New(fs, "migrations")
	if err != nil {
		return err
	}

	dbURL := fmt.Sprintf(
		"clickhouse://%s?username=%s&password=%s&x-multi-statement=true",
		c.Host,
		c.User,
		c.Password,
	)

	m, err := migrate.NewWithSourceInstance("iofs", d, dbURL)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func openCH(c *Config) (driver.Conn, error) {
	conn, err := clickhouse.Open(
		&clickhouse.Options{
			Addr: []string{c.Host},
			Auth: clickhouse.Auth{
				Database: c.DBName,
				Username: c.User,
				Password: c.Password,
			},
		},
	)
	if err == nil {
		err = conn.Ping(context.Background())
	}
	return conn, err
}
