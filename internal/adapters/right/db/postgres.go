package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	host     = "db"
	port     = "5432"
	user     = "postgres"
	password = "postgres"
	dbname   = "1337board"
)

type Postgres struct {
	db *sql.DB
	Repo
}

func NewPostgres() *Postgres {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := sql.Open("pgx", psqlInfo)
	if err != nil {
		log.Fatalf("Failed to open a DB connection: %v", err)
	}

	log.Println("Successfully connected to the database!")

	return &Postgres{
		db:   db,
		Repo: Repo{Conn: db},
	}
}

func (p *Postgres) Close() error {
	if err := p.db.Close(); err != nil {
		return err
	}
	return nil
}
