package mongo

import (
	"context"
	"fmt"

	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/writeconcern"
	"go.uber.org/fx"
)

type MongoClient struct {
	Client *mongo.Client
}

type MongoClientParam struct {
	fx.In
	fx.Lifecycle
	logger.Logger
}

type MongoClientConfig struct {
	URI        string `validate:"required,mongodb_connection_string"`
	ReplicaSet string
}

// NewMongoClient creates a new MongoDB client.
// TODO: Add OTel instrumentation once the library is compatible with the MongoDB driver v2. See <https://github.com/open-telemetry/opentelemetry-go-contrib/issues/6419>
func NewMongoClient(p MongoClientParam, cfg *MongoClientConfig) (*MongoClient, error) {
	opts := options.Client().
		ApplyURI(cfg.URI).
		SetLoggerOptions(
			options.Logger().
				SetComponentLevel(options.LogComponentAll, options.LogLevelInfo).
				SetSink(NewMongoClientLogger(p.Logger)),
		)

	if cfg.ReplicaSet != "" {
		opts = opts.SetReplicaSet(cfg.ReplicaSet)
	}

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return client.Ping(ctx, nil)
		},
		OnStop: func(ctx context.Context) error {
			return client.Disconnect(ctx)
		},
	})

	return &MongoClient{Client: client}, nil
}

func (c *MongoClient) Transaction(ctx context.Context, fc func(sessionCtx context.Context) (interface{}, error)) (interface{}, error) {
	session, err := c.Client.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	txnOpts := options.Transaction().SetWriteConcern(writeconcern.Majority()) // ensure data durability

	return session.WithTransaction(ctx, fc, txnOpts)
}
