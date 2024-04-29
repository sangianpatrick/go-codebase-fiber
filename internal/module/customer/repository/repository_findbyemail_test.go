package repository_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sangianpatrick/go-codebase-fiber/internal/module/customer/repository"
	"github.com/sangianpatrick/go-codebase-fiber/internal/pkg/applogger"
	"github.com/sangianpatrick/go-codebase-fiber/internal/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestFindByEmailSuccess(t *testing.T) {
	SetEnv(t)

	logger := applogger.GetZap()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	c := NewCustomer(time.Now())
	mock.ExpectBegin()
	mock.ExpectPrepare(findByEmailQuery).WillBeClosed()
	mock.ExpectQuery(findByEmailQuery).
		WithArgs(email).
		RowsWillBeClosed().
		WillReturnRows(
			sqlmock.NewRows(
				[]string{
					"id", "email", "firstname", "lastname", "verification_status",
					"member_status", "password", "password_salt", "created_at", "updated_at",
				},
			).AddRow(
				c.ID, c.Email, c.Firstname, c.Lastname, c.VerificationStatus,
				c.MemberStatus, c.Password, c.PasswordSalt, c.CreatedAt, c.UpdatedAt,
			),
		)

	repo := repository.NewRepository(logger, db)
	ctx := context.TODO()
	tx, err := repo.BeginTx(ctx)
	if err != nil {
		t.Fatal(err)
	}
	cresp, err := repo.FindByEmail(ctx, email, tx)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, c.ID, cresp.ID)
	mock.ExpectationsWereMet()
}

func TestFindByEmailErrorWhilePreparingStatement(t *testing.T) {
	SetEnv(t)

	logger := applogger.GetZap()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectPrepare(findByEmailQuery).WillReturnError(fmt.Errorf("error prepare find by email"))

	repo := repository.NewRepository(logger, db)
	ctx := context.TODO()
	tx, err := repo.BeginTx(ctx)
	if err != nil {
		t.Fatal(err)
	}
	_, err = repo.FindByEmail(ctx, email, tx)
	assert.Error(t, err)
	assert.Equal(t, "error prepare find by email", err.Error())
	mock.ExpectationsWereMet()
}

func TestFindByEmailErroConnectionDone(t *testing.T) {
	SetEnv(t)

	logger := applogger.GetZap()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectPrepare(findByEmailQuery).WillBeClosed()
	mock.ExpectQuery(findByEmailQuery).
		WithArgs(email).
		WillReturnError(sql.ErrConnDone)

	repo := repository.NewRepository(logger, db)
	ctx := context.TODO()
	tx, err := repo.BeginTx(ctx)
	if err != nil {
		t.Fatal(err)
	}
	_, err = repo.FindByEmail(ctx, email, tx)
	assert.Error(t, err)
	assert.EqualError(t, sql.ErrConnDone, err.Error())
	mock.ExpectationsWereMet()
}

func TestFindByEmailErroNotFound(t *testing.T) {
	SetEnv(t)

	logger := applogger.GetZap()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectPrepare(findByEmailQuery).WillBeClosed()
	mock.ExpectQuery(findByEmailQuery).
		WithArgs(email).
		WillReturnRows(
			sqlmock.NewRows(
				[]string{
					"id", "email", "firstname", "lastname", "verification_status",
					"member_status", "password", "password_salt", "created_at", "updated_at",
				},
			),
		)

	repo := repository.NewRepository(logger, db)
	ctx := context.TODO()
	tx, err := repo.BeginTx(ctx)
	if err != nil {
		t.Fatal(err)
	}
	_, err = repo.FindByEmail(ctx, email, tx)
	assert.Error(t, err)
	assert.EqualError(t, errors.NotFound, err.Error())
	mock.ExpectationsWereMet()
}

func TestFindByEmailErrorWhenScanning(t *testing.T) {
	SetEnv(t)

	logger := applogger.GetZap()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	c := NewCustomer(time.Now())
	mock.ExpectBegin()
	mock.ExpectPrepare(findByEmailQuery).WillBeClosed()
	mock.ExpectQuery(findByEmailQuery).
		WithArgs(email).
		RowsWillBeClosed().
		WillReturnRows(
			sqlmock.NewRows(
				[]string{
					"id", "email", "firstname", "lastname", "verification_status",
					"member_status", "password", "password_salt", "created_at", "updated_at",
				},
			).AddRow(
				1, 2, c.Firstname, c.Lastname, c.VerificationStatus,
				c.MemberStatus, c.Password, c.PasswordSalt, c.CreatedAt, c.UpdatedAt,
			),
		)

	repo := repository.NewRepository(logger, db)
	ctx := context.TODO()
	tx, err := repo.BeginTx(ctx)
	if err != nil {
		t.Fatal(err)
	}
	_, err = repo.FindByEmail(ctx, email, tx)

	assert.Error(t, err)

	mock.ExpectationsWereMet()
}
