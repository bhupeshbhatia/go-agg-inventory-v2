package service

import (
	"encoding/json"
	"fmt"
	"log"

	mongo "github.com/TerrexTech/go-mongoutils/mongo"
	"github.com/bhupeshbhatia/go-agg-inven-mongo-cmd/model"
	mgo "github.com/mongodb/mongo-go-driver/mongo"
	"github.com/pkg/errors"
)

// //Global variable = aggregate version and ID
// var AggregateVersion = int64(1)
// var AggregateID = int64(1)

type InventoryData struct {
	Product     *model.Inventory
	MongoTable  *mongo.Collection
	SearchField string
	GetValue    interface{}
	FilterName  string
	FilterValue int64
}

func GetInventoryJSON(jsonString []byte) (*model.Inventory, error) {
	var inventory model.Inventory
	err := json.Unmarshal(jsonString, &inventory)
	if err != nil {
		err = errors.Wrap(err, "Unable to unmarshal foodItem into Inventory struct")
		log.Println(err)
		return nil, err
	}
	return &inventory, nil
}

func GetMarshal(inventory *model.Inventory) ([]byte, error) {
	inv, err := json.Marshal(inventory)
	if err != nil {
		err = errors.Wrap(err, "Unable to marshal foodItem into Inventory struct")
		log.Println(err)
		return nil, err
	}
	return inv, nil
}

func AddFood(data InventoryData) (*mgo.InsertOneResult, error) {

	if data.FilterName == "Fruit_ID" && data.Product.FruitID == 0 {
		return nil, errors.New("Error inserting record. No Fruit_ID found")
	}

	insertResult, err := data.MongoTable.InsertOne(data.Product)
	if err != nil {
		err = errors.Wrap(err, "Unable to insert event")
		log.Println(err)
		return nil, err
	}

	fmt.Println(insertResult)
	return insertResult, nil
	// return nil, nil
}

func GetFoodProducts(data InventoryData) ([]interface{}, error) {
	findResults, err := data.MongoTable.FindMap(map[string]interface{}{
		data.SearchField: map[string]interface{}{
			"$eq": data.GetValue,
		},
	})
	if err != nil {
		err = errors.Wrap(err, "Error while fetching food product.")
		log.Println(err)
		return nil, err
	}

	return findResults, nil
}

func UpdateAgg(data InventoryData) (*mgo.UpdateResult, error) {
	filter := &model.Inventory{
		FruitID: data.FilterValue,
	}
	if data.FilterValue == 0 {
		return nil, errors.New("Fruit_ID cannot be 0")
	}

	//how to convert struct to map?

	update := &map[string]interface{}{
		"fruit_id":        data.Product.FruitID,
		"origin":          data.Product.Origin,
		"device_id":       data.Product.DeviceID,
		"sale_price":      data.Product.SalePrice,
		"original_weight": data.Product.OriginalWeight,
	}

	updateResult, err := data.MongoTable.UpdateMany(filter, update)
	if err != nil {
		err = errors.Wrap(err, "Unable to update event")
		log.Println(err)
		return nil, err
	}
	fmt.Println(updateResult)

	return updateResult, nil
}

func DeleteAgg(data InventoryData) (*mgo.DeleteResult, error) {

	if data.FilterName == "Fruit_ID" && data.FilterValue == 0 {
		return nil, errors.New("Error deleting product.")
	}

	if data.FilterName != "Fruit_ID" {
		return nil, errors.New("Error deleting product.")
	}

	deleteResult, err := data.MongoTable.DeleteMany(&model.Inventory{
		FruitID: data.Product.FruitID,
	})
	if err != nil {
		err = errors.Wrap(err, "Unable to delete event")
		log.Println(err)
		return nil, err
	}

	fmt.Println(deleteResult)

	return deleteResult, nil
}

// func findProductById(inventory *model.Inventory, mgTable *mongo.Collection) ([]interface{}, error) {

// 	findResult, err := mgTable.Find(&model.Inventory{
// 		FruitID: inventory.FruitID,
// 	})
// 	if err != nil {
// 		err = errors.Wrap(err, "Unable to find product")
// 		log.Println(err)
// 		return nil, err
// 	}
// 	return findResult, nil
// }

// func findProductByName(inventory *model.Inventory, mgTable *mongo.Collection) ([]interface{}, error) {

// 	findResult, err := mgTable.Find(&model.Inventory{
// 		Name: inventory.Name,
// 	})
// 	if err != nil {
// 		err = errors.Wrap(err, "Unable to find product")
// 		log.Println(err)
// 		return nil, err
// 	}
// 	return findResult, nil
// }

// type TimeRange struct {
// 	Week  int64 `json:"week,omitempty"`
// 	Month int64 `json:"month,omitempty"`
// 	Year  int64 `json:"year,omitempty"`
// }

// func GetTimeJSON(jsonString []byte) (*TimeRange, error) {
// 	var timeRange TimeRange
// 	err := json.Unmarshal(jsonString, &timeRange)
// 	if err != nil {
// 		err = errors.Wrap(err, "Unable to unmarshal foodItem into Inventory struct")
// 		log.Println(err)
// 		return nil, err
// 	}
// 	return &timeRange, nil
// }

// //Of What?
// func FindByDateArrived(inventory *model.Inventory, specificPeriod string, mgTable *mongo.Collection) ([]interface{}, error) {
// 	var result []interface{}
// 	var greaterThanDate time.Time
// 	var lessThanDate time.Time

// 	greaterThanDate = time.Now().AddDate()

// 	switch specificPeriod {
// 	case "week":
// 		greaterThanDate = time.Now().AddDate(0, 0, -7)
// 		lessThanDate = time.Now()
// 	case "twoWeeks":
// 		greaterThanDate = time.Now().AddDate(0, 0, -14)
// 		lessThanDate = time.Now()
// 	case "oneMonth":
// 		greaterThanDate = time.Now().AddDate(0, -1, 0)
// 		lessThanDate = time.Now()
// 	case "threeMonths":
// 		greaterThanDate = time.Now().AddDate(0, -3, 0)
// 		lessThanDate = time.Now()
// 	case "sixMonths":
// 		greaterThanDate = time.Now().AddDate(0, -6, 0)
// 		lessThanDate = time.Now()
// 	case "oneYear":
// 		greaterThanDate = time.Now().AddDate(-1, 0, 0)
// 		lessThanDate = time.Now()
// 	}

// 	result, err := mgTable.FindMap(map[string]interface{}{
// 		"date_arrived": map[string]time.Time{
// 			"gt": greaterThanDate,
// 			"lt": lessThanDate,
// 		},
// 	})
// 	if err != nil {
// 		err = errors.Wrap(err, "Unable to find result for specific dates")
// 		log.Println(err)
// 		return nil, err
// 	}
// 	return result, nil
// }
