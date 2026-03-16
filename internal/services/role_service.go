package services

import (
	"errors"
	"passenger_service_backend/internal/db"
	"passenger_service_backend/internal/dto"
	"passenger_service_backend/internal/models"
	"passenger_service_backend/internal/repository"

	"gorm.io/gorm"
)

type RoleService interface {
	FindAllRole() (*[]models.Role, error)
	FindById(id uint) (*models.Role, error)
	CreateRole(req *dto.RoleInput) (*models.Role, error)
	UpdateRole(id uint, req *dto.RoleInput) (*models.Role, error)
	DeleteRole(id uint) error
	GetPermissions() (*[]models.Permission, error)
	AssignPermissions(roleID uint, permissionIDs []uint) error
	GetRolePermissions(roleID uint) ([]*models.Permission, error)
}

type RoleServiceImpl struct {
	roleRepo       repository.RoleRepository
	permissionRepo repository.PermissionRepository
}

// AssignPermissions implements RoleService.
func (r *RoleServiceImpl) AssignPermissions(roleID uint, permissionIDs []uint) error {
	role, err := r.roleRepo.FindById(roleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("role not found")
		}
		return err
	}

	permissions, err := r.permissionRepo.FindAllById(permissionIDs)
	if err != nil {
		return err
	}

	if len(*permissions) != len(permissionIDs) {
		return errors.New("some permissions not found")
	}
	_, err = r.roleRepo.UpdatePermission(role.ID, permissions)
	if err != nil {
		return err
	}

	return nil
}

// CreateRole implements RoleService.
func (r *RoleServiceImpl) CreateRole(req *dto.RoleInput) (*models.Role, error) {
	_, err := r.roleRepo.FindByName(req.Name)
	if err == nil {
		return nil, errors.New("role with this name already exists")
	}

	role := models.Role{
		Name:         req.Name,
		Description:  req.Description,
		IsSystemRole: false,
	}
	newRole, err := r.roleRepo.Create(role)
	if err != nil {
		return nil, err
	}

	if len(req.PermissionIDs) > 0 {
		permissions, err := r.permissionRepo.FindAllById(req.PermissionIDs)
		if err != nil {
			return nil, err
		}
		_, err = r.roleRepo.AddPermission(newRole.ID, permissions)
		if err != nil {
			return nil, err
		}

	}
	roleData, err := r.roleRepo.FindById(newRole.ID)
	if err != nil {
		return nil, err
	}

	return &roleData, nil
}

// DeleteUser implements RoleService.
func (r *RoleServiceImpl) DeleteRole(id uint) error {

	role, err := r.roleRepo.FindById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("role not found")
		}
		return err
	}

	if role.IsSystemRole {
		return errors.New("cannot delete system role")
	}

	var userCount int64
	if err := db.DB.Model(&models.User{}).Where("role_id = ?", id).Count(&userCount).Error; err != nil {
		return err
	}

	if userCount > 0 {
		return errors.New("cannot delete role that is assigned to users")
	}

	if err := db.DB.Model(&role).Association("Permissions").Clear(); err != nil {
		return err
	}

	return r.roleRepo.Delete(role)
}

// FindAllRole implements RoleService.
func (r *RoleServiceImpl) FindAllRole() (*[]models.Role, error) {
	roles, err := r.roleRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return &roles, nil
}

// FindById implements RoleService.
func (r *RoleServiceImpl) FindById(id uint) (*models.Role, error) {
	role, err := r.roleRepo.FindById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}
	return &role, nil

}

// GetPermissions implements RoleService.
func (r *RoleServiceImpl) GetPermissions() (*[]models.Permission, error) {
	permissions, err := r.permissionRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return &permissions, nil
}

// GetRolePermissions implements RoleService.
func (r *RoleServiceImpl) GetRolePermissions(roleID uint) ([]*models.Permission, error) {
	role, err := r.roleRepo.FindById(roleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err

	}

	return role.Permissions, nil
}

// UpdateRole implements RoleService.
func (r *RoleServiceImpl) UpdateRole(id uint, req *dto.RoleInput) (*models.Role, error) {

	role, err := r.roleRepo.FindById(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}

	if req.Name != "" && req.Name != role.Name {
		if role.IsSystemRole {
			return nil, errors.New("cannot change system role name")
		}

		_, err := r.roleRepo.FindByNameAndId(req.Name, id)
		if err == nil {
			return nil, err
		}
		role.Name = req.Name
	}

	if req.Description != "" {
		role.Description = req.Description
	}

	if err := db.DB.Save(&role).Error; err != nil {
		return nil, err
	}

	if len(req.PermissionIDs) > 0 {

		permissions, err := r.permissionRepo.FindAllById(req.PermissionIDs)
		if err != nil {
			return nil, err
		}
		_, err = r.roleRepo.UpdatePermission(role.ID, permissions)
		if err != nil {
			return nil, err
		}
	}

	newRole, err := r.roleRepo.FindById(role.ID)

	if err != nil {
		return nil, err
	}

	return &newRole, nil
}

func NewRoleService(roleRepo repository.RoleRepository,
	permissionRepo repository.PermissionRepository,
) RoleService {
	return &RoleServiceImpl{roleRepo: roleRepo, permissionRepo: permissionRepo}

}
