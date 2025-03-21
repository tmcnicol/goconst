const eventTypes = [
	// User has been added to the system
	// Indicates that somthing foo
	"USER_INVITED",
	// New task has been created on a project
	"TASK_CREATED",
	// Update to an existing task
	"TASK_UPDATED",
	// New file has been uploaded to the system
	"FILE_UPLOADED",
] as const;

type EventType = typeof eventTypes[number];
