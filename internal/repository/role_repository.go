package repository

import (
	"errors"
	"passenger_service_backend/internal/db"
	"passenger_service_backend/internal/models"

	"gorm.io/gorm"
)

type RoleRepository interface {
	Create(param models.Role) (models.Role, error)
	Update(param models.Role) (models.Role, error)
	Delete(param models.Role) error
	FindById(paramId uint) (models.Role, error)
	FindByName(name string) (models.Role, error)
	FindByNameAndId(name string, id uint) (*models.Role, error)
	AddPermission(id uint, permissions *[]models.Permission) (*models.Role, error)
	UpdatePermission(id uint, permissions *[]models.Permission) (*models.Role, error)
	FindAll() ([]models.Role, error)
}

type RoleRepositoryImpl struct {
}

// UpdatePermission implements RoleRepository.
func (a *RoleRepositoryImpl) UpdatePermission(id uint, permissions *[]models.Permission) (*models.Role, error) {
	role := models.Role{}
	err := db.DB.
		Select("id", "name", "description", "level", "is_system_role", "created_at", "updated_at").
		First(&role, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	if err := db.DB.Model(&role).Association("Permissions").Replace(&permissions); err != nil {
		return nil, err
	}

	// Reload dengan permissions
	err = db.DB.
		Select("id", "name", "description", "level", "is_system_role", "created_at", "updated_at").
		Preload("Permissions", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "resource", "action", "description", "created_at", "updated_at")
		}).
		First(&role, id).Error

	return &role, err
}

// FindByNameAndId implements RoleRepository.
func (a *RoleRepositoryImpl) FindByNameAndId(name string, id uint) (*models.Role, error) {
	var existingRole models.Role
	err := db.DB.
		Select("id", "name", "description", "level", "is_system_role", "created_at", "updated_at").
		Where("name = ? AND id != ?", name, id).
		First(&existingRole).Error

	if err == nil {
		return nil, errors.New("role with this name already exists")
	}

	if err == gorm.ErrRecordNotFound {
		return nil, nil // Tidak ada role dengan nama yang sama, ini valid
	}

	return nil, err
}

// AddPermission implements RoleRepository.
func (a *RoleRepositoryImpl) AddPermission(id uint, permissions *[]models.Permission) (*models.Role, error) {
	role := models.Role{}
	err := db.DB.
		Select("id", "name", "description", "level", "is_system_role", "created_at", "updated_at").
		First(&role, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	if err := db.DB.Model(&role).Association("Permissions").Append(&permissions); err != nil {
		return nil, err
	}

	// Reload dengan permissions
	err = db.DB.
		Select("id", "name", "description", "level", "is_system_role", "created_at", "updated_at").
		Preload("Permissions", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "resource", "action", "description", "created_at", "updated_at")
		}).
		First(&role, id).Error

	return &role, err
}

// FindByName implements RoleRepository.
func (a *RoleRepositoryImpl) FindByName(name string) (models.Role, error) {
	role := models.Role{}
	err := db.DB.
		Select("id", "name", "description", "level", "is_system_role", "created_at", "updated_at").
		First(&role, "name = ?", name).Error

	return role, err
}

// Create implements RoleRepository.
func (a *RoleRepositoryImpl) Create(param models.Role) (models.Role, error) {
	var result models.Role

	db := db.DB

	err := db.Create(&param).Error
	if err != nil {
		return result, err
	}

	err = db.
		Select("id", "name", "description", "level", "is_system_role", "created_at", "updated_at").
		First(&result, param.ID).Error
	return result, err
}

// Delete implements RoleRepository.
func (a *RoleRepositoryImpl) Delete(param models.Role) error {
	return db.DB.Delete(&param).Error
}

// FindAll implements RoleRepository.
func (a *RoleRepositoryImpl) FindAll() ([]models.Role, error) {
	var roles []models.Role
	if err := db.DB.
		Select("id", "name", "description", "level", "is_system_role", "created_at", "updated_at").
		Preload("Permissions", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "resource", "action", "description", "created_at", "updated_at")
		}).
		Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// FindById implements RoleRepository.
func (a *RoleRepositoryImpl) FindById(paramId uint) (models.Role, error) {
	role := models.Role{}
	err := db.DB.
		Select("id", "name", "description", "level", "is_system_role", "created_at", "updated_at").
		Preload("Permissions", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "resource", "action", "description", "created_at", "updated_at")
		}).
		First(&role, "id = ?", paramId).Error

	return role, err
}

// Update implements RoleRepository.
func (a *RoleRepositoryImpl) Update(param models.Role) (models.Role, error) {
	var result models.Role

	db := db.DB

	err := db.Model(&param).Updates(param).Error
	if err != nil {
		return result, err
	}

	err = db.
		Select("id", "name", "description", "level", "is_system_role", "created_at", "updated_at").
		First(&result, param.ID).Error
	return result, err
}

func NewRoleRepository() RoleRepository {
	return &RoleRepositoryImpl{}
}
