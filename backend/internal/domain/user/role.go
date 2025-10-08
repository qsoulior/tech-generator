package user_domain

type Role string

const (
	RoleRead     Role = "read"
	RoleWrite    Role = "write"
	RoleMaintain Role = "maintain"
)

var roleSet = map[Role]struct{}{
	RoleRead:     {},
	RoleWrite:    {},
	RoleMaintain: {},
}

func (r Role) Valid() bool {
	_, found := roleSet[r]
	return found
}
