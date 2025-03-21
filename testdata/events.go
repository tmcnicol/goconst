package testdata

//go:generate goconst --type eventType
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

type anotherType string

const (
	// User has been added to the system
	USER_INVITED anotherType = "user_invited"
	// New task has been created on a project
	TASK_CREATED anotherType = "task_created"
	// Update to an existing task
	TASK_UPDATED anotherType = "task_updated"
	// New file has been uploaded to the system
	FILE_UPLOADED anotherType = "file_uploaded"
)
