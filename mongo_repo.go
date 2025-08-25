package audit

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongoRepository implements the AuditRepository interface using MongoDB
type mongoRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
	config     *Config
}

// NewMongoRepository creates a new MongoDB repository
func NewMongoRepository(config *Config) (AuditRepository, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	clientOptions := options.Client().
		ApplyURI(config.MongoURI).
		SetMaxPoolSize(config.MaxPoolSize).
		SetMinPoolSize(config.MinPoolSize).
		SetConnectTimeout(config.ConnectTimeout)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), config.ConnectTimeout)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	collection := client.Database(config.DatabaseName).Collection(config.CollectionName)

	repo := &mongoRepository{
		client:     client,
		collection: collection,
		config:     config,
	}

	// Create indexes if enabled
	if config.EnableIndexes {
		if err := repo.EnsureIndexes(context.Background()); err != nil {
			return nil, fmt.Errorf("failed to create indexes: %w", err)
		}
	}

	return repo, nil
}

// Insert inserts a new audit entry
func (r *mongoRepository) Insert(ctx context.Context, entry AuditEntry) error {
	// Set timestamp if not provided
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}

	// Generate ID if not provided
	if entry.ID.IsZero() {
		entry.ID = primitive.NewObjectID()
	}

	var err error
	for attempt := 0; attempt <= r.config.MaxRetries; attempt++ {
		_, err = r.collection.InsertOne(ctx, entry)
		if err == nil {
			return nil
		}

		if attempt < r.config.MaxRetries {
			time.Sleep(r.config.RetryDelay)
		}
	}

	return fmt.Errorf("failed to insert audit entry after %d attempts: %w", r.config.MaxRetries+1, err)
}

// FindByQuery finds audit entries based on query parameters
func (r *mongoRepository) FindByQuery(ctx context.Context, query AuditQuery) (*AuditQueryResult, error) {
	filter := r.buildFilter(query)

	// Get total count
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to count documents: %w", err)
	}

	// Build options
	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}) // Sort by timestamp descending

	if query.Limit > 0 {
		opts.SetLimit(int64(query.Limit))
	}
	if query.Offset > 0 {
		opts.SetSkip(int64(query.Offset))
	}

	// Execute query
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer cursor.Close(ctx)

	var entries []AuditEntry
	if err := cursor.All(ctx, &entries); err != nil {
		return nil, fmt.Errorf("failed to decode results: %w", err)
	}

	hasMore := false
	if query.Limit > 0 && int64(query.Offset+len(entries)) < total {
		hasMore = true
	}

	return &AuditQueryResult{
		Entries: entries,
		Total:   total,
		HasMore: hasMore,
	}, nil
}

// FindByID finds an audit entry by its ID
func (r *mongoRepository) FindByID(ctx context.Context, id string) (*AuditEntry, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	filter := bson.M{"_id": objectID}

	var entry AuditEntry
	err = r.collection.FindOne(ctx, filter).Decode(&entry)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find audit entry: %w", err)
	}

	return &entry, nil
}

// FindByResource finds audit entries for a specific resource
func (r *mongoRepository) FindByResource(ctx context.Context, resourceType, resourceID string, limit int) ([]AuditEntry, error) {
	filter := bson.M{
		"resource.type": resourceType,
		"resource.id":   resourceID,
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}})

	if limit > 0 {
		opts.SetLimit(int64(limit))
	}

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find resource history: %w", err)
	}
	defer cursor.Close(ctx)

	var entries []AuditEntry
	if err := cursor.All(ctx, &entries); err != nil {
		return nil, fmt.Errorf("failed to decode results: %w", err)
	}

	return entries, nil
}

// FindByActor finds audit entries for a specific actor
func (r *mongoRepository) FindByActor(ctx context.Context, actorID string, actorType ActorType, limit int) ([]AuditEntry, error) {
	filter := bson.M{
		"actor.id": actorID,
	}

	if actorType != "" {
		filter["actor.type"] = actorType
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}})

	if limit > 0 {
		opts.SetLimit(int64(limit))
	}

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find actor history: %w", err)
	}
	defer cursor.Close(ctx)

	var entries []AuditEntry
	if err := cursor.All(ctx, &entries); err != nil {
		return nil, fmt.Errorf("failed to decode results: %w", err)
	}

	return entries, nil
}

// EnsureIndexes creates necessary database indexes
func (r *mongoRepository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "actor.id", Value: 1}, {Key: "timestamp", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "actor.type", Value: 1}, {Key: "timestamp", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "actor.session_id", Value: 1}, {Key: "timestamp", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "resource.type", Value: 1}, {Key: "resource.id", Value: 1}, {Key: "timestamp", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "action", Value: 1}, {Key: "timestamp", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "timestamp", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "actor.id", Value: 1}, {Key: "actor.type", Value: 1}, {Key: "timestamp", Value: -1}},
		},
	}

	opts := options.CreateIndexes().SetMaxTime(30 * time.Second)
	_, err := r.collection.Indexes().CreateMany(ctx, indexes, opts)
	if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	return nil
}

// Close closes the repository connection
func (r *mongoRepository) Close(ctx context.Context) error {
	return r.client.Disconnect(ctx)
}

// buildFilter builds a MongoDB filter from AuditQuery
func (r *mongoRepository) buildFilter(query AuditQuery) bson.M {
	filter := bson.M{}

	if query.ActorID != "" {
		filter["actor.id"] = query.ActorID
	}
	if query.ActorType != "" {
		filter["actor.type"] = query.ActorType
	}
	if query.SessionID != "" {
		filter["actor.session_id"] = query.SessionID
	}
	if len(query.Actions) > 0 {
		filter["action"] = bson.M{"$in": query.Actions}
	}
	if query.ResourceType != "" {
		filter["resource.type"] = query.ResourceType
	}
	if query.ResourceID != "" {
		filter["resource.id"] = query.ResourceID
	}
	if query.Success != nil {
		filter["success"] = *query.Success
	}

	// Time range filter
	if query.StartTime != nil || query.EndTime != nil {
		timeFilter := bson.M{}
		if query.StartTime != nil {
			timeFilter["$gte"] = *query.StartTime
		}
		if query.EndTime != nil {
			timeFilter["$lte"] = *query.EndTime
		}
		filter["timestamp"] = timeFilter
	}

	return filter
}
