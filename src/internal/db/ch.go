package db

import (
	"context"
	"embed"
	"fmt"
	"strings"

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

func (ch *ClickhouseDB) insertBatch(
	table string,
	nextRow func() *[]any,
) error {
	ctx := context.Background()

	row := nextRow()
	firstRowLen := len(*row)
	phs := make([]string, firstRowLen)
	for i := 0; i < firstRowLen; i++ {
		phs[i] = "?"
	}
	query := fmt.Sprintf(
		"INSERT INTO %s VALUES (%s)",
		table,
		strings.Join(phs, ","),
	)

	batch, err := ch.conn.PrepareBatch(ctx, query)
	if err != nil {
		return err
	}

	for row != nil {
		rLen := len(*row)
		if rLen != firstRowLen {
			ch.logger.Error.Printf(
				"Row length: %d mismatches with first row length: %d\n",
				rLen,
				firstRowLen,
			)
		}
		err := batch.Append(*row...)
		if err != nil {
			ch.logger.Error.Printf(
				"Error adding to batch in: %s: %s\n",
				table,
				err,
			)
		}
		row = nextRow()
	}

	return batch.Send()
}
