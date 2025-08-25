package audit

import (
	"context"
)

// AuditRepository defines the interface for audit data storage
type AuditRepository interface {
	// Insert inserts a new audit entry
	Insert(ctx context.Context, entry AuditEntry) error

	// FindByQuery finds audit entries based on query parameters
	FindByQuery(ctx context.Context, query AuditQuery) (*AuditQueryResult, error)

	// FindByID finds an audit entry by its ID
	FindByID(ctx context.Context, id string) (*AuditEntry, error)

	// FindByResource finds audit entries for a specific resource
	FindByResource(ctx context.Context, resourceType, resourceID string, limit int) ([]AuditEntry, error)

	// FindByActor finds audit entries for a specific actor
	FindByActor(ctx context.Context, actorID string, actorType ActorType, limit int) ([]AuditEntry, error)

	// EnsureIndexes creates necessary database indexes for optimal performance
	EnsureIndexes(ctx context.Context) error

	// Close closes the repository connection
	Close(ctx context.Context) error
}
