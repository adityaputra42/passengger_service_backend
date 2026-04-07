package services

import (
	"passenger_service_backend/internal/models"
	"passenger_service_backend/internal/repository"

	"github.com/google/uuid"
)

type RBACService interface {
	CheckPermission(uid uuid.UUID, resource, action string) (bool, error)
	CheckPermissionOrOwn(uid uuid.UUID, resource, action string, resourceID uint) (bool, error)

	HasRole(uid uuid.UUID, roleName string) (bool, error)  // hierarchy
	HasExactRole(uid uuid.UUID, roleId uint) (bool, error) // STRICT

	GetUserRole(uid uuid.UUID) (*models.Role, error)
	CanManageUser(managerID, targetUserID uuid.UUID) (bool, error)
}

type RBACServiceImpl struct {
	repo repository.RBACRepository
}

// HasExactRole implements [RBACService].
func (s *RBACServiceImpl) HasExactRole(uid uuid.UUID, roleId uint) (bool, error) {
	userRole, err := s.repo.GetUserRole(uid)
	if err != nil {
		return false, err
	}

	return userRole.ID == roleId, nil
}

func NewRBACService(repo repository.RBACRepository) RBACService {
	return &RBACServiceImpl{repo: repo}
}

// =======================
// PERMISSION
// =======================

func (s *RBACServiceImpl) CheckPermission(
	uid uuid.UUID,
	resource string,
	action string,
) (bool, error) {
	return s.repo.HasPermission(uid, resource, action)
}

func (s *RBACServiceImpl) CheckPermissionOrOwn(
	uid uuid.UUID,
	resource string,
	action string,
	resourceID uint,
) (bool, error) {

	if ok, _ := s.CheckPermission(uid, resource, action); ok {
		return true, nil
	}

	if ok, _ := s.CheckPermission(uid, resource, action+"_own"); ok {
		return s.repo.IsOwner(uid, resource, resourceID)
	}

	return false, nil
}

// =======================
// ROLE
// =======================

func (s *RBACServiceImpl) HasRole(uid uuid.UUID, roleName string) (bool, error) {
	userRole, err := s.repo.GetUserRole(uid)
	if err != nil {
		return false, err
	}

	userLevel := s.GetRoleHierarchyLevel(userRole.Name)
	requiredLevel := s.GetRoleHierarchyLevel(roleName)

	return userLevel >= requiredLevel, nil
}

func (s *RBACServiceImpl) GetUserRole(uid uuid.UUID) (*models.Role, error) {
	return s.repo.GetUserRole(uid)
}

// =======================
// USER MANAGEMENT
// =======================

func (s *RBACServiceImpl) CanManageUser(
	managerID uuid.UUID,
	targetuid uuid.UUID,
) (bool, error) {

	managerRole, err := s.GetUserRole(managerID)
	if err != nil {
		return false, err
	}

	targetRole, err := s.GetUserRole(targetuid)
	if err != nil {
		return false, err
	}

	return s.GetRoleHierarchyLevel(managerRole.Name) >
		s.GetRoleHierarchyLevel(targetRole.Name), nil
}

// =======================
// ROLE HIERARCHY
// =======================

func (s *RBACServiceImpl) GetRoleHierarchyLevel(roleName string) int {
	hierarchy := map[string]int{
		"Super Admin": 4,
		"Admin":       3,
		"Vendor":      2,
		"Customer":    1,
	}

	if level, ok := hierarchy[roleName]; ok {
		return level
	}
	return 0
}
