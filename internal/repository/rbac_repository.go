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

func NewRBACRepository(db *gorm.DB) *RBACRepositoryImpl {
	return &RBACRepositoryImpl{db: db}
}
func (r *RBACRepositoryImpl) GetUserRole(uid uuid.UUID) (*models.Role, error) {
	var user models.User

	if err := r.db.
		Select("uid", "role_id").
		Preload("Role", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "description", "level", "is_system_role", "created_at", "updated_at").
				Where("deleted_at IS NULL")
		}).
		Where("uid = ?", uid).
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user.Role, nil
}

func (r *RBACRepositoryImpl) HasPermission(
	uid uuid.UUID,
	resource string,
	action string,
) (bool, error) {
	var count int64

	err := r.db.Model(&models.User{}).
		Joins("JOIN roles ON roles.id = users.role_id AND roles.deleted_at IS NULL").
		Joins("JOIN role_permissions ON role_permissions.role_id = roles.id").
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id AND permissions.deleted_at IS NULL").
		Where("users.uid = ?", uid).
		Where(
			"permissions.name = ? OR (permissions.resource = ? AND permissions.action = ?)",
			resource+"."+action, resource, action,
		).
		Count(&count).Error

	return count > 0, err
}

func (r *RBACRepositoryImpl) IsOwner(
	uid uuid.UUID,
	resource string,
	resourceID uint,
) (bool, error) {
	switch resource {
	case "pnr", "booking":
		var count int64
		err := r.db.Model(&models.PNRContact{}).
			Joins("JOIN users ON users.email = pnr_contacts.email").
			Where("users.uid = ? AND pnr_contacts.pnr_id = ?", uid, resourceID).
			Count(&count).Error
		return count > 0, err

	case "passenger":
		var count int64
		err := r.db.Model(&models.PNRPassenger{}).
			Joins("JOIN pnrs ON pnrs.id = pnr_passengers.pnr_id").
			Joins("JOIN pnr_contacts ON pnr_contacts.pnr_id = pnrs.id").
			Joins("JOIN users ON users.email = pnr_contacts.email").
			Where("users.uid = ? AND pnr_passengers.id = ?", uid, resourceID).
			Count(&count).Error
		return count > 0, err

	default:
		return false, nil
	}
}
