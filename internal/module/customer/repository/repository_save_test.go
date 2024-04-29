package repository_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sangianpatrick/go-codebase-fiber/internal/module/customer/repository"
	"github.com/sangianpatrick/go-codebase-fiber/internal/pkg/applogger"
	"github.com/stretchr/testify/assert"
)

func TestSaveErrorWhileExecQuery(t *testing.T) {
	SetEnv(t)

	logger := applogger.GetZap()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	now := time.Now()
	c := NewCustomer(now)

	mock.ExpectBegin()
	mock.ExpectPrepare(insertQuery).
		WillBeClosed()
	mock.ExpectExec(insertQuery).
		WithArgs(
			c.ID, c.Email, c.Firstname, c.Lastname, c.VerificationStatus,
			c.MemberStatus, c.Password, c.PasswordSalt, c.CreatedAt, c.UpdatedAt,
		).
		WillReturnError(fmt.Errorf("exec error"))

	repo := repository.NewRepository(logger, db)
	ctx := context.TODO()
	tx, err := repo.BeginTx(ctx)
	if err != nil {
		t.Fatal(err)
	}

	err = repo.Save(ctx, c, tx)
	assert.Error(t, err)
	assert.Equal(t, "exec error", err.Error())

	mock.ExpectationsWereMet()
}

func TestSaveErrorWhilePreparingStatement(t *testing.T) {
	SetEnv(t)

	logger := applogger.GetZap()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	now := time.Now()
	c := NewCustomer(now)

	mock.ExpectBegin()
	mock.ExpectPrepare(insertQuery).
		WillBeClosed().
		WillReturnError(fmt.Errorf("prepare error"))

	repo := repository.NewRepository(logger, db)
	ctx := context.TODO()
	tx, err := repo.BeginTx(ctx)
	if err != nil {
		t.Fatal(err)
	}

	err = repo.Save(ctx, c, tx)
	assert.Error(t, err)
	assert.Equal(t, "prepare error", err.Error())

	mock.ExpectationsWereMet()
}
func TestSaveSuccess(t *testing.T) {
	SetEnv(t)

	logger := applogger.GetZap()

	now := time.Now()
	c := NewCustomer(now)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectPrepare(insertQuery).WillBeClosed()
	mock.ExpectExec(insertQuery).
		WithArgs(
			c.ID, c.Email, c.Firstname, c.Lastname, c.VerificationStatus,
			c.MemberStatus, c.Password, c.PasswordSalt, c.CreatedAt, c.UpdatedAt,
		).
		WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectCommit()

	repo := repository.NewRepository(logger, db)
	ctx := context.TODO()
	tx, err := repo.BeginTx(ctx)
	if err != nil {
		t.Fatal(err)
	}

	err = repo.Save(ctx, c, tx)
	if err != nil {
		t.Fatal(err)
	}

	err = repo.CommitTx(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectationsWereMet()
}
