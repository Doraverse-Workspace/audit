package audit

import (
	"time"
)

// Config represents the configuration for the audit module
type Config struct {
	// MongoDB connection settings
	MongoURI       string `json:"mongo_uri" yaml:"mongo_uri"`
	DatabaseName   string `json:"database_name" yaml:"database_name"`
	CollectionName string `json:"collection_name" yaml:"collection_name"`

	// Connection pool settings
	MaxPoolSize    uint64        `json:"max_pool_size" yaml:"max_pool_size"`
	MinPoolSize    uint64        `json:"min_pool_size" yaml:"min_pool_size"`
	ConnectTimeout time.Duration `json:"connect_timeout" yaml:"connect_timeout"`

	// Retry settings
	MaxRetries int           `json:"max_retries" yaml:"max_retries"`
	RetryDelay time.Duration `json:"retry_delay" yaml:"retry_delay"`

	// Performance settings
	BatchSize     int  `json:"batch_size" yaml:"batch_size"`
	EnableIndexes bool `json:"enable_indexes" yaml:"enable_indexes"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		MongoURI:       "mongodb://localhost:27017",
		DatabaseName:   "audit",
		CollectionName: "audit_logs",
		MaxPoolSize:    100,
		MinPoolSize:    5,
		ConnectTimeout: 10 * time.Second,
		MaxRetries:     3,
		RetryDelay:     time.Second,
		BatchSize:      1000,
		EnableIndexes:  true,
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.MongoURI == "" {
		return ErrInvalidConfig{Field: "MongoURI", Message: "cannot be empty"}
	}
	if c.DatabaseName == "" {
		return ErrInvalidConfig{Field: "DatabaseName", Message: "cannot be empty"}
	}
	if c.CollectionName == "" {
		return ErrInvalidConfig{Field: "CollectionName", Message: "cannot be empty"}
	}
	if c.MaxRetries < 0 {
		return ErrInvalidConfig{Field: "MaxRetries", Message: "cannot be negative"}
	}
	if c.RetryDelay < 0 {
		return ErrInvalidConfig{Field: "RetryDelay", Message: "cannot be negative"}
	}
	if c.BatchSize <= 0 {
		return ErrInvalidConfig{Field: "BatchSize", Message: "must be positive"}
	}
	return nil
}

// ErrInvalidConfig represents a configuration validation error
type ErrInvalidConfig struct {
	Field   string
	Message string
}

func (e ErrInvalidConfig) Error() string {
	return "invalid config field '" + e.Field + "': " + e.Message
}
