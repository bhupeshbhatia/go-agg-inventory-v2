package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/bhupeshbhatia/go-agg-inventory-v2/connectDB"
	"github.com/bhupeshbhatia/go-agg-inventory-v2/model"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/pkg/errors"
)

type InvSearch struct {
	MaxTime          int64 `bson:"max_time,omitempty" json:"max_time,omitempty"`
	TimePeriodInDays int64 `bson:"time_days,omitempty" json:"time_days,omitempty"`
	TimeInHours      int64 `bson:"time_hours,omitempty json:"time_hours,omitempty"`
}

type InvDashboard struct {
	ProdName     string  `bson:"prod_name,omitempty" json:"prod_name,omitempty"`
	ProdWeight   float64 `bson:"prod_weight,omitempty" json:"prod_weight,omitempty"`
	TotalWeight  float64 `bson:"prod_weight,omitempty" json:"total_weight,omitempty"`
	SoldWeight   float64 `bson:"prod_weight,omitempty" json:"sold_weight,omitempty"`
	WasteWeight  float64 `bson:"prod_weight,omitempty" json:"waste_weight,omitempty"`
	ProductSold  int64   `bson:"prod_weight,omitempty" json:"product_sold,omitempty"`
	DonateWeight float64 `bson:"prod_weight,omitempty" json:"donate_product,omitempty"`
}

func LoadDataInMongo(w http.ResponseWriter, r *http.Request) {
	// DB connection
	Db, err := connectDB.ConfirmDbExists()
	if err != nil {
		err = errors.Wrap(err, "Mongo client unable to connect")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	inventory := []model.Inventory{}
	for i := 0; i < 100; i++ {
		inventory = append(inventory, GenerateDataForInv())
	}

	for _, val := range inventory {
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

	_, err = json.Marshal(&inventory)
	if err != nil {
		log.Println(err)
	}
	w.WriteHeader(http.StatusOK)
}

func LoadInventoryTable(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.Wrap(err, "Unable to read the request body")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	inventory := SearchByDays(body) //Just need max time

	if len(*inventory) > 0 {
		invJSON, err := json.Marshal(inventory)
		if err != nil {
			err = errors.Wrap(err, "Unable to marshal foodItem into Inventory struct")
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(invJSON)
	}
	w.WriteHeader(http.StatusInternalServerError)

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

//need end and start time (start in days)
func SearchTimeRange(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.Wrap(err, "Unable to read the request body")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	inventory := SearchByDays(body) //Just need max time

	if len(*inventory) > 0 {
		invJSON, err := json.Marshal(inventory)
		if err != nil {
			err = errors.Wrap(err, "Unable to marshal foodItem into Inventory struct")
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(invJSON)
	}

	// var marshalInventory []byte
	// for _, v := range findResults {
	// 	marshalInventory, err = json.Marshal(v)
	// 	if err != nil {
	// 		err = errors.Wrap(err, "Unable to marshal foodItem into Inventory struct")
	// 		log.Println(err)
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		return
	// 	}
	// 	w.Write(marshalInventory)
	// }
	// w.WriteHeader(http.StatusOK)
}

//find results for timestamp field within a specified time range
func SearchByDays(req []byte) *[]model.Inventory {
	Db, err := connectDB.ConfirmDbExists()
	if err != nil {
		err = errors.Wrap(err, "Mongo client unable to connect")
		log.Println(err)
		return nil
	}

	searchInv := []InvSearch{}
	var findResults []interface{}
	err = json.Unmarshal(req, &searchInv)
	if err != nil {
		err = errors.Wrap(err, "Unable to unmarshal foodItem into Inventory struct")
		log.Println(err)
		return nil
	}

	for _, searchVal := range searchInv {
		if searchVal.TimePeriodInDays != 0 {
			endTime := time.Unix(searchVal.MaxTime, 0)
			startTime := endTime.AddDate(0, 0, -(int(searchVal.TimePeriodInDays))).Unix()

			findResults, err = Db.Collection.FindMap(map[string]interface{}{

				"timestamp": map[string]*int64{
					"$lt": &searchVal.MaxTime,
					"$gt": &startTime,
				},
			})
		} else {
			findResults, err = Db.Collection.FindMap(map[string]interface{}{

				"timestamp": map[string]*int64{
					"$lt": &searchVal.MaxTime,
				},
			})
		}
		if err != nil {
			err = errors.Wrap(err, "Error while fetching product.")
			log.Println(err)
			return nil
		}
	}

	inventory := []model.Inventory{}

	for _, v := range findResults {
		resultInv := v.(*model.Inventory)
		inventory = append(inventory, *resultInv)
	}

	return &inventory

	// log.Println(len(findResults))
	// return findResults
}

//find results for timestamp field within a specified time range
func SearchByHours(req []byte) *[]model.Inventory {
	Db, err := connectDB.ConfirmDbExists()
	if err != nil {
		err = errors.Wrap(err, "Mongo client unable to connect")
		log.Println(err)
		return nil
	}

	searchInv := []InvSearch{}
	var findResults []interface{}
	err = json.Unmarshal(req, &searchInv)
	if err != nil {
		err = errors.Wrap(err, "Unable to unmarshal foodItem into Inventory struct")
		log.Println(err)
		return nil
	}
	for _, searchVal := range searchInv {
		if searchVal.TimeInHours != 0 {
			startTime := searchVal.MaxTime - 3600

			findResults, err = Db.Collection.FindMap(map[string]interface{}{

				"timestamp": map[string]*int64{
					"$lt": &searchVal.MaxTime,
					"$gt": &startTime,
				},
			})
		}
		if err != nil {
			err = errors.Wrap(err, "Error while fetching product.")
			log.Println(err)
			return nil
		}
	}

	inventory := []model.Inventory{}

	for _, v := range findResults {
		resultInv := v.(*model.Inventory)
		inventory = append(inventory, *resultInv)
	}

	return &inventory

	// log.Println(len(findResults))
	// return findResults
}

func TotalWeightSoldWasteDonatePerDay(w http.ResponseWriter, r *http.Request) {
	var totalWeight float64
	var soldWeight float64
	var wasteWeight float64
	var donateWeight float64

	//Get timestamp from frontend
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.Wrap(err, "Unable to read the request body")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Get list of inventory from SearchDay
	inventory := SearchByDays(body) //Just need max time

	if inventory == nil {
		log.Println(errors.New("Unable to get anything back from SearchWithTime function"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, v := range *inventory {
		totalWeight = v.TotalWeight + totalWeight
		soldWeight = v.SoldWeight + soldWeight
		wasteWeight = v.WasteWeight + wasteWeight
		donateWeight = v.DonateWeight + donateWeight
	}

	totalResult, err := json.Marshal(InvDashboard{
		TotalWeight:  totalWeight,
		SoldWeight:   soldWeight,
		WasteWeight:  wasteWeight,
		DonateWeight: donateWeight,
	})
	if err != nil {
		err = errors.Wrap(err, "Unable to create response body")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(totalResult)
	w.WriteHeader(http.StatusOK)
}

func TotalProductSoldGraph(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.Wrap(err, "Unable to read the request body")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	inventory := SearchByHours(body) //Just need max time

	// if len(*inventory) > 0 {

	// }

	if inventory == nil {
		log.Println(errors.New("Unable to get anything back from SearchWithTime function"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var soldWeight float64

	for _, v := range *inventory {
		soldWeight = v.SoldWeight + soldWeight
	}

	totalResult, err := json.Marshal(InvDashboard{
		SoldWeight: soldWeight,
	})
	if err != nil {
		err = errors.Wrap(err, "Unable to create response body")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(totalResult)
	w.WriteHeader(http.StatusOK)
}

// type bla struct {
// 	ID string `bson:"_id", json:"_id"`,
// 	Total float32 `bson:"total", json:"total"`
// }

//Need end time
func DistributionByWeight(w http.ResponseWriter, r *http.Request) {
	Db, err := connectDB.ConfirmDbExists()
	if err != nil {
		err = errors.Wrap(err, "Mongo client unable to connect")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pipeline := bson.NewArray(
		bson.VC.Document(
			bson.NewDocument(
				bson.EC.SubDocumentFromElements(
					"$group",
					bson.EC.String("_id", "$name"),
					bson.EC.SubDocumentFromElements(
						"total",
						bson.EC.String("$sum", "$total_weight"),
					),
				),
			),
		),
	)
	aggResults, err := Db.Collection.Aggregate(pipeline)
	if err != nil {
		log.Fatalln(err)
	}
	// log.Println(aggResults)

	dist := []InvDashboard{}

	for _, v := range aggResults {
		value := v.(map[string]interface{})
		strValue := value["_id"].(string)
		secValue := value["total"].(float64)
		dist = append(dist, InvDashboard{
			ProdName:   strValue,
			ProdWeight: secValue,
		})
	}

	DistWeightByte, err := json.Marshal(&dist)
	if err != nil {
		err = errors.Wrap(err, "Unable to marshal distribution")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(DistWeightByte)
	w.WriteHeader(http.StatusOK)
}

//need end and start time (start period has to be in number of days)
// func GetInvForToday(w http.ResponseWriter, r *http.Request) {
// 	body, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		err = errors.Wrap(err, "Unable to read the request body")
// 		log.Println(err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	Db, err := connectDB.ConfirmDbExists()
// 	if err != nil {
// 		err = errors.Wrap(err, "Mongo client unable to connect")
// 		log.Println(err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	searchInv := &InvSearch{}
// 	err = json.Unmarshal(body, searchInv)
// 	if err != nil {
// 		err = errors.Wrap(err, "Unable to unmarshal foodItem into Inventory struct")
// 		log.Println(err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	endTime := time.Unix(searchInv.MaxTime, 0)
// 	startTime := endTime.AddDate(0, 0, -(int(searchInv.TimePeriodInDays))).Unix()

// 	findResults, err := Db.Collection.FindMap(map[string]interface{}{

// 		"timestamp": map[string]*int64{
// 			"$lt": &searchInv.MaxTime,
// 			"$gt": &startTime,
// 		},
// 	})
// 	if err != nil {
// 		err = errors.Wrap(err, "Error while fetching product.")
// 		log.Println(err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	log.Println(len(findResults))
// 	resultCount := strconv.Itoa(len(findResults))
// 	w.Write([]byte(resultCount))
// 	w.WriteHeader(http.StatusOK)

// 	// pipeline := bson.NewArray(
// 	// 	bson.VC.DocumentFromElements(
// 	// 		bson.EC.SubDocumentFromElements(
// 	// 			"$match",
// 	// 			bson.EC.SubDocumentFromElements(
// 	// 				"timestamp",
// 	// 				bson.EC.Int64("$gte", startTime),
// 	// 				bson.EC.Int64("$lte", searchInv.MaxTime),
// 	// 			),
// 	// 		),
// 	// 	),
// 	// 	bson.VC.DocumentFromElements(
// 	// 		bson.EC.SubDocumentFromElements(
// 	// 			"$group",
// 	// 			bson.EC.SubDocumentFromElements(
// 	// 				"_id",
// 	// 				bson.EC.String("_id", nil),
// 	// 			),
// 	// 			bson.EC.SubDocumentFromElements(
// 	// 				"count",
// 	// 				bson.EC.count("count", 1),
// 	// 			),
// 	// 		),
// 	// 	),
// 	// )
// 	// aggResults, err := Db.Collection.Aggregate(pipeline)
// 	// if err != nil {
// 	// 	log.Fatalln(err)
// 	// }
// 	// log.Println(len(aggResults))
// 	// for _, r := range aggResults {
// 	// 	// log.Println(r.(*item))
// 	// }
// }

// //need end time
// func TotalInventory(w http.ResponseWriter, r *http.Request) {
// 	body, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		err = errors.Wrap(err, "Unable to read the request body")
// 		log.Println(err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	Db, err := connectDB.ConfirmDbExists()
// 	if err != nil {
// 		err = errors.Wrap(err, "Mongo client unable to connect")
// 		log.Println(err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	searchInv := &InvSearch{}
// 	err = json.Unmarshal(body, searchInv)
// 	if err != nil {
// 		err = errors.Wrap(err, "Unable to unmarshal foodItem into Inventory struct")
// 		log.Println(err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	findResults, err := Db.Collection.FindMap(map[string]interface{}{

// 		"timestamp": map[string]*int64{
// 			"$lt": &searchInv.MaxTime,
// 		},
// 	})
// 	if err != nil {
// 		err = errors.Wrap(err, "Error while fetching product.")
// 		log.Println(err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	resultCount := strconv.Itoa(len(findResults))
// 	w.Write([]byte(resultCount))
// 	w.WriteHeader(http.StatusOK)
// }

// //need to think about how to show it in per hour - require start and end time
// func AvgInventoryPerHour(w http.ResponseWriter, r *http.Request) {
// 	body, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		err = errors.Wrap(err, "Unable to read the request body")
// 		log.Println(err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	Db, err := connectDB.ConfirmDbExists()
// 	if err != nil {
// 		err = errors.Wrap(err, "Mongo client unable to connect")
// 		log.Println(err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	searchInv := &InvSearch{}
// 	err = json.Unmarshal(body, searchInv)
// 	if err != nil {
// 		err = errors.Wrap(err, "Unable to unmarshal foodItem into Inventory struct")
// 		log.Println(err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	startTime := searchInv.MaxTime - 3600

// 	//here we need to send one hour, weekly, monthly

// 	findResults, err := Db.Collection.FindMap(map[string]interface{}{

// 		"timestamp": map[string]*int64{
// 			"$lte": &searchInv.MaxTime,
// 			"$gte": &startTime,
// 		},
// 	})
// 	if err != nil {
// 		err = errors.Wrap(err, "Error while fetching product.")
// 		log.Println(err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	resultCount := strconv.Itoa(len(findResults))
// 	w.Write([]byte(resultCount))
// 	w.WriteHeader(http.StatusOK)
// }

// // func TotalWeightSoldPerFruit(w http.ResponseWriter, r *http.Request) {
// // 	body, err := ioutil.ReadAll(r.Body)
// // 	if err != nil {
// // 		err = errors.Wrap(err, "Unable to read the request body")
// // 		log.Println(err)
// // 		w.WriteHeader(http.StatusInternalServerError)
// // 		return
// // 	}

// // 	Db, err := connectDB.ConfirmDbExists()
// // 	if err != nil {
// // 		err = errors.Wrap(err, "Mongo client unable to connect")
// // 		log.Println(err)
// // 		w.WriteHeader(http.StatusInternalServerError)
// // 		return
// // 	}

// // 	searchInv := &InvSearch{}
// // 	err = json.Unmarshal(body, searchInv)
// // 	if err != nil {
// // 		err = errors.Wrap(err, "Unable to unmarshal foodItem into Inventory struct")
// // 		log.Println(err)
// // 		w.WriteHeader(http.StatusInternalServerError)
// // 		return
// // 	}

// // 	startTime := searchInv.MaxTime - 3600

// // 	//here we need to send one hour, weekly, monthly

// // 	findResults, err := Db.Collection.FindMap(map[string]interface{}{

// // 		"timestamp": map[string]*int64{
// // 			"$lte": &searchInv.MaxTime,
// // 			"$gte": &startTime,
// // 		},
// // 	})
// // 	if err != nil {
// // 		err = errors.Wrap(err, "Error while fetching product.")
// // 		log.Println(err)
// // 		w.WriteHeader(http.StatusInternalServerError)
// // 		return
// // 	}
// // 	resultCount := strconv.Itoa(len(findResults))
// // 	w.Write([]byte(resultCount))
// // 	w.WriteHeader(http.StatusOK)
// // }
