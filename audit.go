// Package audit provides a comprehensive audit logging system for Go applications.
// It supports MongoDB as the storage backend and provides a fluent API for logging
// various types of actions performed by different types of actors on resources.
//
// Installation:
//
//	go get github.com/Doraverse-Workspace/audit
//
// Basic usage:
//
//	import "github.com/Doraverse-Workspace/audit"
//
//	config := audit.DefaultConfig()
//	config.MongoURI = "mongodb://localhost:27017"
//	config.DatabaseName = "myapp"
//
//	service, err := audit.NewService(config)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer service.Close(context.Background())
//
//	// Set as default service for builders
//	audit.SetDefaultService(service)
//
//	// Log user login
//	err = audit.NewAuditBuilder().
//		Login().
//		User("user123", "John Doe").
//		Session("session456").
//		Resource("session", "session456", "User Session").
//		IPAddress("192.168.1.1").
//		Success(true).
//		Log(context.Background())
//
//	// Query audit history
//	query := audit.AuditQuery{
//		ActorID: "user123",
//		Actions: []audit.AuditAction{audit.ActionLogin, audit.ActionLogout},
//		Limit:   50,
//	}
//	result, err := service.GetHistory(context.Background(), query)
package audit

import (
	"context"
	"fmt"
)

// Version information
const (
	Version = "1.0.0"
)

// Package-level convenience functions

// Initialize initializes the audit module with the given configuration
// and sets it as the default service for builders.
func Initialize(config *Config) error {
	service, err := NewService(config)
	if err != nil {
		return fmt.Errorf("failed to initialize audit module: %w", err)
	}

	SetDefaultService(service)
	return nil
}

// InitializeWithDefaults initializes the audit module with default configuration
// and the specified MongoDB URI and database name.
func InitializeWithDefaults(mongoURI, databaseName string) error {
	config := DefaultConfig()
	config.MongoURI = mongoURI
	config.DatabaseName = databaseName

	return Initialize(config)
}

// LogAction is a convenience function to log an audit entry using the default service
func LogAction(ctx context.Context, entry AuditEntry) error {
	if defaultService == nil {
		return ErrNoServiceConfigured{}
	}
	return defaultService.LogAction(ctx, entry)
}

// GetHistory is a convenience function to get audit history using the default service
func GetHistory(ctx context.Context, query AuditQuery) (*AuditQueryResult, error) {
	if defaultService == nil {
		return nil, ErrNoServiceConfigured{}
	}
	return defaultService.GetHistory(ctx, query)
}

// GetByID is a convenience function to get an audit entry by ID using the default service
func GetByID(ctx context.Context, id string) (*AuditEntry, error) {
	if defaultService == nil {
		return nil, ErrNoServiceConfigured{}
	}
	return defaultService.GetByID(ctx, id)
}

// GetResourceHistory is a convenience function to get resource history using the default service
func GetResourceHistory(ctx context.Context, resourceType, resourceID string, limit int) ([]AuditEntry, error) {
	if defaultService == nil {
		return nil, ErrNoServiceConfigured{}
	}
	return defaultService.GetResourceHistory(ctx, resourceType, resourceID, limit)
}

// GetActorHistory is a convenience function to get actor history using the default service
func GetActorHistory(ctx context.Context, actorID string, actorType ActorType, limit int) ([]AuditEntry, error) {
	if defaultService == nil {
		return nil, ErrNoServiceConfigured{}
	}
	return defaultService.GetActorHistory(ctx, actorID, actorType, limit)
}

// Shutdown gracefully shuts down the default audit service
func Shutdown(ctx context.Context) error {
	if defaultService == nil {
		return nil
	}
	return defaultService.Close(ctx)
}
