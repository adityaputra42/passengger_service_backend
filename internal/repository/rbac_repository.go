package repository

import (
	"passenger_service_backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RBACRepository interface {
	GetUserRole(uid uuid.UUID) (*models.Role, error)
	HasPermission(uid uuid.UUID, resource, action string) (bool, error)
	IsOwner(uid uuid.UUID, resource string, resourceID uint) (bool, error)
}

type RBACRepositoryImpl struct {
	db *gorm.DB
}

func NewRBACRepository(db *gorm.DB) RBACRepository {
	return &RBACRepositoryImpl{db: db}
}

// GetUserRole - Ambil role user
func (r *RBACRepositoryImpl) GetUserRole(uid uuid.UUID) (*models.Role, error) {
	var user models.User
	if err := r.db.
		Select("id", "role_id").
		Preload("Role", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "description", "level", "is_system_role", "created_at", "updated_at")
		}).
		First(&user, uid).Error; err != nil {
		return nil, err
	}
	return &user.Role, nil
}

// HasPermission - Cek permission via join table (FAST)
func (r *RBACRepositoryImpl) HasPermission(
	uid uuid.UUID,
	resource string,
	action string,
) (bool, error) {
	var count int64

	err := r.db.
		Table("users").
		Select("1"). // Only select constant 1 for counting
		Joins("JOIN roles ON roles.id = users.role_id").
		Joins("JOIN role_permissions ON role_permissions.role_id = roles.id").
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
		Where("users.uid = ?", uid).
		Where("(permissions.name = ? OR (permissions.resource = ? AND permissions.action = ?))",
			resource+"."+action, resource, action).
		Count(&count).Error

	return count > 0, err
}

// IsOwner - Cek ownership untuk *_own permission
func (r *RBACRepositoryImpl) IsOwner(
	uid uuid.UUID,
	resource string,
	resourceID uint,
) (bool, error) {
	switch resource {
	case "orders":
		var count int64
		err := r.db.
			Table("orders").
			Select("1").
			Where("id = ? AND uid = ?", resourceID, uid).
			Count(&count).Error
		return count > 0, err

	case "transactions":
		var count int64
		err := r.db.
			Table("transactions").
			Select("1").
			Where("tx_id = ? AND uid = ?", resourceID, uid).
			Count(&count).Error
		return count > 0, err

	default:
		return false, nil
	}
}
