package audit

import (
	"context"
	"fmt"
)

// AuditService defines the interface for audit operations
type AuditService interface {
	// LogAction logs an audit entry
	LogAction(ctx context.Context, entry AuditEntry) error

	// GetHistory retrieves audit history based on query parameters
	GetHistory(ctx context.Context, query AuditQuery) (*AuditQueryResult, error)

	// GetByID retrieves an audit entry by its ID
	GetByID(ctx context.Context, id string) (*AuditEntry, error)

	// GetResourceHistory retrieves audit history for a specific resource
	GetResourceHistory(ctx context.Context, resourceType, resourceID string, limit int) ([]AuditEntry, error)

	// GetActorHistory retrieves audit history for a specific actor
	GetActorHistory(ctx context.Context, actorID string, actorType ActorType, limit int) ([]AuditEntry, error)

	// Close closes the service and underlying connections
	Close(ctx context.Context) error
}

// auditService implements the AuditService interface
type auditService struct {
	repo AuditRepository
}

// NewService creates a new audit service with the given configuration
func NewService(config *Config) (AuditService, error) {
	repo, err := NewMongoRepository(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}

	return &auditService{
		repo: repo,
	}, nil
}

// NewServiceWithRepository creates a new audit service with a custom repository
func NewServiceWithRepository(repo AuditRepository) AuditService {
	return &auditService{
		repo: repo,
	}
}

// LogAction logs an audit entry
func (s *auditService) LogAction(ctx context.Context, entry AuditEntry) error {
	if err := s.validateAuditEntry(entry); err != nil {
		return fmt.Errorf("invalid audit entry: %w", err)
	}

	return s.repo.Insert(ctx, entry)
}

// GetHistory retrieves audit history based on query parameters
func (s *auditService) GetHistory(ctx context.Context, query AuditQuery) (*AuditQueryResult, error) {
	if err := s.validateAuditQuery(query); err != nil {
		return nil, fmt.Errorf("invalid query: %w", err)
	}

	return s.repo.FindByQuery(ctx, query)
}

// GetByID retrieves an audit entry by its ID
func (s *auditService) GetByID(ctx context.Context, id string) (*AuditEntry, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	return s.repo.FindByID(ctx, id)
}

// GetResourceHistory retrieves audit history for a specific resource
func (s *auditService) GetResourceHistory(ctx context.Context, resourceType, resourceID string, limit int) ([]AuditEntry, error) {
	if resourceType == "" {
		return nil, fmt.Errorf("resource type cannot be empty")
	}
	if resourceID == "" {
		return nil, fmt.Errorf("resource ID cannot be empty")
	}
	if limit < 0 {
		return nil, fmt.Errorf("limit cannot be negative")
	}

	return s.repo.FindByResource(ctx, resourceType, resourceID, limit)
}

// GetActorHistory retrieves audit history for a specific actor
func (s *auditService) GetActorHistory(ctx context.Context, actorID string, actorType ActorType, limit int) ([]AuditEntry, error) {
	if actorID == "" {
		return nil, fmt.Errorf("actor ID cannot be empty")
	}
	if limit < 0 {
		return nil, fmt.Errorf("limit cannot be negative")
	}

	return s.repo.FindByActor(ctx, actorID, actorType, limit)
}

// Close closes the service and underlying connections
func (s *auditService) Close(ctx context.Context) error {
	return s.repo.Close(ctx)
}

// validateAuditEntry validates an audit entry
func (s *auditService) validateAuditEntry(entry AuditEntry) error {
	if entry.Action == "" {
		return fmt.Errorf("action cannot be empty")
	}
	if entry.Actor.ID == "" {
		return fmt.Errorf("actor ID cannot be empty")
	}
	if entry.Actor.Type == "" {
		return fmt.Errorf("actor type cannot be empty")
	}
	if entry.Resource.Type == "" {
		return fmt.Errorf("resource type cannot be empty")
	}
	if entry.Resource.ID == "" {
		return fmt.Errorf("resource ID cannot be empty")
	}

	// Validate actor type
	switch entry.Actor.Type {
	case ActorTypeUser, ActorTypeSystem, ActorTypeService, ActorTypeAPI, ActorTypeAdmin:
		// Valid actor types
	default:
		return fmt.Errorf("invalid actor type: %s", entry.Actor.Type)
	}

	// Validate action
	switch entry.Action {
	case ActionCreate, ActionUpdate, ActionDelete, ActionLogin, ActionLogout, ActionView, ActionExport:
		// Valid actions
	default:
		return fmt.Errorf("invalid action: %s", entry.Action)
	}

	return nil
}

// validateAuditQuery validates an audit query
func (s *auditService) validateAuditQuery(query AuditQuery) error {
	if query.Limit < 0 {
		return fmt.Errorf("limit cannot be negative")
	}
	if query.Offset < 0 {
		return fmt.Errorf("offset cannot be negative")
	}
	if query.StartTime != nil && query.EndTime != nil {
		if query.StartTime.After(*query.EndTime) {
			return fmt.Errorf("start time cannot be after end time")
		}
	}

	// Validate actor type if provided
	if query.ActorType != "" {
		switch query.ActorType {
		case ActorTypeUser, ActorTypeSystem, ActorTypeService, ActorTypeAPI, ActorTypeAdmin:
			// Valid actor types
		default:
			return fmt.Errorf("invalid actor type: %s", query.ActorType)
		}
	}

	// Validate actions if provided
	for _, action := range query.Actions {
		switch action {
		case ActionCreate, ActionUpdate, ActionDelete, ActionLogin, ActionLogout, ActionView, ActionExport:
			// Valid actions
		default:
			return fmt.Errorf("invalid action: %s", action)
		}
	}

	return nil
}
