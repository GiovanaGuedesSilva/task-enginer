package common

// Project types

type ProjectStatus string

const (
	ProjectStatusActive   ProjectStatus = "active"
	ProjectStatusInactive ProjectStatus = "inactive"
	ProjectStatusArchived ProjectStatus = "archived"
	ProjectStatusDeleted  ProjectStatus = "deleted"
)

type ProjectPriority string

const (
	ProjectPriorityLow    ProjectPriority = "low"
	ProjectPriorityMedium ProjectPriority = "medium"
	ProjectPriorityHigh   ProjectPriority = "high"
	ProjectPriorityUrgent ProjectPriority = "urgent"
)

// Task types

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusCancelled  TaskStatus = "cancelled"
)

type TaskPriority string

const (
	TaskPriorityLow    TaskPriority = "low"
	TaskPriorityMedium TaskPriority = "medium"
	TaskPriorityHigh   TaskPriority = "high"
	TaskPriorityUrgent TaskPriority = "urgent"
)

// User types

type UserRole string

const (
	UserRoleAdmin    UserRole = "admin"
	UserRoleManager  UserRole = "manager"
	UserRoleMember   UserRole = "member"
	UserRoleObserver UserRole = "observer"
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
)
