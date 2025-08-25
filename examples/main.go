package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Doraverse-Workspace/audit"
)

func main() {
	Example()

	// ExampleAdvancedConfiguration()
}

func Example() {
	// Initialize with default configuration
	err := audit.InitializeWithDefaults("mongodb://localhost:27017", "myapp")
	if err != nil {
		log.Fatal(err)
	}
	defer audit.Shutdown(context.Background())

	ctx := context.Background()

	// Example 1: Log user login
	err = audit.NewAuditBuilder().
		Login().
		User("user123", "John Doe").
		Session("session456").
		Resource("session", "session456", "User Session").
		IPAddress("192.168.1.1").
		UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36").
		Success(true).
		Log(ctx)
	if err != nil {
		log.Printf("Failed to log user login: %v", err)
	}

	// Example 2: Log data update with changes
	err = audit.NewAuditBuilder().
		Update().
		User("user123", "John Doe").
		Session("session456").
		Resource("user", "user456", "Jane Smith Profile").
		AddChange("email", "old@example.com", "new@example.com").
		AddChange("status", "inactive", "active").
		Metadata("source", "admin_panel").
		Metadata("reason", "account_activation").
		Success(true).
		Log(ctx)
	if err != nil {
		log.Printf("Failed to log user update: %v", err)
	}

	// Example 3: Log system action
	err = audit.NewAuditBuilder().
		Update().
		System("cleanup_service", "Data Cleanup Service").
		Resource("user", "user789", "Inactive User").
		Metadata("reason", "automated_cleanup").
		Metadata("retention_policy", "90_days").
		Success(true).
		Log(ctx)
	if err != nil {
		log.Printf("Failed to log system action: %v", err)
	}

	// Example 4: Log failed operation
	err = audit.NewAuditBuilder().
		Delete().
		User("user123", "John Doe").
		Session("session456").
		Resource("document", "doc789", "Important Document").
		Error(fmt.Errorf("insufficient permissions")).
		Log(ctx)
	if err != nil {
		log.Printf("Failed to log failed operation: %v", err)
	}

	// Example 5: Query audit history
	startTime := time.Now().Add(-24 * time.Hour)
	endTime := time.Now()

	query := audit.AuditQuery{
		ActorID:   "user123",
		ActorType: audit.ActorTypeUser,
		Actions:   []audit.AuditAction{audit.ActionLogin, audit.ActionUpdate, audit.ActionDelete},
		StartTime: &startTime,
		EndTime:   &endTime,
		Limit:     50,
	}

	result, err := audit.GetHistory(ctx, query)
	if err != nil {
		log.Printf("Failed to get audit history: %v", err)
		return
	}

	fmt.Printf("Found %d audit entries (total: %d, has more: %v)\n",
		len(result.Entries), result.Total, result.HasMore)

	// Example 6: Get resource history
	resourceHistory, err := audit.GetResourceHistory(ctx, "user", "user456", 10)
	if err != nil {
		log.Printf("Failed to get resource history: %v", err)
		return
	}

	fmt.Printf("Found %d entries for resource user:user456\n", len(resourceHistory))

	// Example 7: Get actor history
	actorHistory, err := audit.GetActorHistory(ctx, "user123", audit.ActorTypeUser, 20)
	if err != nil {
		log.Printf("Failed to get actor history: %v", err)
		return
	}

	fmt.Printf("Found %d entries for actor user123\n", len(actorHistory))
}

func ExampleAdvancedConfiguration() {
	// Custom configuration
	config := audit.DefaultConfig()
	config.MongoURI = "mongodb://localhost:27017"
	config.DatabaseName = "myapp_audit"
	config.CollectionName = "audit_trail"
	config.MaxRetries = 5
	config.RetryDelay = 2 * time.Second

	// Create service with custom config
	service, err := audit.NewService(config)
	if err != nil {
		log.Fatal(err)
	}
	defer service.Close(context.Background())

	ctx := context.Background()

	// Use builder with specific service
	err = audit.NewAuditBuilderWithService(service).
		Action(audit.ActionExport).
		Actor("api_client_v2", audit.ActorTypeAPI, "Mobile App v2.0").
		Resource("report", "monthly_sales_2024", "Monthly Sales Report").
		Metadata("format", "pdf").
		Metadata("filters", map[string]any{
			"date_range": "2024-01-01 to 2024-01-31",
			"department": "sales",
		}).
		IPAddress("203.0.113.42").
		Success(true).
		Log(ctx)

	if err != nil {
		log.Printf("Failed to log export action: %v", err)
	}

	// Complex query with multiple filters
	query := audit.AuditQuery{
		ActorType:    audit.ActorTypeAPI,
		ResourceType: "report",
		Actions:      []audit.AuditAction{audit.ActionExport, audit.ActionView},
		Success:      boolPtr(true),
		Limit:        100,
		Offset:       0,
	}

	result, err := service.GetHistory(ctx, query)
	if err != nil {
		log.Printf("Failed to get filtered history: %v", err)
		return
	}

	for _, entry := range result.Entries {
		fmt.Printf("Export by %s (%s) at %s - %s:%s\n",
			entry.Actor.Name, entry.Actor.Type,
			entry.Timestamp.Format(time.RFC3339),
			entry.Resource.Type, entry.Resource.ID)
	}
}

func boolPtr(b bool) *bool {
	return &b
}
