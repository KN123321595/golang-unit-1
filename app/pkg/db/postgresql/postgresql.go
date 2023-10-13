package postgresql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func GetDBConnection(user, password, host, port, dbname string) (*sqlx.DB, error) {
	connString := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", user, password, host, port, dbname)

	db, err := sqlx.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("error connection to database: %w", err)
	}

	return db, nil
}
