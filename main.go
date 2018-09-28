package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bhupeshbhatia/go-agg-inventory-v2/service"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/urfave/negroni"
)

const AGGREGATE_ID = 2

func ErrorStackTrace(err error) string {
	return fmt.Sprintf("%+v\n", err)
}

func initRoutes() *mux.Router {
	router := mux.NewRouter()
	router = setAuthenticationRoute(router)
	return router
}

func setAuthenticationRoute(router *mux.Router) *mux.Router {
	router.HandleFunc("/add-product", service.AddInventory).Methods("POST", "OPTIONS")
	router.HandleFunc("/update-product", service.UpdateInventory).Methods("POST", "OPTIONS")
	router.HandleFunc("/delete-product", service.DeleteInventory).Methods("POST", "OPTIONS")
	router.HandleFunc("/search-range", service.TimeSearchInTable).Methods("POST", "OPTIONS")
	router.HandleFunc("/create-data", service.LoadDataInMongo).Methods("GET", "OPTIONS")
	router.HandleFunc("/load-table", service.LoadInventoryTable).Methods("POST", "OPTIONS")
	router.HandleFunc("/dist-weight", service.DistributionByWeight).Methods("GET", "OPTIONS")
	router.HandleFunc("/twsalewaste", service.TotalWeightSoldWasteDonatePerDay).Methods("POST", "OPTIONS")
	router.HandleFunc("/search-table", service.SearchInvTable).Methods("POST", "OPTIONS")

	router.HandleFunc("/perhr-sale", service.ProdSoldPerHour).Methods("POST", "OPTIONS")

	return router
}

func main() {
	err := godotenv.Load()
	if err != nil {
		err = errors.Wrap(err,
			".env file not found, env-vars will be read as set in environment",
		)
		log.Println(err)
	}

	// headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	// originsOk := handlers.AllowedOrigins([]string{("ORIGIN_ALLOWED")})
	// methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// service.GetProductCount()

	router := initRoutes()
	n := negroni.Classic()
	n.UseHandler(router)
	// http.ListenAndServe(":8080", handlers.CORS(originsOk, headersOk, methodsOk)(n))
	http.ListenAndServe(":8080", n)
}

//----------------------------------------------------------

//Calling Mongo
// Db, err := connectDB.ConfirmDbExists()
// if err != nil {
// 	err = errors.Wrap(err, "Mongo client unable to connect")
// 	log.Println(err)
// }

// mgCollection := Db.Collection

// mockData := mockdata.JsonForAddProduct()
// inventory, err := service.GetInventoryJSON([]byte(mockData))
// if err != nil {
// 	err = errors.Wrap(err, "Unable to unmarshal into Inventory struct")
// 	log.Println(err)
// }

// //Adding timestamp to inventory
// inventory.Timestamp = time.Now()

// insertData := &service.InventoryData{
// 	Product:     inventory,
// 	MongoTable:  mgCollection,
// 	FilterName:  "Fruit_ID",
// 	FilterValue: inventory.FruitID,
// }

// insertResult, err := service.AddProduct(*insertData)
// if err != nil {
// 	err = errors.Wrap(err, "Unable to insert event")
// 	log.Println(err)
// }

//Consumer
// config := consumer.Config{
// 	ConsumerGroup: "inventory.consumer.persistence",
// 	KafkaBrokers:  []string{"kafka:9092"},
// 	Topics:        []string{"event.rns_eventstore.events"},
// }

// kafka.KafkaConsumer(config)

// eventTopic := "events.rns_eventstore.events." + strconv.Itoa(AGGREGATE_ID)

// eventConfig := consumer.Config{
// 	ConsumerGroup: "events.rns_eventstore.eventsresponse",
// 	KafkaBrokers:  []string{"kafka:9092"},
// 	Topics:        []string{eventTopic},
// }

// kafka.KafkaConsumer(eventConfig)

// //Kafka Producer
// aggregateID := es.EventStoreQuery{
// 	AggregateID:      AGGREGATE_ID,
// 	AggregateVersion: 2,
// 	YearBucket:       2018,
// }

// producerJSON, err := json.Marshal(aggregateID)
// if err != nil {
// 	log.Println(err)
// }

// kafka.KafkaProducer(string(producerJSON))

// asyncProducer, err := kafka.ResponseProducer(kafka.KafAdapter{
// 	Address:         []string{"kafka:9092"},
// 	ProducerResChan: input,
// 	ResponseTopic:   "events.rns_eventstore.eventsquery",
// })
// if err != nil {
// 	err = errors.Wrap(err, "Unable to create producer")
// 	log.Println(ErrorStackTrace(err))
// }

// asyncProducer.EnableLogging()

// go func() {
// 	for err := asyncProducer.Errors() {
// 	  log.Println(err)
// 	}
//   }()

// aggregateID := es.EventStoreQuery{
// 	AggregateID:      AGGREGATE_ID,
// 	AggregateVersion: 2,
// 	YearBucket:       2018,
// }

// inputJson, err := json.Marshal(aggregateID)
// if err != nil {
// 	log.Println(err)
// }

// go func() {
// 	fmt.Println(asyncProducer)
// 	input <- &es.KafkaResponse{
// 		AggregateID: AGGREGATE_ID,
// 		Input:       string(inputJson),
// 	}
// }()
// time.Sleep(10 * time.Second)

//---------------------------------------------------------------------------------

// Creates a KafkaIO from KafkaAdapter based on set environment variables.

// input := make(chan *model.KafkaResponse)

// _, err := kafka.ResponseProducer(kafka.KafAdapter{
// 	Address:         []string{"kafka:9092"},
// 	ProducerResChan: input,
// 	ResponseTopic:   "test",
// })
// if err != nil {
// 	err = errors.Wrap(err, "Unable to create producer")
// 	log.Println(ErrorStackTrace(err))
// }

// go func() {
// 	// fmt.Println(produce)
// 	input <- &model.KafkaResponse{
// 		AggregateID: 1,
// 		Input:       "NOOOOOOOOOOO",
// 	}

// }()

//Consumer
// adap := kafka.KafkaConAdapter{
// 	Address:        []string{"kafka:9092"},
// 	ConsumerGroup:  "monitoring",
// 	ConsumerTopics: []string{"KafkaProducerTest"},
// }
// kio, err := kafka.Consume(&adap)
// go func() {
// 	for err := range kio.ConsumerErrChan {
// 		log.Println(err)
// 	}
// }()
// if err != nil {
// 	err = errors.Wrap(err, "Unable to create producer")
// 	log.Println(ErrorStackTrace(err))
// }

// var test []byte

// for msg := range kio.ConsumerMsgChan {
// 	log.Println(msg)
// 	test = msg.Value
// }

// fmt.Print(test)

// func CreateClientAndCollection() *mongo.Collection {
// 	client, err := connectDB.CreateClient()
// 	if err != nil {
// 		err = errors.Wrap(err, "Unable to get Mongo collection")
// 		log.Println(ErrorStackTrace(err))
// 	}

// 	mgTable, err := connectDB.CreateCollection(client, "users", "rns_aggregates")
// 	if err != nil {
// 		err = errors.Wrap(err, "Unable to insert in mongo")
// 		log.Println(ErrorStackTrace(err))
// 	}

// 	// aggVersion, err := events.GetMaxAggregateVersion(mgTable, aggregateID)
// 	// if err != nil {
// 	// 	err = errors.Wrap(err, "Mongo version not received")
// 	// 	log.Println(ErrorStackTrace(err))
// 	// }
// 	// return aggVersion

// 	return mgTable
// }

// mgTable := CreateClientAndCollection()
// inventory, err := service.GetInventoryJSON([]byte(mockdata.JsonForGetJSONString()))
// if err != nil {
// 	err = errors.Wrap(err, "Unable to unmarshal foodItem into Inventory struct")
// 	log.Println(err)
// }

// fmt.Println(inventory)

// _, err := service.GetInventoryJSON([]byte(mockdata.JsonForAddProduct()))
// if err != nil {
// 	err = errors.Wrap(err, "Unable to unmarshal addProduct json into Inventory struct")
// 	log.Println(err)
// }

// fmt.Printf("%+v", inventoryData)

// inv, err := service.GetMarshal(mockdata.InventoryMock())
// if err != nil {
// 	err = errors.Wrap(err, "Unable to unmarshal addProduct json into Inventory struct")
// 	log.Println(err)
// }

// testJson, err := service.GetInventoryJSON(inv)
// if err != nil {
// 	err = errors.Wrap(err, "Unable to unmarshal addProduct json into Inventory struct")
// 	log.Println(err)
// }

// fmt.Printf("%+v", testJson)

// timeWhenInserted := time.Now()

// inventoryInsert.Timestamp = timeWhenInserted

// insertData := &service.InventoryData{
// 	Product:     inventoryInsert,
// 	MongoTable:  mgTable,
// 	FilterName:  "Fruit_ID",
// 	FilterValue: inventoryInsert.FruitID,
// }

// insertResult, err := service.AddFood(*insertData)
// if err != nil {
// 	err = errors.Wrap(err, "Unable to unmarshal addProduct json into Inventory struct")
// 	log.Println(err)
// }
// fmt.Println("Insert: ", insertResult)

//=======================================================================================
//==========================================KAFKA==================================
