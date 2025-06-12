package db

import (
	"context"
	"database/sql"
	"errors"
	"transaction-backend/domain"

	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(dsn string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	schema := `CREATE TABLE IF NOT EXISTS transactions (
		id UUID PRIMARY KEY,
		city TEXT,
		latitude DOUBLE PRECISION,
		longitude DOUBLE PRECISION,
		time TIMESTAMP,
		price INTEGER,
		type TEXT
	);`

	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	return &PostgresRepository{db: db}, nil
}

func (r *PostgresRepository) Save(ctx context.Context, txn domain.Transaction) error {
	if !txn.PropertyType.IsValid() {
		return errors.New("invalid property type")
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO transactions (id, city, latitude, longitude, time, price, type)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		txn.ID, txn.City, txn.Latitude, txn.Longitude, txn.Time, txn.Price, string(txn.PropertyType),
	)
	return err
}

func (r *PostgresRepository) GetAll(ctx context.Context) ([]domain.Transaction, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, city, latitude, longitude, time, price, type FROM transactions`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txns []domain.Transaction
	for rows.Next() {
		var txn domain.Transaction
		var typeStr string
		err := rows.Scan(&txn.ID, &txn.City, &txn.Latitude, &txn.Longitude, &txn.Time, &txn.Price, &typeStr)
		if err != nil {
			continue
		}
		if pt, err := domain.NewPropertyType(typeStr); err == nil {
			txn.PropertyType = pt
			txns = append(txns, txn)
		}
	}
	return txns, nil
}
