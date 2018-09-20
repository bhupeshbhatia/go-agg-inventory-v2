package controller

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/bhupeshbhatia/go-agg-inventory-v2/service"
	"github.com/pkg/errors"
)

func GetProductHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.Wrap(err, "Unable to read the request body")
		log.Println(err)
		return
	}

	// inventory, err := service.GetInventoryFromJSON(body)
	// if err != nil {
	// 	err = errors.Wrap(err, "Unable to unmarshal request body into Inventory struct")
	// 	log.Println(err)
	// 	return
	// }

	// mongoColl, err := service.GetMongoCollection()
	// if err != nil {
	// 	err = errors.Wrap(err, "Unable to connect to Mongo")
	// 	log.Println(err)
	// 	return
	// }

	_, err = service.GetFoodProducts(service.InventoryData{
		SearchField: string(body),
	})
	if err != nil {
		err = errors.Wrap(err, "Unable to find product")
		log.Println(err)
		return
	}
}

func AddProductHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.Wrap(err, "Unable to read the request body")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println(string(body))
	inventory, err := service.GetInventoryFromJSON(body)
	if err != nil {
		err = errors.Wrap(err, "Unable to unmarshal request body into Inventory struct")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	mongoColl, err := service.GetMongoCollection()
	if err != nil {
		err = errors.Wrap(err, "Unable to connect to Mongo")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = service.AddProduct(service.InventoryData{
		Product:         inventory,
		MongoCollection: mongoColl.Collection,
		FilterByName:    "item_id",
	})
	if err != nil {
		err = errors.Wrap(err, "Unable to add product")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
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

	inventory, err := service.GetInventoryFromJSON(body)
	if err != nil {
		err = errors.Wrap(err, "Unable to unmarshal request body into Inventory struct")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	mongoColl, err := service.GetMongoCollection()
	if err != nil {
		err = errors.Wrap(err, "Unable to connect to Mongo")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	inventory.Timestamp = time.Now().Unix()

	_, err = service.UpdateProduct(service.InventoryData{
		Product:         inventory,
		MongoCollection: mongoColl.Collection,
		FilterByItemId:  inventory.ItemID,
	})
	if err != nil {
		err = errors.Wrap(err, "Unable to add product")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.Wrap(err, "Unable to read the request body")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	inventory, err := service.GetInventoryFromJSON(body)
	if err != nil {
		err = errors.Wrap(err, "Unable to unmarshal request body into Inventory struct")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	mongoColl, err := service.GetMongoCollection()
	if err != nil {
		err = errors.Wrap(err, "Unable to connect to Mongo")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = service.DeleteProduct(service.InventoryData{
		Product:         inventory,
		MongoCollection: mongoColl.Collection,
		FilterByName:    "item_id",
	})
	if err != nil {
		err = errors.Wrap(err, "Unable to add product")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetProductRangeHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.Wrap(err, "Unable to read the request body")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	inventory, err := service.GetInventoryFromJSON(body)
	if err != nil {
		err = errors.Wrap(err, "Unable to unmarshal request body into Inventory struct")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println(inventory.ExpiryDate)

	st := time.Unix(inventory.ExpiryDate, 0)
	startTime := st.AddDate(0, 0, -20)

	stTime := startTime.Unix()

	t := time.Unix(stTime, 0)
	yesterdayTime := t.AddDate(0, 0, -40)
	yesTime := yesterdayTime.Unix()

	log.Println(startTime)

	mongoColl, err := service.GetMongoCollection()
	if err != nil {
		err = errors.Wrap(err, "Unable to connect to Mongo")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	getResults, err := service.GetProductInRange(service.InventoryData{
		Product:         inventory,
		MongoCollection: mongoColl.Collection,
		FilterByName:    "expiry_date",
		GetValue:        inventory.ExpiryDate,
		StartDate:       stTime,
		YesterdayTime:   yesTime,
	})
	if err != nil {
		err = errors.Wrap(err, "Unable to add product")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println(getResults)

	// for _, v := range getResults {
	// 	marshaledJSON := service.GetMarshal(*model.Inventory(v))
	// }

	w.WriteHeader(http.StatusOK)
}
