package testpackage

//go:generate ../goconst --type role --out data/role.gen.ts
type role string

const (
	// ADMIN is able to manage the account
	ADMIN role = "admin"
	// Full permissions on a project
	PROJECT_MANAGER role = "project_manager"
	// Manage and update parts of tasks
	PROJECT_MEMBER role = "project_member"
)
