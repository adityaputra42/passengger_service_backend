package repository

import (
	"context"
	"errors"
	"passenger_service_backend/internal/cache"
	"passenger_service_backend/internal/models"
	"time"

	"github.com/google/uuid"
)

// cachedRBACRepository wraps *RBACRepositoryImpl (concrete, not interface).
// Wire injects the concrete type so it has an unambiguous provider.
type CachedRBACRepository struct {
	inner *RBACRepositoryImpl // ← concrete, not RBACRepository interface
	cache *cache.Client
}

func NewCachedRBACRepository(
	inner *RBACRepositoryImpl,
	cache *cache.Client,
) *CachedRBACRepository {
	return &CachedRBACRepository{
		inner: inner,
		cache: cache,
	}
}

func (r *CachedRBACRepository) GetUserRole(uid uuid.UUID) (*models.Role, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	key := cache.KeyUserRole(uid)

	var role models.Role
	err := r.cache.Get(ctx, key, &role)
	if err == nil {
		return &role, nil
	}
	if !errors.Is(err, cache.ErrCacheMiss) {
		// Redis degraded — fall through to DB, don't fail the request
		_ = err
	}

	result, err := r.inner.GetUserRole(uid)
	if err != nil {
		return nil, err
	}

	_ = r.cache.Set(ctx, key, result, cache.TTLRBACRole*time.Second)
	return result, nil
}

func (r *CachedRBACRepository) HasPermission(uid uuid.UUID, resource, action string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	key := cache.KeyUserPermission(uid, resource, action)

	type permResult struct {
		Allowed bool `json:"allowed"`
	}

	var cached permResult
	err := r.cache.Get(ctx, key, &cached)
	if err == nil {
		return cached.Allowed, nil
	}
	if !errors.Is(err, cache.ErrCacheMiss) {
		_ = err
	}

	allowed, err := r.inner.HasPermission(uid, resource, action)
	if err != nil {
		return false, err
	}

	_ = r.cache.Set(ctx, key, permResult{Allowed: allowed}, cache.TTLRBACPermission*time.Second)
	return allowed, nil
}

func (r *CachedRBACRepository) IsOwner(uid uuid.UUID, resource string, resourceID uint) (bool, error) {
	// Intentionally not cached — see rbac_repository.go for reasoning
	return r.inner.IsOwner(uid, resource, resourceID)
}

// InvalidateUser removes all RBAC cache entries for a user.
// Call from RoleService.AssignPermissions and UserService.Update (role change).
func (r *CachedRBACRepository) InvalidateUser(uid uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Bust exact role key
	if err := r.cache.Del(ctx, cache.KeyUserRole(uid)); err != nil {
		return err
	}
	// Bust all permission keys for this user (dynamic resource+action combinations)
	return r.cache.DelPattern(ctx, "rbac:permission:"+uid.String()+":*")
}

var _ RBACRepository = (*CachedRBACRepository)(nil)
