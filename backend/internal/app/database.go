package app

import (
	"backend/pkg/common"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"backend/pkg/config"
)

func createDB(cfg *config.Config) (*sqlx.DB, func(), error) {
	switch cfg.DB.Driver {
	case "sqlite":
		return createSQLiteDB(cfg)
	case "mysql":
		return createMySQLDB(cfg)
	case "postgres":
		return createPostgresDB(cfg)
	default:
		return nil, nil, fmt.Errorf("unsupported db driver: %s", cfg.DB.Driver)
	}
}

func createSQLiteDB(cfg *config.Config) (*sqlx.DB, func(), error) {
	if err := common.EnsurePath(cfg.DB.SQLite.Path); err != nil {
		return nil, nil, err
	}
	dsn := fmt.Sprintf("file:%s?_foreign_keys=off", cfg.DB.SQLite.Path)

	db, err := sqlx.Open("sqlite3", dsn)
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		db.Close()
	}
	pragmas := []string{
		"PRAGMA journal_mode = WAL;",
		"PRAGMA synchronous = NORMAL;",
		"PRAGMA foreign_keys = ON;",
		"PRAGMA temp_store = MEMORY;",
	}

	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			return nil, cleanup, err
		}
	}

	// SQLite 强烈建议
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if err := db.Ping(); err != nil {
		return nil, cleanup, err
	}

	return db, cleanup, nil
}

func createMySQLDB(cfg *config.Config) (*sqlx.DB, func(), error) {
	c := cfg.DB.MySQL

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
		c.Params,
	)

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		db.Close()
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)

	if err := db.Ping(); err != nil {
		return nil, cleanup, err
	}

	return db, cleanup, nil
}

func createPostgresDB(cfg *config.Config) (*sqlx.DB, func(), error) {
	c := cfg.DB.Postgres

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.DBName,
		c.SSLMode,
	)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		db.Close()
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)

	if err := db.Ping(); err != nil {
		return nil, cleanup, err
	}

	return db, cleanup, nil
}

func runMigrate(db *sqlx.DB, cfg *config.Config) error {
	m := cfg.DB.Migrate
	if !m.Enabled {
		return nil
	}

	filename := fmt.Sprintf("%s.sql", cfg.DB.Driver)
	path := filepath.Join(m.Dir, filename)

	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read migrate file failed: %w", err)
	}

	sql := strings.TrimSpace(string(content))
	if sql == "" {
		return nil
	}

	if _, err := db.Exec(sql); err != nil {
		return fmt.Errorf("exec migrate sql failed: %w", err)
	}

	return nil
}
