const roles = [
		// ADMIN is able to manage the account
	
	"ADMIN",
		// Full permissions on a project
	
	"PROJECT_MANAGER",
		// Manage and update parts of tasks
	
	"PROJECT_MEMBER",
] as const;

type Role = typeof roles[number];
