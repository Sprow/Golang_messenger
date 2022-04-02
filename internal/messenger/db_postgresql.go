package messenger

import (
	"database/sql"
	_ "embed"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

const (
	HOST = "localhost"
	PORT = 5432
)

type Config struct {
	Username string
	Password string
	DBName   string
}

type PostgreSqlDB struct {
	Conn *sql.DB
}

func Initialize(cfg Config) (PostgreSqlDB, error) {
	log.Println("Initializing")
	db := PostgreSqlDB{}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Username, cfg.Password, HOST, PORT, cfg.DBName)
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return db, err
	}

	db.Conn = conn
	err = db.Conn.Ping()
	if err != nil {
		return db, err
	}
	log.Println("Database connection established")
	err = migrate(db)
	if err != nil {
		log.Println(err)
	}
	log.Println("migrate complete")
	return db, nil
}

//go:embed migrate_postgresql.sql
var migrationSql string

func migrate(db PostgreSqlDB) error {
	_, err := db.Conn.Exec(migrationSql)
	return err
}
