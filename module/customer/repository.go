package customer

import (
	"context"
	"database/sql"

	"github.com/sangianpatrick/go-codebase-fiber/pkg/applogger"
	"github.com/sangianpatrick/go-codebase-fiber/pkg/errors"
	"go.uber.org/zap"
)

type Repository interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
	RollbackTx(ctx context.Context, tx *sql.Tx) error
	CommitTx(ctx context.Context, tx *sql.Tx) error

	Save(ctx context.Context, c Customer, tx *sql.Tx) error
	FindByEmail(ctx context.Context, email string, tx *sql.Tx) (Customer, error)
}

type sqlCommand interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

type repository struct {
	logger *applogger.ZapLogger
	db     *sql.DB
}

func NewRepository(logger *applogger.ZapLogger, db *sql.DB) Repository {
	return &repository{
		logger: logger,
		db:     db,
	}
}

// BeginTx implements Repository.
func (r *repository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Error(ctx, err.Error(), zap.Error(err))
		return nil, err
	}

	return tx, nil
}

// CommitTx implements Repository.
func (r *repository) CommitTx(ctx context.Context, tx *sql.Tx) error {
	if err := tx.Commit(); err != nil {
		r.logger.Error(ctx, err.Error(), zap.Error(err))
		return err
	}

	return nil
}

// RollbackTx implements Repository.
func (r *repository) RollbackTx(ctx context.Context, tx *sql.Tx) error {
	if err := tx.Rollback(); err != nil {
		r.logger.Error(ctx, err.Error(), zap.Error(err))
		return err
	}

	return nil
}

func (r *repository) query(ctx context.Context, query string, cmd sqlCommand, args ...interface{}) ([]Customer, error) {
	stmt, err := cmd.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := make([]Customer, 0)
	for rows.Next() {
		var c Customer
		err := rows.Scan(
			&c.ID, &c.Email, &c.Firstname, &c.Lastname, &c.VerificationStatus,
			&c.MemberStatus, &c.Password, &c.PasswordSalt, &c.CreatedAt, &c.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		data = append(data, c)
	}

	return data, nil
}

func (r *repository) exec(ctx context.Context, query string, cmd sqlCommand, args ...interface{}) error {
	stmt, err := cmd.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) FindByEmail(ctx context.Context, email string, tx *sql.Tx) (Customer, error) {
	var cmd sqlCommand = r.db

	if tx != nil {
		cmd = tx
	}

	query := `
		SELECT
			id, email, firstname, lastname, verification_status,
			member_status, password, password_salt, created_at, updated_at
		FROM customer
		WHERE
			email = $1
		LIMIT 1
	`

	data, err := r.query(ctx, query, cmd, email)
	if err != nil {
		r.logger.Error(ctx, "", zap.Error(err))
		return Customer{}, err
	}

	if len(data) < 1 {
		return Customer{}, errors.NotFound
	}

	return data[0], nil
}

// Save implements Repository.
func (r *repository) Save(ctx context.Context, c Customer, tx *sql.Tx) error {
	var cmd sqlCommand = r.db

	if tx != nil {
		cmd = tx
	}

	query := `
		INSERT INTO customer
		(
			id, email, firstname, lastname, verification_status, member_status, password, password_salt, created_at, updated_at
		)
		VALUES
		(
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)
	`

	err := r.exec(ctx, query, cmd,
		c.ID, c.Email, c.Firstname, c.Lastname, c.VerificationStatus,
		c.MemberStatus, c.Password, c.PasswordSalt, c.CreatedAt, c.UpdatedAt,
	)
	if err != nil {
		r.logger.Error(ctx, "", zap.Error(err))
		return err
	}

	return nil
}
