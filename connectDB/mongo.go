package connectDB

import (
	mongo "github.com/TerrexTech/go-mongoutils/mongo"
	"github.com/bhupeshbhatia/go-agg-inven-mongo-cmd/model"
)

// CreateClient creates a MongoDB-Client.
func CreateClient() (*mongo.Client, error) {
	// Would ideally set these config-params as environment vars
	config := mongo.ClientConfig{
		Hosts:               []string{"localhost:27017"},
		Username:            "root",
		Password:            "root",
		TimeoutMilliseconds: 5000,
	}

	// ====> MongoDB Client
	client, err := mongo.NewClient(config)
	// Let the parent functions handle error, always -.-
	// (Even though in these examples, we won't, for simplicity)
	return client, err
}

// createCollection demonstrates creating the collection and the associated database.
func CreateCollection(client *mongo.Client, name string, database string) (*mongo.Collection, error) {
	// ====> Collection Configuration
	conn := &mongo.ConnectionConfig{
		Client:  client,
		Timeout: 5000,
	}
	// Index Configuration
	indexConfigs := []mongo.IndexConfig{
		mongo.IndexConfig{
			ColumnConfig: []mongo.IndexColumnConfig{
				mongo.IndexColumnConfig{
					Name:        "fruit_id",
					IsDescOrder: false,
				},
			},
			IsUnique: true,
			Name:     "fruit_id_index",
		},
	}

	// ====> Create New Collection
	c := &mongo.Collection{
		Connection:   conn,
		Name:         name,
		Database:     database,
		SchemaStruct: &model.Inventory{},
		Indexes:      indexConfigs,
	}
	return mongo.EnsureCollection(c)
}
