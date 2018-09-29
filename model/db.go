package model

import (
	"encoding/json"
	"log"

	mongo "github.com/TerrexTech/go-mongoutils/mongo"
	"github.com/pkg/errors"
)

type Datastore interface {
	Collection() *mongo.Collection
	CreateDataMongo()
	PopulateInvTable(body []byte) ([]byte, error)
}

type Db struct {
	collection *mongo.Collection
	*Db
}

type DbConfig struct {
	Hosts               []string
	Username            string
	Password            string
	TimeoutMilliseconds uint32
	Database            string
	Collection          string
}

//===============================================
type InvSearch struct {
	EndDate   int64  `bson:"end_date,omitempty" json:"end_date,omitempty"`
	StartDate int64  `bson:"start_date,omitempty" json:"start_date,omitempty"`
	SearchKey string `bson:"search_key,omitempty" json:"search_key,omitempty"`
	SearchVal string `bson:"search_val,omitempty" json:"search_val,omitempty"`
}

func ConfirmDbExists(dbConfig DbConfig) (*Db, error) {
	clientConfig := mongo.ClientConfig{
		Hosts:               dbConfig.Hosts,
		Username:            dbConfig.Username,
		Password:            dbConfig.Password,
		TimeoutMilliseconds: 3000,
	}

	// ====> MongoDB Client
	client, err := mongo.NewClient(clientConfig)
	if err != nil {
		err = errors.Wrap(err, "Mongo client not available")
		return nil, err
	}

	conn := &mongo.ConnectionConfig{
		Client:  client,
		Timeout: 5000,
	}

	// Index Configuration
	indexConfigs := []mongo.IndexConfig{
		mongo.IndexConfig{
			ColumnConfig: []mongo.IndexColumnConfig{
				mongo.IndexColumnConfig{
					Name:        "item_id",
					IsDescOrder: false,
				},
			},
			IsUnique: true,
			Name:     "item_id_index",
		},
	}

	// ====> Create New Collection
	colConfig := &mongo.Collection{
		Connection:   conn,
		Name:         dbConfig.Collection,
		Database:     dbConfig.Database,
		SchemaStruct: &Inventory{},
		Indexes:      indexConfigs,
	}
	c, err := mongo.EnsureCollection(colConfig)
	if err != nil {
		err = errors.Wrap(err, "Error creating DB-client")
		return nil, err
	}

	return &Db{
		collection: c,
	}, nil

}

// Collection returns the currrent MongoDB collection
func (d *Db) Collection() *mongo.Collection {
	return d.collection
}

func (db *Db) CreateDataMongo() ([]byte, error) {
	newInventory := []Inventory{}
	for i := 0; i < 100; i++ {
		newInventory = append(newInventory, GenerateDataForInv())
	}

	for _, v := range newInventory {
		insertResult, err := db.collection.InsertOne(v)
		if err != nil {
			err = errors.Wrap(err, "Unable to insert data")
			log.Println(err)
			return nil, err
		}
		log.Println(insertResult)
	}

	marInv, err := json.Marshal(&newInventory)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return marInv, nil
}

func (db *Db) PopulateInvTable(body []byte) ([]byte, error) {
	searchUpLimit := InvSearch{}

	//Read the body - so unmarshal
	err := json.Unmarshal(body, &searchUpLimit)
	if err != nil {
		err = errors.Wrap(err, "Unable to unmarshal - SearchOneTime")
		log.Println(err)
		return nil, err
	}

	//Find
	findResults, err := db.collection.Find(map[string]interface{}{
		"timestamp": map[string]int64{
			"$lte": searchUpLimit.EndDate,
		},
	})
	if err != nil {
		err = errors.Wrap(err, "Error while fetching product.")
		log.Println(err)
		return nil, err
	}

	inventory := []Inventory{}

	for _, v := range findResults {
		resultInv := v.(*Inventory)
		inventory = append(inventory, *resultInv)
	}

	//Marshal into []byte
	totalResult, err := json.Marshal(&inventory)
	if err != nil {
		err = errors.Wrap(err, "Unable to create response body")
		log.Println(err)
		return nil, err
	}
	return totalResult, nil
}

// func (db *Db) SearchByEndDate(endDate int64) *[]Inventory {

// 	findResults, err := db.collection.Find(map[string]interface{}{

// 		"timestamp": map[string]int64{
// 			"$lte": endDate,
// 		},
// 	})
// 	if err != nil {
// 		err = errors.Wrap(err, "Error while fetching product.")
// 		log.Println(err)
// 		return nil
// 	}

// 	inventory := []Inventory{}

// 	for _, v := range findResults {
// 		resultInv := v.(*Inventory)
// 		inventory = append(inventory, *resultInv)
// 	}

// 	return &inventory
// }
