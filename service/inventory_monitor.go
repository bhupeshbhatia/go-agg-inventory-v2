package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/bhupeshbhatia/go-agg-inventory-v2/mockdata"

	"github.com/bhupeshbhatia/go-agg-inventory-v2/connectDB"
	"github.com/bhupeshbhatia/go-agg-inventory-v2/model"
	"github.com/pkg/errors"
)

type InvSearch struct {
	SearchTime int64 `json:"search_time"`
	TimePeriod int   `json:"time_period"`
}

func LoadData(w http.ResponseWriter, r *http.Request) {

	// file := strings.NewReader(mockdata.Testing())
	// var inv []model.Inventory
	// if err := json.NewDecoder(file).Decode(&inv); err != nil {
	// 	err = errors.Wrap(err, "Unable to load data")
	// 	log.Println(err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// DB connection
	Db, err := connectDB.ConfirmDbExists()
	if err != nil {
		err = errors.Wrap(err, "Mongo client unable to connect")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Convert body of type []byte into type []model.Inventory{}
	inventory := []model.Inventory{}
	err = json.Unmarshal(mockdata.StartUpLoadData(), &inventory)
	if err != nil {
		err = errors.Wrap(err, "Unable to unmarshal product into Inventory struct")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, val := range inventory {
		insertResult, err := Db.Collection.InsertOne(val)
		if err != nil {
			err = errors.Wrap(err, "Unable to insert event")
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Println(insertResult)

	}

	w.WriteHeader(http.StatusOK)
}

func GetBatchData(w http.ResponseWriter, r *http.Request) {
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
	})
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

func AddProductHandler(w http.ResponseWriter, r *http.Request) {
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

func UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
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
	if inventory.ItemID.String() == "" {
		w.WriteHeader(http.StatusInternalServerError)
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

func DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
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

func SearchInRange(w http.ResponseWriter, r *http.Request) {
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

	endTime := time.Unix(searchInv.SearchTime, 0)
	startTime := endTime.AddDate(0, 0, -(searchInv.TimePeriod)).Unix()

	findResults, err := Db.Collection.FindMap(map[string]interface{}{

		"timestamp": map[string]*int64{
			"$lt": &searchInv.SearchTime,
			"$gt": &startTime,
		},
	})
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
