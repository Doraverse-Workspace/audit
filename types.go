package audit

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuditEntry represents a single audit log entry
type AuditEntry struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
	Action    AuditAction        `bson:"action" json:"action"`
	Actor     Actor              `bson:"actor" json:"actor"`
	Resource  AuditResource      `bson:"resource" json:"resource"`
	Changes   []FieldChange      `bson:"changes,omitempty" json:"changes,omitempty"`
	Metadata  map[string]any     `bson:"metadata,omitempty" json:"metadata,omitempty"`
	IPAddress string             `bson:"ip_address,omitempty" json:"ip_address,omitempty"`
	UserAgent string             `bson:"user_agent,omitempty" json:"user_agent,omitempty"`
	Success   bool               `bson:"success" json:"success"`
	ErrorMsg  string             `bson:"error_msg,omitempty" json:"error_msg,omitempty"`
}

// Actor represents who/what performed the action
type Actor struct {
	ID        string    `bson:"id" json:"id"`                                     // actor identifier
	Type      ActorType `bson:"type" json:"type"`                                 // type of actor
	Name      string    `bson:"name,omitempty" json:"name,omitempty"`             // human-readable name
	SessionID string    `bson:"session_id,omitempty" json:"session_id,omitempty"` // session if applicable
}

// ActorType defines the type of actor performing an action
type ActorType string

const (
	ActorTypeUser    ActorType = "user"    // human user
	ActorTypeSystem  ActorType = "system"  // system/automated process
	ActorTypeService ActorType = "service" // service account
	ActorTypeAPI     ActorType = "api"     // API client/application
	ActorTypeAdmin   ActorType = "admin"   // admin user
)

// AuditAction defines the type of action being audited
type AuditAction string

const (
	ActionCreate AuditAction = "create"
	ActionUpdate AuditAction = "update"
	ActionDelete AuditAction = "delete"
	ActionLogin  AuditAction = "login"
	ActionLogout AuditAction = "logout"
	ActionView   AuditAction = "view"
	ActionExport AuditAction = "export"
)

// AuditResource represents the target resource of the action
type AuditResource struct {
	Type string `bson:"type" json:"type"`                     // e.g., "user", "session", "product"
	ID   string `bson:"id" json:"id"`                         // resource identifier
	Name string `bson:"name,omitempty" json:"name,omitempty"` // human-readable name
}

// FieldChange represents a change to a specific field
type FieldChange struct {
	Field    string `bson:"field" json:"field"`
	OldValue any    `bson:"old_value,omitempty" json:"old_value,omitempty"`
	NewValue any    `bson:"new_value,omitempty" json:"new_value,omitempty"`
}

// AuditQuery represents query parameters for searching audit logs
type AuditQuery struct {
	ActorID      string        `json:"actor_id,omitempty"`
	ActorType    ActorType     `json:"actor_type,omitempty"`
	SessionID    string        `json:"session_id,omitempty"`
	Actions      []AuditAction `json:"actions,omitempty"`
	ResourceType string        `json:"resource_type,omitempty"`
	ResourceID   string        `json:"resource_id,omitempty"`
	StartTime    *time.Time    `json:"start_time,omitempty"`
	EndTime      *time.Time    `json:"end_time,omitempty"`
	Success      *bool         `json:"success,omitempty"`
	Limit        int           `json:"limit,omitempty"`
	Offset       int           `json:"offset,omitempty"`
}

// AuditQueryResult represents the result of an audit query
type AuditQueryResult struct {
	Entries []AuditEntry `json:"entries"`
	Total   int64        `json:"total"`
	HasMore bool         `json:"has_more"`
}
