package shared

// Role is a user role identifier. Keep this minimal; module-specific roles
// (if any) belong in the owning module.
type Role string

const (
	// RoleAdmin is a privileged role.
	RoleAdmin Role = "admin"
	// RoleUser is the default role.
	RoleUser Role = "user"
	// RoleService is reserved for machine/service accounts.
	RoleService Role = "service"
)

// Roles is the canonical ordered list of known roles.
var Roles = []Role{RoleAdmin, RoleUser, RoleService}

// IsValid reports whether the role is a known role.
func (r Role) IsValid() bool {
	for _, known := range Roles {
		if r == known {
			return true
		}
	}
	return false
}

// HasRole reports whether the given roles include this role.
func (r Role) HasRole(roles []Role) bool {
	for _, rl := range roles {
		if rl == r {
			return true
		}
	}
	return false
}
