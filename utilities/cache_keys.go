package utilities

import "fmt"

// Cache keys constants (kept minimal and generic)
const (
	UserCachePrefix = "user:"
)

// UserCacheKey builds the cache key for a user entity by ID.
func UserCacheKey(userID uint) string {
	return fmt.Sprintf("%s%d", UserCachePrefix, userID)
}
