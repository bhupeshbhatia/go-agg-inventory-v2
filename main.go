package main

import (
	"fmt"
	"log"

	mongo "github.com/TerrexTech/go-mongoutils/mongo"
	"github.com/bhupeshbhatia/go-agg-inven-mongo-cmd/connectDB"
	"github.com/bhupeshbhatia/go-agg-inven-mongo-cmd/mockdata"
	"github.com/bhupeshbhatia/go-agg-inven-mongo-cmd/service"
	"github.com/pkg/errors"
)

func ErrorStackTrace(err error) string {
	return fmt.Sprintf("%+v\n", err)
}

func main() {
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

	inv, err := service.GetMarshal(mockdata.InventoryMock())
	if err != nil {
		err = errors.Wrap(err, "Unable to unmarshal addProduct json into Inventory struct")
		log.Println(err)
	}

	testJson, err := service.GetInventoryJSON(inv)
	if err != nil {
		err = errors.Wrap(err, "Unable to unmarshal addProduct json into Inventory struct")
		log.Println(err)
	}

	fmt.Printf("%+v", testJson)

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

	// input := make(chan *model.KafkaResponse)

	// _, err := kafka.CreateProducer(kafka.KafAdapter{
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

	// //Consumer
	// adap := kafka.KafkaConAdapter{
	// 	Address:        []string{"kafka:9092"},
	// 	ConsumerGroup:  "monitoring",
	// 	ConsumerTopics: []string{"KafkaProducerTest"},
	// }
	// kio, err := kafka.CreateConsumer(&adap)
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

}

func CreateClientAndCollection() *mongo.Collection {
	client, err := connectDB.CreateClient()
	if err != nil {
		err = errors.Wrap(err, "Unable to get Mongo collection")
		log.Println(ErrorStackTrace(err))
	}

	mgTable, err := connectDB.CreateCollection(client, "users", "rns_aggregates")
	if err != nil {
		err = errors.Wrap(err, "Unable to insert in mongo")
		log.Println(ErrorStackTrace(err))
	}

	// aggVersion, err := events.GetMaxAggregateVersion(mgTable, aggregateID)
	// if err != nil {
	// 	err = errors.Wrap(err, "Mongo version not received")
	// 	log.Println(ErrorStackTrace(err))
	// }
	// return aggVersion

	return mgTable
}
