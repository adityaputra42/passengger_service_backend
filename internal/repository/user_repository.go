package repository

import (
	"context"
	"passenger_service_backend/internal/db"
	"passenger_service_backend/internal/dto"
	"passenger_service_backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, param models.User) (models.User, error)
	Update(ctx context.Context, param *models.User) (models.User, error)
	Delete(ctx context.Context, param models.User) error
	FindByUid(ctx context.Context, id uuid.UUID) (models.User, error)
	FindByEmail(ctx context.Context, email string) (models.User, error)
	FindByUsernameOrEmail(ctx context.Context, identifier string) (models.User, error)
	FindAll(ctx context.Context, param dto.UserListRequest) (*dto.UserListResponse, error)
}

type UserRepositoryImpl struct{}

func (u *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (models.User, error) {
	user := models.User{}

	err := db.DB.WithContext(ctx).
		Select("uid", "email", "full_name", "password_hash", "role_id", "created_at", "updated_at").
		Preload("Role", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "description", "level", "is_system_role", "created_at", "updated_at")
		}).
		Preload("Role.Permissions", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "resource", "action", "description", "created_at", "updated_at")
		}).
		Where("email = ?", email).
		First(&user).Error

	return user, err
}

func (u *UserRepositoryImpl) FindByUsernameOrEmail(ctx context.Context, identifier string) (models.User, error) {
	user := models.User{}

	err := db.DB.WithContext(ctx).
		Select("uid", "email", "full_name", "password_hash", "role_id", "created_at", "updated_at").
		Preload("Role", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "description", "level", "is_system_role", "created_at", "updated_at")
		}).
		Preload("Role.Permissions", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "resource", "action", "description", "created_at", "updated_at")
		}).
		Where("email = ?", identifier).
		First(&user).Error

	return user, err
}

func (u *UserRepositoryImpl) Create(ctx context.Context, param models.User) (models.User, error) {
	var result models.User

	db := db.DB.WithContext(ctx)

	if err := db.Create(&param).Error; err != nil {
		return result, err
	}

	err := db.
		Select("uid", "email", "full_name", "password_hash", "role_id", "created_at", "updated_at").
		Preload("Role", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "description", "level", "is_system_role", "created_at", "updated_at")
		}).
		Preload("Role.Permissions", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "resource", "action", "description", "created_at", "updated_at")
		}).
		First(&result, param.UID).Error

	return result, err
}

func (u *UserRepositoryImpl) Delete(ctx context.Context, param models.User) error {
	return db.DB.WithContext(ctx).Delete(&param).Error
}

func (u *UserRepositoryImpl) FindAll(ctx context.Context, param dto.UserListRequest) (*dto.UserListResponse, error) {
	offset := (param.Page - 1) * param.Limit

	query := db.DB.WithContext(ctx).
		Model(&models.User{}).
		Select("uid", "email", "full_name", "password_hash", "role_id", "created_at", "updated_at").
		Preload("Role", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "description", "level", "is_system_role", "created_at", "updated_at")
		})

	if param.SortBy == "" {
		param.SortBy = "created_at desc"
	}

	query = query.Order(param.SortBy)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var users []models.User
	if err := query.Offset(offset).Limit(param.Limit).Find(&users).Error; err != nil {
		return nil, err
	}

	result := dto.ToUserListResponse(users, total, param.Page, param.Limit)
	return result, nil
}

func (u *UserRepositoryImpl) FindByUid(ctx context.Context, id uuid.UUID) (models.User, error) {
	var user models.User
	err := db.DB.WithContext(ctx).
		Select("uid", "email", "full_name", "password_hash", "role_id", "created_at", "updated_at").
		Preload("Role", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "description", "level", "is_system_role", "created_at", "updated_at")
		}).
		Preload("Role.Permissions", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "resource", "action", "description", "created_at", "updated_at")
		}).
		First(&user, id).Error
	return user, err
}

/*
============================

	UPDATE USER

============================
*/
func (u *UserRepositoryImpl) Update(ctx context.Context, param *models.User) (models.User, error) {
	var result models.User

	db := db.DB.WithContext(ctx)

	if err := db.Save(param).Error; err != nil {
		return result, err
	}

	err := db.
		Select("uid", "email", "full_name", "password_hash", "role_id", "created_at", "updated_at").
		Preload("Role", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "description", "level", "is_system_role", "created_at", "updated_at")
		}).
		Preload("Role.Permissions", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "resource", "action", "description", "created_at", "updated_at")
		}).
		First(&result, param.UID).Error

	return result, err
}

/*
============================

	CONSTRUCTOR

============================
*/
func NewUserReposiory() UserRepository {
	return &UserRepositoryImpl{}
}
