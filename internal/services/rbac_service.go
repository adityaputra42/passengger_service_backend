package services

import (
	"passenger_service_backend/internal/models"
	"passenger_service_backend/internal/repository"
)

type RBACService interface {
	CheckPermission(userID uint, resource, action string) (bool, error)
	CheckPermissionOrOwn(userID uint, resource, action string, resourceID uint) (bool, error)

	HasRole(userID uint, roleName string) (bool, error)      // hierarchy
	HasExactRole(userID uint, roleName string) (bool, error) // STRICT

	GetUserRole(userID uint) (*models.Role, error)
	CanManageUser(managerID, targetUserID uint) (bool, error)
}

type RBACServiceImpl struct {
	repo repository.RBACRepository
}

// HasExactRole implements [RBACService].
func (s *RBACServiceImpl) HasExactRole(userID uint, roleName string) (bool, error) {
	userRole, err := s.repo.GetUserRole(userID)
	if err != nil {
		return false, err
	}

	return userRole.Name == roleName, nil
}

func NewRBACService(repo repository.RBACRepository) RBACService {
	return &RBACServiceImpl{repo: repo}
}

// =======================
// PERMISSION
// =======================

func (s *RBACServiceImpl) CheckPermission(
	userID uint,
	resource string,
	action string,
) (bool, error) {
	return s.repo.HasPermission(userID, resource, action)
}

func (s *RBACServiceImpl) CheckPermissionOrOwn(
	userID uint,
	resource string,
	action string,
	resourceID uint,
) (bool, error) {

	if ok, _ := s.CheckPermission(userID, resource, action); ok {
		return true, nil
	}

	if ok, _ := s.CheckPermission(userID, resource, action+"_own"); ok {
		return s.repo.IsOwner(userID, resource, resourceID)
	}

	return false, nil
}

// =======================
// ROLE
// =======================

func (s *RBACServiceImpl) HasRole(userID uint, roleName string) (bool, error) {
	userRole, err := s.repo.GetUserRole(userID)
	if err != nil {
		return false, err
	}

	userLevel := s.GetRoleHierarchyLevel(userRole.Name)
	requiredLevel := s.GetRoleHierarchyLevel(roleName)

	return userLevel >= requiredLevel, nil
}

func (s *RBACServiceImpl) GetUserRole(userID uint) (*models.Role, error) {
	return s.repo.GetUserRole(userID)
}

// =======================
// USER MANAGEMENT
// =======================

func (s *RBACServiceImpl) CanManageUser(
	managerID uint,
	targetUserID uint,
) (bool, error) {

	managerRole, err := s.GetUserRole(managerID)
	if err != nil {
		return false, err
	}

	targetRole, err := s.GetUserRole(targetUserID)
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
