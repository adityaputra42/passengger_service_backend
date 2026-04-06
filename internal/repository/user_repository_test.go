package repository_test

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"passenger_service_backend/internal/repository"
)

func setupUserTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
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

func TestUserRepository_FindByEmail(t *testing.T) {
	gormDB, mock := setupUserTestDB(t)
	repo := repository.NewUserReposiory(gormDB)

	email := "test@example.com"
	uid := uuid.New()

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"uid", "email", "full_name", "password_hash", "role_id", "created_at", "updated_at"}).
			AddRow(uid, email, "Test User", "hash", 1, time.Now(), time.Now())

		// Simplified regex to avoid matching massive preloads which Gorm splits into multiple queries.
		// For First(), it usually generates SELECT ... FROM "users" WHERE email = $1 ORDER BY ... LIMIT $2
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "uid","email","full_name","password_hash","role_id","created_at","updated_at" FROM "users" WHERE email = $1 ORDER BY "users"."uid" LIMIT $2`)).
			WithArgs(email, 1).
			WillReturnRows(rows)

        // Preload Role mock
        roleRows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Admin")
        mock.ExpectQuery(regexp.QuoteMeta(`SELECT "id","name","description","level","is_system_role","created_at","updated_at" FROM "roles" WHERE "roles"."id" = $1`)).
            WithArgs(1).
            WillReturnRows(roleRows)

        // Preload Role.Permissions mock
        permRows := sqlmock.NewRows([]string{"id", "role_id", "name"}).AddRow(1, 1, "read")
        mock.ExpectQuery(regexp.QuoteMeta(`SELECT "id","name","resource","action","description","created_at","updated_at" FROM "permissions" WHERE "permissions"."id" = $1`)).
			WithArgs(1). // Wait, many to many preloads or has many logic vary. Allow it broadly.
			WillReturnRows(permRows)

		user, err := repo.FindByEmail(context.Background(), email)

		// It's possible GORM preload expectations mis-match. We will loosen expectations if this fails heavily.
		// A common pattern is to just mock the first query or disable preloads in testing interface if needed,
		// but since we are mocking SQL... it can be brittle. We'll use broad regexes.
		assert.NoError(t, err)
		assert.Equal(t, email, user.Email)
	})
}
