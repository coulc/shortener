package storage

import (
	"database/sql"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)


var (
	URLExist = errors.New("The url already exists.")
	URLNotFound = errors.New("Cannot find this url.")
)

type SQLiteStorage struct {
	db *sql.DB
}

func NewSQLiteStorage(dbPath string) (*SQLiteStorage,error) {
	// dir := filepath.Dir(dbPath)
	// if err := os.MkdirAll(dir, 0755); err != nil {
	// 	slog.Error("Failed to create database directory", "err", err, "dir", dir)
	// 	return nil, err
	// }
	if err := EnsureDir(dbPath);err != nil {
		return  nil,err
	}

	db,err := sql.Open("sqlite3",dbPath)	
	if err != nil {
		return nil,err
	}

	if err = db.Ping(); err != nil {
		return nil,err
	}
	slog.Info("Database connection successful.")

	// sqlite configuration
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	sqlStr := `
	CREATE TABLE IF NOT EXISTS urls (
		short_code TEXT PRIMARY KEY,
		long_url TEXT NOT NULL UNIQUE,
		created_at INTEGER,
		visit_count INTEGER
	);
	`
	_,err = db.Exec(sqlStr)
	if err != nil {
		return nil,err
	}

	return &SQLiteStorage{ db: db },nil
}


func (s *SQLiteStorage) Save(m *URLMapping) error {
	sqlStr := "SELECT short_code, long_url, created_at, visit_count FROM urls WHERE long_url = ?"
	row := s.db.QueryRow(sqlStr,m.LongURL)

	var existing URLMapping
	err := row.Scan(&existing.ShortCode, &existing.LongURL, &existing.CreatedAt, &existing.VisitCount)
	if err == nil {
		m.ShortCode = existing.ShortCode
		m.LongURL = existing.LongURL
		m.CreatedAt = existing.CreatedAt
		return URLExist
	}

	sqlStr = "INSERT INTO urls(short_code, long_url, created_at, visit_count) VALUES (?, ?, ?, ?)"	

	_,err = s.db.Exec(sqlStr,m.ShortCode,m.LongURL,m.CreatedAt,m.VisitCount)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLiteStorage) Get(shortCode string) (*URLMapping,error) {
	sqlStr := "SELECT short_code, long_url, created_at, visit_count FROM urls WHERE short_code = ?"

	row:= s.db.QueryRow(sqlStr,shortCode)

	var m URLMapping

	err := row.Scan(&m.ShortCode,&m.LongURL,&m.CreatedAt,&m.VisitCount)
	if err != nil {
		return nil,URLNotFound
	}
	return &m,nil
}

func (s *SQLiteStorage)IncrementVisit(shortCode string) error {
	sqlStr := "UPDATE urls SET visit_count = visit_count + 1 WHERE short_code = ?"
	_,err := s.db.Exec(sqlStr,shortCode)
	if err != nil {
		return err
	}
	return nil
}


func (s *SQLiteStorage)Delete(shortCode string) error {
	sqlStr := "DELETE FROM urls WHERE short_code = ?"
	res,err := s.db.Exec(sqlStr,shortCode)
	if err != nil {
		return err
	}
	rows,_ := res.RowsAffected()
	if rows == 0 {
		return URLNotFound
	}
	return nil
}

func (s *SQLiteStorage)Close() error {
	return s.db.Close()
}

func EnsureDir(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		slog.Error("Failed to create database directory", "err", err, "dir", dir)
		return err
	}
	return nil
}

