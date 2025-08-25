# Audit Module

A comprehensive audit logging system for Go applications with MongoDB storage backend.

## Features

- **Comprehensive Audit Logging**: Track create, update, delete, login, logout, view, and export actions
- **Actor-based Architecture**: Support for different actor types (user, system, service, API, admin)
- **Resource Tracking**: Track actions on any type of resource with detailed metadata
- **Change Tracking**: Record field-level changes for update operations
- **Flexible Querying**: Rich query interface with filtering by actor, resource, time range, and more
- **MongoDB Integration**: Optimized for MongoDB with automatic indexing
- **Fluent API**: Easy-to-use builder pattern for logging audit entries
- **Session Tracking**: Link actions to user sessions for better traceability

## Installation

```bash
go get github.com/Doraverse-Workspace/audit
```

## Quick Start

```go
package main

import (
    "context"
    "log"
    
    "github.com/Doraverse-Workspace/audit"
)

func main() {
    // Initialize with MongoDB
    err := audit.InitializeWithDefaults("mongodb://localhost:27017", "myapp")
    if err != nil {
        log.Fatal(err)
    }
    defer audit.Shutdown(context.Background())

    ctx := context.Background()

    // Log user login
    err = audit.NewAuditBuilder().
        Login().
        User("user123", "John Doe").
        Session("session456").
        Resource("session", "session456", "User Session").
        IPAddress("192.168.1.1").
        Success(true).
        Log(ctx)
    if err != nil {
        log.Printf("Failed to log audit: %v", err)
    }

    // Query audit history
    query := audit.AuditQuery{
        ActorID: "user123",
        Limit:   10,
    }
    
    result, err := audit.GetHistory(ctx, query)
    if err != nil {
        log.Printf("Failed to get history: %v", err)
        return
    }
    
    log.Printf("Found %d audit entries", len(result.Entries))
}
```

## Core Concepts

### Actor Types

- `ActorTypeUser`: Human users
- `ActorTypeSystem`: Automated system processes  
- `ActorTypeService`: Service accounts
- `ActorTypeAPI`: API clients/applications
- `ActorTypeAdmin`: Administrative users

### Actions

- `ActionCreate`: Resource creation
- `ActionUpdate`: Resource modification
- `ActionDelete`: Resource deletion
- `ActionLogin`: User authentication
- `ActionLogout`: User session termination
- `ActionView`: Resource access/viewing
- `ActionExport`: Data export operations

## Usage Examples

### Basic Logging

```go
// User action with session
audit.NewAuditBuilder().
    Update().
    User("user123", "John Doe").
    Session("session456").
    Resource("user", "user456", "Jane Smith").
    AddChange("email", "old@example.com", "new@example.com").
    AddChange("status", "inactive", "active").
    Success(true).
    Log(ctx)

// System action
audit.NewAuditBuilder().
    Delete().
    System("cleanup_service", "Data Cleanup Service").
    Resource("file", "temp123", "Temporary File").
    Metadata("reason", "retention_policy").
    Success(true).
    Log(ctx)

// Failed operation
audit.NewAuditBuilder().
    Create().
    User("user123", "John Doe").
    Resource("document", "doc789", "Important Document").
    Error(fmt.Errorf("insufficient permissions")).
    Log(ctx)
```

### Querying History

```go
// Query by actor
query := audit.AuditQuery{
    ActorID:   "user123",
    ActorType: audit.ActorTypeUser,
    Actions:   []audit.AuditAction{audit.ActionLogin, audit.ActionLogout},
    Limit:     50,
}
result, err := audit.GetHistory(ctx, query)

// Query by resource
resourceHistory, err := audit.GetResourceHistory(ctx, "user", "user456", 10)

// Query by time range
startTime := time.Now().Add(-24 * time.Hour)
endTime := time.Now()
query = audit.AuditQuery{
    StartTime: &startTime,
    EndTime:   &endTime,
    Success:   &[]bool{true}[0], // Only successful operations
}
```

### Advanced Configuration

```go
config := audit.DefaultConfig()
config.MongoURI = "mongodb://localhost:27017"
config.DatabaseName = "myapp_audit"
config.CollectionName = "audit_trail"
config.MaxRetries = 5
config.RetryDelay = 2 * time.Second
config.MaxPoolSize = 50

service, err := audit.NewService(config)
if err != nil {
    log.Fatal(err)
}

// Use specific service instance
audit.NewAuditBuilderWithService(service).
    Export().
    API("mobile_app", "Mobile App v1.0").
    Resource("report", "sales_2024", "Sales Report").
    Log(ctx)
```

## Data Structure

### AuditEntry

```go
type AuditEntry struct {
    ID        primitive.ObjectID `json:"id"`
    Timestamp time.Time          `json:"timestamp"`
    Action    AuditAction        `json:"action"`
    Actor     Actor              `json:"actor"`
    Resource  AuditResource      `json:"resource"`
    Changes   []FieldChange      `json:"changes,omitempty"`
    Metadata  map[string]any     `json:"metadata,omitempty"`
    IPAddress string             `json:"ip_address,omitempty"`
    UserAgent string             `json:"user_agent,omitempty"`
    Success   bool               `json:"success"`
    ErrorMsg  string             `json:"error_msg,omitempty"`
}
```

### Actor

```go
type Actor struct {
    ID        string    `json:"id"`
    Type      ActorType `json:"type"`
    Name      string    `json:"name,omitempty"`
    SessionID string    `json:"session_id,omitempty"`
}
```

## MongoDB Indexes

The module automatically creates the following indexes for optimal query performance:

- `actor.id + timestamp` (descending)
- `actor.type + timestamp` (descending)  
- `actor.session_id + timestamp` (descending)
- `resource.type + resource.id + timestamp` (descending)
- `action + timestamp` (descending)
- `timestamp` (descending)
- `actor.id + actor.type + timestamp` (descending)

## Error Handling

The module provides detailed error messages and proper error wrapping. Common errors include:

- Configuration validation errors
- MongoDB connection errors  
- Invalid audit entry validation errors
- Query parameter validation errors

## Performance Considerations

- Uses connection pooling for MongoDB
- Implements retry logic with exponential backoff
- Optimized indexes for common query patterns
- Batch operations support for high-volume scenarios
- Configurable timeouts and limits

## License

[Add your license information here]
