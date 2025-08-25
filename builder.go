package audit

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuditBuilder provides a fluent interface for building and logging audit entries
type AuditBuilder struct {
	entry   AuditEntry
	service AuditService
}

// defaultService is the global audit service instance
var defaultService AuditService

// SetDefaultService sets the global audit service instance
func SetDefaultService(service AuditService) {
	defaultService = service
}

// NewAuditBuilder creates a new audit builder
func NewAuditBuilder() *AuditBuilder {
	return &AuditBuilder{
		entry: AuditEntry{
			ID:        primitive.NewObjectID(),
			Timestamp: time.Now().UTC(),
			Metadata:  make(map[string]any),
		},
		service: defaultService,
	}
}

// NewAuditBuilderWithService creates a new audit builder with a specific service
func NewAuditBuilderWithService(service AuditService) *AuditBuilder {
	return &AuditBuilder{
		entry: AuditEntry{
			ID:        primitive.NewObjectID(),
			Timestamp: time.Now().UTC(),
			Metadata:  make(map[string]any),
		},
		service: service,
	}
}

// Action sets the action type
func (b *AuditBuilder) Action(action AuditAction) *AuditBuilder {
	b.entry.Action = action
	return b
}

// Actor sets the actor information
func (b *AuditBuilder) Actor(actorID string, actorType ActorType, actorName string) *AuditBuilder {
	b.entry.Actor = Actor{
		ID:   actorID,
		Type: actorType,
		Name: actorName,
	}
	return b
}

// ActorWithSession sets the actor information including session ID
func (b *AuditBuilder) ActorWithSession(actorID string, actorType ActorType, actorName, sessionID string) *AuditBuilder {
	b.entry.Actor = Actor{
		ID:        actorID,
		Type:      actorType,
		Name:      actorName,
		SessionID: sessionID,
	}
	return b
}

// Session sets the session ID for the current actor
func (b *AuditBuilder) Session(sessionID string) *AuditBuilder {
	b.entry.Actor.SessionID = sessionID
	return b
}

// Resource sets the target resource information
func (b *AuditBuilder) Resource(resourceType, resourceID, resourceName string) *AuditBuilder {
	b.entry.Resource = AuditResource{
		Type: resourceType,
		ID:   resourceID,
		Name: resourceName,
	}
	return b
}

// AddChange adds a field change to the audit entry
func (b *AuditBuilder) AddChange(field string, oldValue, newValue any) *AuditBuilder {
	if b.entry.Changes == nil {
		b.entry.Changes = make([]FieldChange, 0)
	}
	b.entry.Changes = append(b.entry.Changes, FieldChange{
		Field:    field,
		OldValue: oldValue,
		NewValue: newValue,
	})
	return b
}

// Metadata adds metadata to the audit entry
func (b *AuditBuilder) Metadata(key string, value any) *AuditBuilder {
	if b.entry.Metadata == nil {
		b.entry.Metadata = make(map[string]any)
	}
	b.entry.Metadata[key] = value
	return b
}

// IPAddress sets the IP address
func (b *AuditBuilder) IPAddress(ip string) *AuditBuilder {
	b.entry.IPAddress = ip
	return b
}

// UserAgent sets the user agent
func (b *AuditBuilder) UserAgent(ua string) *AuditBuilder {
	b.entry.UserAgent = ua
	return b
}

// Success sets the success status
func (b *AuditBuilder) Success(success bool) *AuditBuilder {
	b.entry.Success = success
	return b
}

// Error sets the error message and marks the operation as failed
func (b *AuditBuilder) Error(err error) *AuditBuilder {
	if err != nil {
		b.entry.ErrorMsg = err.Error()
		b.entry.Success = false
	}
	return b
}

// Timestamp sets a custom timestamp
func (b *AuditBuilder) Timestamp(timestamp time.Time) *AuditBuilder {
	b.entry.Timestamp = timestamp
	return b
}

// Build returns the built audit entry without logging it
func (b *AuditBuilder) Build() AuditEntry {
	return b.entry
}

// Log logs the audit entry using the configured service
func (b *AuditBuilder) Log(ctx context.Context) error {
	if b.service == nil {
		return ErrNoServiceConfigured{}
	}
	return b.service.LogAction(ctx, b.entry)
}

// Convenience methods for common actor types

// User sets the actor as a user
func (b *AuditBuilder) User(userID, userName string) *AuditBuilder {
	return b.Actor(userID, ActorTypeUser, userName)
}

// UserWithSession sets the actor as a user with session
func (b *AuditBuilder) UserWithSession(userID, userName, sessionID string) *AuditBuilder {
	return b.ActorWithSession(userID, ActorTypeUser, userName, sessionID)
}

// System sets the actor as a system process
func (b *AuditBuilder) System(systemID, systemName string) *AuditBuilder {
	return b.Actor(systemID, ActorTypeSystem, systemName)
}

// Service sets the actor as a service account
func (b *AuditBuilder) Service(serviceID, serviceName string) *AuditBuilder {
	return b.Actor(serviceID, ActorTypeService, serviceName)
}

// API sets the actor as an API client
func (b *AuditBuilder) API(apiID, apiName string) *AuditBuilder {
	return b.Actor(apiID, ActorTypeAPI, apiName)
}

// Admin sets the actor as an admin user
func (b *AuditBuilder) Admin(adminID, adminName string) *AuditBuilder {
	return b.Actor(adminID, ActorTypeAdmin, adminName)
}

// AdminWithSession sets the actor as an admin user with session
func (b *AuditBuilder) AdminWithSession(adminID, adminName, sessionID string) *AuditBuilder {
	return b.ActorWithSession(adminID, ActorTypeAdmin, adminName, sessionID)
}

// Convenience methods for common actions

// Create sets the action to create
func (b *AuditBuilder) Create() *AuditBuilder {
	return b.Action(ActionCreate)
}

// Update sets the action to update
func (b *AuditBuilder) Update() *AuditBuilder {
	return b.Action(ActionUpdate)
}

// Delete sets the action to delete
func (b *AuditBuilder) Delete() *AuditBuilder {
	return b.Action(ActionDelete)
}

// Login sets the action to login
func (b *AuditBuilder) Login() *AuditBuilder {
	return b.Action(ActionLogin)
}

// Logout sets the action to logout
func (b *AuditBuilder) Logout() *AuditBuilder {
	return b.Action(ActionLogout)
}

// View sets the action to view
func (b *AuditBuilder) View() *AuditBuilder {
	return b.Action(ActionView)
}

// Export sets the action to export
func (b *AuditBuilder) Export() *AuditBuilder {
	return b.Action(ActionExport)
}

// ErrNoServiceConfigured represents an error when no audit service is configured
type ErrNoServiceConfigured struct{}

func (e ErrNoServiceConfigured) Error() string {
	return "no audit service configured: use SetDefaultService() or NewAuditBuilderWithService()"
}
