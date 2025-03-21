package testpackage

//go:generate ../goconst --type eventType --out data/eventType.gen.ts
type eventType string

const (
	// User has been added to the system
	// Indicates that somthing foo
	USER_INVITED eventType = "user_invited"
	// New task has been created on a project
	TASK_CREATED eventType = "task_created"
	// Update to an existing task
	TASK_UPDATED eventType = "task_updated"
	// New file has been uploaded to the system
	FILE_UPLOADED eventType = "file_uploaded"
)
