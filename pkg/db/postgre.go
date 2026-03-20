package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgreConf struct {
	host     string
	port     int
	user     string
	password string
	dbName   string
	schema   string
	sslMode  string
}

type PostgreOption func(*PostgreConf)

func WithHost(host string) PostgreOption {
	return func(pc *PostgreConf) {
		pc.host = host
	}
}

func WithPort(port int) PostgreOption {
	return func(pc *PostgreConf) {
		pc.port = port
	}
}

func WithUser(user string) PostgreOption {
	return func(pc *PostgreConf) {
		pc.user = user
	}
}

func WithPassword(password string) PostgreOption {
	return func(pc *PostgreConf) {
		pc.password = password
	}
}

func WithDbName(dbName string) PostgreOption {
	return func(pc *PostgreConf) {
		pc.dbName = dbName
	}
}

func WithSchema(schema string) PostgreOption {
	return func(pc *PostgreConf) {
		pc.schema = schema
	}
}

func WithSSLMode(sslMode string) PostgreOption {
	return func(pc *PostgreConf) {
		pc.sslMode = sslMode
	}
}

type postgreDB struct {
	db *gorm.DB
}

func NewPostgre(opts ...PostgreOption) (Database, error) {
	cfg := &PostgreConf{
		sslMode: "disable",
		schema:  "public",
	}

	for _, opt := range opts {
		opt(cfg)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s search_path=%s",
		cfg.host, cfg.port, cfg.user, cfg.password, cfg.dbName, cfg.sslMode, cfg.schema,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("db: failed to connect postgres: %w", err)
	}

	return &postgreDB{db: db}, nil
}

// Close implements [Database].
func (p *postgreDB) Close() error {
	sql, err := p.db.DB()
	if err != nil {
		return err
	}
	return sql.Close()
}

// GetDB implements [Database].
func (p *postgreDB) GetDB() *gorm.DB {
	return p.db
}

// Ping implements [Database].
func (p *postgreDB) Ping() error {
	sql, err := p.db.DB()
	if err != nil {
		return err
	}
	return sql.Ping()
}
