package repository

import (
	"math"
	"passenger_service_backend/internal/db"
	"passenger_service_backend/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(param models.User) (models.User, error)
	Update(param *models.User) (models.User, error)
	Delete(param models.User) error
	FindById(id uint) (models.User, error)
	FindByEmail(email string) (models.User, error)
	FindByUsernameOrEmail(identifier string) (models.User, error)
	FindAll(param models.UserListRequest) (*models.UserListResponse, error)
}

type UserRepositoryImpl struct{}

/*
============================

	FIND BY EMAIL

============================
*/
func (u *UserRepositoryImpl) FindByEmail(email string) (models.User, error) {
	user := models.User{}

	err := db.DB.
		Select("id", "username", "email", "password_hash", "first_name", "last_name", "role_id", "is_active", "created_at", "updated_at").
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

func (u *UserRepositoryImpl) FindByUsernameOrEmail(identifier string) (models.User, error) {
	user := models.User{}

	err := db.DB.
		Select("id", "username", "email", "password_hash", "first_name", "last_name", "role_id", "is_active", "created_at", "updated_at").
		Preload("Role", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "description", "level", "is_system_role", "created_at", "updated_at")
		}).
		Preload("Role.Permissions", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "resource", "action", "description", "created_at", "updated_at")
		}).
		Where("email = ? OR username = ?", identifier, identifier).
		First(&user).Error

	return user, err
}

/*
============================

	CREATE USER

============================
*/
func (u *UserRepositoryImpl) Create(param models.User) (models.User, error) {
	var result models.User

	db := db.DB

	if err := db.Create(&param).Error; err != nil {
		return result, err
	}

	err := db.
		Select("id", "username", "email", "password_hash", "first_name", "last_name", "role_id", "is_active", "created_at", "updated_at").
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

	DELETE USER

============================
*/
func (u *UserRepositoryImpl) Delete(param models.User) error {
	return db.DB.Delete(&param).Error
}

/*
============================

	FIND ALL (PAGINATION)

============================
*/
func (u *UserRepositoryImpl) FindAll(param models.UserListRequest) (*models.UserListResponse, error) {
	offset := (param.Page - 1) * param.Limit

	query := db.DB.
		Model(&models.User{}).
		Select("id", "username", "email", "first_name", "last_name", "role_id", "is_active", "created_at", "updated_at").
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

	userResponses := make([]models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = *user.ToResponse()
	}

	totalPages := int(math.Ceil(float64(total) / float64(param.Limit)))

	return &models.UserListResponse{
		Users:      userResponses,
		Total:      total,
		Page:       param.Page,
		Limit:      param.Limit,
		TotalPages: totalPages,
	}, nil
}

/*
============================

	FIND BY ID

============================
*/
func (u *UserRepositoryImpl) FindById(id uint) (models.User, error) {
	var user models.User
	err := db.DB.
		Select("id", "username", "email", "password_hash", "first_name", "last_name", "role_id", "is_active", "created_at", "updated_at").
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
func (u *UserRepositoryImpl) Update(param *models.User) (models.User, error) {
	var result models.User

	db := db.DB

	if err := db.Save(param).Error; err != nil {
		return result, err
	}

	err := db.
		Select("id", "username", "email", "password_hash", "first_name", "last_name", "role_id", "is_active", "created_at", "updated_at").
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
