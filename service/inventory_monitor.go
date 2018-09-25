package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo/findopt"

	"github.com/bhupeshbhatia/go-agg-inventory-v2/connectDB"
	"github.com/bhupeshbhatia/go-agg-inventory-v2/model"
	"github.com/pkg/errors"
)

type InvSearch struct {
	MaxTime          int64 `json:"max_time"`
	TimePeriodInDays int   `json:"time_period"`
}

func LoadDataInMongo(w http.ResponseWriter, r *http.Request) {

	// file := strings.NewReader(mockdata.Testing())
	// var inv []model.Inventory
	// if err := json.NewDecoder(file).Decode(&inv); err != nil {
	// 	err = errors.Wrap(err, "Unable to load data")
	// 	log.Println(err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// DB connection
	// Db, err := connectDB.ConfirmDbExists()
	// if err != nil {
	// 	err = errors.Wrap(err, "Mongo client unable to connect")
	// 	log.Println(err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	inventory := []model.Inventory{}
	for i := 0; i < 100; i++ {
		inventory = append(inventory, GenerateDataForInv())
	}

	// invBody := []byte(mockdata.StartUpLoadData())

	// //Convert body of type []byte into type []model.Inventory{}
	// inventory := []model.Inventory{}
	// err = json.Unmarshal(invBody, &inventory)
	// if err != nil {
	// 	err = errors.Wrap(err, "Unable to unmarshal product into Inventory struct")
	// 	log.Println(err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// for _, val := range inventory {
	// 	log.Println(val.ItemID)
	// 	insertResult, err := Db.Collection.InsertOne(val)
	// 	if err != nil {
	// 		err = errors.Wrap(err, "Unable to insert event")
	// 		log.Println(err)
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		return
	// 	}
	// 	log.Println(insertResult)

	// }

	jsonWithInvData, err := json.Marshal(&inventory)
	if err != nil {
		log.Println(err)
	}
	w.Write(jsonWithInvData)
	w.WriteHeader(http.StatusOK)
}

func GetInvFromMongo(w http.ResponseWriter, r *http.Request) {
	//Mongo collection
	Db, err := connectDB.ConfirmDbExists()
	if err != nil {
		err = errors.Wrap(err, "Mongo client unable to connect")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	timestamp := time.Now().Unix()
	log.Println(timestamp)

	//find results
	findResults, err := Db.Collection.FindMap(map[string]interface{}{

		"timestamp": map[string]int64{
			"$lt": timestamp,
		},
	},
		findopt.Limit(100),
	)
	if err != nil {
		err = errors.Wrap(err, "Error while fetching product.")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println(findResults)

	var marshalInventory []byte
	for _, v := range findResults {
		marshalInventory, err = json.Marshal(v)
		if err != nil {
			err = errors.Wrap(err, "Unable to marshal foodItem into Inventory struct")
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(marshalInventory)
	}
	w.WriteHeader(http.StatusOK)

}

func AddInventory(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.Wrap(err, "Unable to read the request body")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//DB connection
	Db, err := connectDB.ConfirmDbExists()
	if err != nil {
		err = errors.Wrap(err, "Mongo client unable to connect")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println(string(body))

	//Convert body of type []byte into type []model.Inventory{}
	inventory := []model.Inventory{}
	err = json.Unmarshal(body, &inventory)
	if err != nil {
		err = errors.Wrap(err, "Unable to unmarshal foodItem into Inventory struct")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, val := range inventory {
		if val.ItemID.String() != "" { //need to change this
			val.Timestamp = time.Now().Unix()
			log.Println(val.ItemID)
			insertResult, err := Db.Collection.InsertOne(val)
			if err != nil {
				err = errors.Wrap(err, "Unable to insert event")
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			log.Println(insertResult)
		}
	}

	w.WriteHeader(http.StatusOK)
}

func UpdateInventory(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.Wrap(err, "Unable to read the request body")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	Db, err := connectDB.ConfirmDbExists()
	if err != nil {
		err = errors.Wrap(err, "Mongo client unable to connect")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	inventory := &model.Inventory{}
	//Convert body of type []byte into type []model.Inventory{}
	err = json.Unmarshal(body, inventory)
	if err != nil {
		err = errors.Wrap(err, "Unable to unmarshal foodItem into Inventory struct")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Filter required for Update
	filter := &model.Inventory{
		ItemID: inventory.ItemID,
	}

	//Confirm that uuid is not empty
	if inventory.ItemID.String() == "" {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("UUID is empty")
		return
	}

	//Adding the timestamp
	nowTime := time.Now().Unix()
	inventory.Timestamp = nowTime

	update := &map[string]interface{}{
		"item_id":      inventory.ItemID,
		"name":         inventory.Name,
		"origin":       inventory.Origin,
		"date_arrived": inventory.DateArrived,
		"device_id":    inventory.DeviceID,
		"price":        inventory.Price,
		"total_weight": inventory.TotalWeight,
		"location":     inventory.Location,
	}

	updateResult, err := Db.Collection.UpdateMany(filter, update)
	if err != nil {
		err = errors.Wrap(err, "Unable to update event")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(updateResult)

	if updateResult.ModifiedCount > 0 {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func DeleteInventory(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.Wrap(err, "Unable to read the request body")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	Db, err := connectDB.ConfirmDbExists()
	if err != nil {
		err = errors.Wrap(err, "Mongo client unable to connect")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	inventory := &model.Inventory{}
	//Convert body of type []byte into type []model.Inventory{}
	err = json.Unmarshal(body, inventory)
	if err != nil {
		err = errors.Wrap(err, "Unable to unmarshal foodItem into Inventory struct")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if inventory.ItemID.String() == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	deleteResult, err := Db.Collection.DeleteMany(&model.Inventory{
		ItemID: inventory.ItemID,
	})
	if err != nil {
		err = errors.Wrap(err, "Unable to delete event")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if deleteResult.DeletedCount > 0 {
		w.WriteHeader(http.StatusOK)
	} else {
		log.Println("Unable to delete")
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func SearchInTimeRange(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.Wrap(err, "Unable to read the request body")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	Db, err := connectDB.ConfirmDbExists()
	if err != nil {
		err = errors.Wrap(err, "Mongo client unable to connect")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	searchInv := &InvSearch{}
	err = json.Unmarshal(body, searchInv)
	if err != nil {
		err = errors.Wrap(err, "Unable to unmarshal foodItem into Inventory struct")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	endTime := time.Unix(searchInv.MaxTime, 0)
	startTime := endTime.AddDate(0, 0, -(searchInv.TimePeriodInDays)).Unix()

	findResults, err := Db.Collection.FindMap(map[string]interface{}{

		"timestamp": map[string]*int64{
			"$lt": &searchInv.MaxTime,
			"$gt": &startTime,
		},
	},
		findopt.Limit(100),
	)
	if err != nil {
		err = errors.Wrap(err, "Error while fetching product.")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println(findResults)

	var marshalInventory []byte
	for _, v := range findResults {
		marshalInventory, err = json.Marshal(v)
		if err != nil {
			err = errors.Wrap(err, "Unable to marshal foodItem into Inventory struct")
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(marshalInventory)
	}
	w.WriteHeader(http.StatusOK)
}

func GetInvForToday(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.Wrap(err, "Unable to read the request body")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	Db, err := connectDB.ConfirmDbExists()
	if err != nil {
		err = errors.Wrap(err, "Mongo client unable to connect")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	searchInv := &InvSearch{}
	err = json.Unmarshal(body, searchInv)
	if err != nil {
		err = errors.Wrap(err, "Unable to unmarshal foodItem into Inventory struct")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	endTime := time.Unix(searchInv.MaxTime, 0)
	startTime := endTime.AddDate(0, 0, -(searchInv.TimePeriodInDays)).Unix()

	findResults, err := Db.Collection.FindMap(map[string]interface{}{

		"timestamp": map[string]*int64{
			"$lt": &searchInv.MaxTime,
			"$gt": &startTime,
		},
	})
	if err != nil {
		err = errors.Wrap(err, "Error while fetching product.")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println(len(findResults))
	resultCount := strconv.Itoa(len(findResults))
	w.Write([]byte(resultCount))
	w.WriteHeader(http.StatusOK)

	// pipeline := bson.NewArray(
	// 	bson.VC.DocumentFromElements(
	// 		bson.EC.SubDocumentFromElements(
	// 			"$match",
	// 			bson.EC.SubDocumentFromElements(
	// 				"timestamp",
	// 				bson.EC.Int64("$gte", startTime),
	// 				bson.EC.Int64("$lte", searchInv.MaxTime),
	// 			),
	// 		),
	// 	),
	// 	bson.VC.DocumentFromElements(
	// 		bson.EC.SubDocumentFromElements(
	// 			"$group",
	// 			bson.EC.SubDocumentFromElements(
	// 				"_id",
	// 				bson.EC.String("_id", nil),
	// 			),
	// 			bson.EC.SubDocumentFromElements(
	// 				"count",
	// 				bson.EC.count("count", 1),
	// 			),
	// 		),
	// 	),
	// )
	// aggResults, err := Db.Collection.Aggregate(pipeline)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// log.Println(len(aggResults))
	// for _, r := range aggResults {
	// 	// log.Println(r.(*item))
	// }
}

func TotalInventory(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.Wrap(err, "Unable to read the request body")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	Db, err := connectDB.ConfirmDbExists()
	if err != nil {
		err = errors.Wrap(err, "Mongo client unable to connect")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	searchInv := &InvSearch{}
	err = json.Unmarshal(body, searchInv)
	if err != nil {
		err = errors.Wrap(err, "Unable to unmarshal foodItem into Inventory struct")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	findResults, err := Db.Collection.FindMap(map[string]interface{}{

		"timestamp": map[string]*int64{
			"$lt": &searchInv.MaxTime,
		},
	})
	if err != nil {
		err = errors.Wrap(err, "Error while fetching product.")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resultCount := strconv.Itoa(len(findResults))
	w.Write([]byte(resultCount))
	w.WriteHeader(http.StatusOK)
}

func AvgInventoryPerHour(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.Wrap(err, "Unable to read the request body")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	Db, err := connectDB.ConfirmDbExists()
	if err != nil {
		err = errors.Wrap(err, "Mongo client unable to connect")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	searchInv := &InvSearch{}
	err = json.Unmarshal(body, searchInv)
	if err != nil {
		err = errors.Wrap(err, "Unable to unmarshal foodItem into Inventory struct")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	startTime := searchInv.MaxTime - 3600

	//here we need to send one hour, weekly, monthly

	findResults, err := Db.Collection.FindMap(map[string]interface{}{

		"timestamp": map[string]*int64{
			"$lte": &searchInv.MaxTime,
			"$gte": &startTime,
		},
	})
	if err != nil {
		err = errors.Wrap(err, "Error while fetching product.")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resultCount := strconv.Itoa(len(findResults))
	w.Write([]byte(resultCount))
	w.WriteHeader(http.StatusOK)
}

func TotalWeightSoldPerFruit(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.Wrap(err, "Unable to read the request body")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	Db, err := connectDB.ConfirmDbExists()
	if err != nil {
		err = errors.Wrap(err, "Mongo client unable to connect")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	searchInv := &InvSearch{}
	err = json.Unmarshal(body, searchInv)
	if err != nil {
		err = errors.Wrap(err, "Unable to unmarshal foodItem into Inventory struct")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	startTime := searchInv.MaxTime - 3600

	//here we need to send one hour, weekly, monthly

	findResults, err := Db.Collection.FindMap(map[string]interface{}{

		"timestamp": map[string]*int64{
			"$lte": &searchInv.MaxTime,
			"$gte": &startTime,
		},
	})
	if err != nil {
		err = errors.Wrap(err, "Error while fetching product.")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resultCount := strconv.Itoa(len(findResults))
	w.Write([]byte(resultCount))
	w.WriteHeader(http.StatusOK)
}
