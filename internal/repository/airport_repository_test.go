package repository_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"passenger_service_backend/internal/models"
	"passenger_service_backend/internal/repository"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	require.NoError(t, err)

	return gormDB, mock
}

func TestAirportRepository_Create(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	repo := repository.NewAirportRepository(gormDB)

	id := uuid.New()
	airport := &models.Airport{
		ID:       id,
		Code:     "CGK",
		Name:     "Soekarno-Hatta",
		City:     "Jakarta",
		Country:  "Indonesia",
		Timezone: "Asia/Jakarta",
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "airports" ("code","name","city","country","timezone","id") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)).
			WithArgs(airport.Code, airport.Name, airport.City, airport.Country, airport.Timezone, airport.ID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(airport.ID))

		err := repo.Create(context.Background(), airport)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAirportRepository_FindByID(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	repo := repository.NewAirportRepository(gormDB)

	id := uuid.New()

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "code", "name", "city", "country", "timezone"}).
			AddRow(id, "CGK", "Soekarno-Hatta", "Jakarta", "Indonesia", "Asia/Jakarta")

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "airports" WHERE id = $1 ORDER BY "airports"."id" LIMIT $2`)).
			WithArgs(id, 1).
			WillReturnRows(rows)

		airport, err := repo.FindByID(context.Background(), id)

		assert.NoError(t, err)
		assert.NotNil(t, airport)
		assert.Equal(t, id, airport.ID)
		assert.Equal(t, "CGK", airport.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAirportRepository_FindByCode(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	repo := repository.NewAirportRepository(gormDB)

	id := uuid.New()

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "code", "name", "city", "country", "timezone"}).
			AddRow(id, "CGK", "Soekarno-Hatta", "Jakarta", "Indonesia", "Asia/Jakarta")

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "airports" WHERE code = $1 ORDER BY "airports"."id" LIMIT $2`)).
			WithArgs("CGK", 1).
			WillReturnRows(rows)

		airport, err := repo.FindByCode(context.Background(), "CGK")

		assert.NoError(t, err)
		assert.NotNil(t, airport)
		assert.Equal(t, "CGK", airport.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAirportRepository_FindAll(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	repo := repository.NewAirportRepository(gormDB)

	id1 := uuid.New()
	id2 := uuid.New()

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "code", "name", "city", "country", "timezone"}).
			AddRow(id1, "CGK", "Soekarno-Hatta", "Jakarta", "Indonesia", "Asia/Jakarta").
			AddRow(id2, "DPS", "Ngurah Rai", "Denpasar", "Indonesia", "Asia/Makassar")

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "airports" ORDER BY country, city`)).
			WillReturnRows(rows)

		airports, err := repo.FindAll(context.Background())

		assert.NoError(t, err)
		assert.Len(t, airports, 2)
		assert.Equal(t, "CGK", airports[0].Code)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAirportRepository_Search(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	repo := repository.NewAirportRepository(gormDB)

	id := uuid.New()

	t.Run("Success", func(t *testing.T) {
		query := "Jak"
		like := "%" + query + "%"
		
		rows := sqlmock.NewRows([]string{"id", "code", "name", "city", "country", "timezone"}).
			AddRow(id, "CGK", "Soekarno-Hatta", "Jakarta", "Indonesia", "Asia/Jakarta")

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "airports" WHERE code ILIKE $1 OR name ILIKE $2 OR city ILIKE $3 OR country ILIKE $4 LIMIT $5`)).
			WithArgs(like, like, like, like, 20).
			WillReturnRows(rows)

		airports, err := repo.Search(context.Background(), query)

		assert.NoError(t, err)
		assert.Len(t, airports, 1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAirportRepository_Update(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	repo := repository.NewAirportRepository(gormDB)

	id := uuid.New()
	airport := &models.Airport{
		ID:       id,
		Code:     "CGK",
		Name:     "Soekarno-Hatta Baru",
		City:     "Jakarta",
		Country:  "Indonesia",
		Timezone: "Asia/Jakarta",
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "airports" SET "code"=$1,"name"=$2,"city"=$3,"country"=$4,"timezone"=$5 WHERE "id" = $6`)).
			WithArgs(airport.Code, airport.Name, airport.City, airport.Country, airport.Timezone, airport.ID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Update(context.Background(), airport)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAirportRepository_Delete(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	repo := repository.NewAirportRepository(gormDB)

	id := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "airports" WHERE id = $1`)).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete(context.Background(), id)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
