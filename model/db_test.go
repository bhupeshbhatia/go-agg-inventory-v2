package model

import (
	ctx "context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"testing"
	"time"

	mongo "github.com/TerrexTech/go-mongoutils/mongo"
	"github.com/TerrexTech/uuuid"
	mgo "github.com/mongodb/mongo-go-driver/mongo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

func TestBooks(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Inventory Suite")
}

// newTimeoutContext creates a new WithTimeout context with specified timeout.
func newTimeoutContext(timeout uint32) (ctx.Context, ctx.CancelFunc) {
	return ctx.WithTimeout(
		ctx.Background(),
		time.Duration(timeout)*time.Millisecond,
	)
}

type Env struct {
	DbTest Datastore
}

type mockDb struct {
	collection *mongo.Collection
}

var _ = Describe("Mongo service test", func() {
	var (
		// jsonString string
		mgTable         *mongo.Collection
		client          mongo.Client
		resourceTimeout uint32
		testDatabase    string
		clientConfig    mongo.ClientConfig
		c               *mongo.Collection
		dataCount       int
		unixTime        int64
	)

	testDatabase = "rns_test"
	resourceTimeout = 3000
	dataCount = 5

	dropTestDatabase := func() {
		client, err := mongo.NewClient(clientConfig)
		Expect(err).ToNot(HaveOccurred())

		dbCtx, dbCancel := newTimeoutContext(resourceTimeout)
		err = client.Database(testDatabase).Drop(dbCtx)
		dbCancel()
		Expect(err).ToNot(HaveOccurred())

		err = client.Disconnect()
		Expect(err).ToNot(HaveOccurred())
	}

	BeforeEach(func() {
		clientConfig = mongo.ClientConfig{
			Hosts:               []string{"127.0.0.1:27017"},
			Username:            "root",
			Password:            "root",
			TimeoutMilliseconds: 5000,
		}
		client, err := mongo.NewClient(clientConfig)
		Expect(err).ToNot(HaveOccurred())

		conn := &mongo.ConnectionConfig{
			Client:  client,
			Timeout: 1000,
		}

		c = &mongo.Collection{
			Connection:   conn,
			Name:         "invtest",
			Database:     testDatabase,
			SchemaStruct: &Inventory{},
		}

		mgTable, err = mongo.EnsureCollection(c)
		Expect(err).ToNot(HaveOccurred())

		inv := []Inventory{}
		for i := 0; i < dataCount; i++ {
			inv = append(inv, GenerateDataForInv())
		}

		var ins []interface{}
		for _, v := range inv {
			ins = append(ins, v)
		}

		_, err = mgTable.InsertMany(ins)
		Expect(err).ToNot(HaveOccurred())

		unixTime = time.Now().AddDate(0, 0, 10).Unix()

	})

	AfterEach(func() {
		dropTestDatabase()
		err := client.Disconnect()
		Expect(err).ToNot(HaveOccurred())
	})

	It("Create data to insert in Mongo", func() {
		inv := []Inventory{}
		for i := 0; i < dataCount; i++ {
			inv = append(inv, GenerateDataForInv())
		}

		var ins []interface{}
		for _, v := range inv {
			ins = append(ins, v)
		}

		_, err := mgTable.InsertMany(ins)
		Expect(err).ToNot(HaveOccurred())
	})

	It("Creates UPC of length 12", func() {
		num := GenFakeBarcode("upc")
		strNum := strconv.Itoa(int(num))
		Expect(len(strNum)).Should(Equal(12))

	})

	It("Creates SKU of length 12", func() {
		num := GenFakeBarcode("sku")
		strNum := strconv.Itoa(int(num))
		Expect(len(strNum)).Should(Equal(8))

	})

	It("creates a json from []inventory{}", func() {
		rsUuid, err := uuuid.NewV4()
		dUuid, err := uuuid.NewV4()
		Expect(err).ToNot(HaveOccurred())
		_, err = json.Marshal(&Inventory{
			RsCustomerID: rsUuid,
			DeviceID:     dUuid,
		})
		Expect(err).ToNot(HaveOccurred())
	})

	It("Should unmarshal without error", func() {
		// timeToStr := strconv.Itoa(int(unixTime))
		req := fmt.Sprintf("{\"end_date\": %d}", unixTime)
		// req := `{"end_date": 1540011330}`
		searchUpLimit := InvSearch{}
		err := json.Unmarshal([]byte(req), &searchUpLimit)
		Expect(err).ToNot(HaveOccurred())
		Expect(searchUpLimit.EndDate).To(Equal(unixTime))
	})

	It("Should unmarshal with error if json is not built properly", func() {
		req := fmt.Sprintf("[{\"end_date\": %d}]", unixTime)
		// req := `[{"end_date": 1540011330}]`
		searchUpLimit := InvSearch{}
		err := json.Unmarshal([]byte(req), &searchUpLimit)
		Expect(err).To(HaveOccurred())
	})

	It("Should get no results from db if no match found", func() {
		req := `[{"end_date": 15400}]`
		searchUpLimit := []InvSearch{}
		err := json.Unmarshal([]byte(req), &searchUpLimit)
		Expect(err).ToNot(HaveOccurred())

		var findResults []interface{}
		for _, v := range searchUpLimit {
			findResults, err = mgTable.Find(map[string]interface{}{
				"timestamp": map[string]int64{
					"$lte": v.EndDate,
				},
			})
		}

		Expect(err).ToNot(HaveOccurred())
		Expect(findResults).To(HaveLen(0))
	})

	It("Should not give an error if end_date is 0", func() {
		req := `[{"end_date": 0}]`
		searchUpLimit := []InvSearch{}
		err := json.Unmarshal([]byte(req), &searchUpLimit)
		Expect(err).ToNot(HaveOccurred())
		Expect(searchUpLimit[0].EndDate).To(Equal(int64(0)))

		var findResults []interface{}
		for _, v := range searchUpLimit {
			findResults, err = mgTable.Find(map[string]interface{}{
				"timestamp": map[string]int64{
					"$lte": v.EndDate,
				},
			})
		}
		Expect(err).ToNot(HaveOccurred())
		Expect(findResults).To(HaveLen(0))
	})

	It("Should get results from db", func() {
		req := fmt.Sprintf("[{\"end_date\": %d}]", unixTime)

		// req := `{"end_date": 1540191465}`
		searchUpLimit := []InvSearch{}
		err := json.Unmarshal([]byte(req), &searchUpLimit)
		Expect(err).ToNot(HaveOccurred())

		var findResults []interface{}
		for _, v := range searchUpLimit {
			findResults, err = mgTable.Find(map[string]interface{}{
				"timestamp": map[string]int64{
					"$lte": v.EndDate,
				},
			})
		}
		Expect(err).ToNot(HaveOccurred())
		Expect(findResults).To(HaveLen(dataCount))
	})

	It("Should get marshalled results", func() {
		req := fmt.Sprintf("[{\"end_date\": %d}]", unixTime)
		// req := `{"end_date": 1540011330}`
		searchUpLimit := []InvSearch{}
		err := json.Unmarshal([]byte(req), &searchUpLimit)
		Expect(err).ToNot(HaveOccurred())

		var findResults []interface{}
		for _, v := range searchUpLimit {
			findResults, err = mgTable.Find(map[string]interface{}{
				"timestamp": map[string]int64{
					"$lte": v.EndDate,
				},
			})
		}
		Expect(err).ToNot(HaveOccurred())
		Expect(findResults).To(HaveLen(dataCount))

		inventory := []Inventory{}

		for _, v := range findResults {
			resultInv := v.(*Inventory)
			inventory = append(inventory, *resultInv)
		}

		//Marshal into []byte
		totalResult, err := json.Marshal(&inventory)
		Expect(err).ToNot(HaveOccurred())
		log.Println(reflect.TypeOf(totalResult))
	})

	It("Should not get an error when passing zero value", func() {
		req := `[{"end_date": 0}]`
		searchUpLimit := []InvSearch{}
		err := json.Unmarshal([]byte(req), &searchUpLimit)
		Expect(err).ToNot(HaveOccurred())

		var findResults []interface{}
		for _, v := range searchUpLimit {
			findResults, err = mgTable.Find(map[string]interface{}{
				"timestamp": map[string]int64{
					"$lte": v.EndDate,
				},
			})
		}
		Expect(err).ToNot(HaveOccurred())
		Expect(findResults).To(HaveLen(0))

		inventory := []Inventory{}

		for _, v := range findResults {
			resultInv := v.(*Inventory)
			inventory = append(inventory, *resultInv)
		}

		//Marshal into []byte
		totalResult, err := json.Marshal(&inventory)
		Expect(err).ToNot(HaveOccurred())
		log.Println(reflect.TypeOf(totalResult))
	})

	It("Should not give an if all values in InvSearch are empty", func() {
		req := `[{"end_date": 0, "start_date": 0, "search_key": "", "search_val": ""}]`
		db := &Db{mgTable}
		env := &Env{db}
		_, err := env.DbTest.SearchDb([]byte(req))
		Expect(err).ToNot(HaveOccurred())
	})

	// //THIS NEEDS TO BE WORKED ON
	// // It("Should give an error if searchKey and value have wrong types", func() {
	// // 	req := `{"search_key": "sku", "search_val": 2}`
	// // 	db := &Db{mgTable}
	// // 	env := &Env{db}
	// // 	_, err := env.DbTest.SearchDb([]byte(req))
	// // 	Expect(err).To(HaveOccurred())
	// // })

	It("Should not add inventory to db", func() {
		req := `{"upc": 0, "sku": 0, "name": "", "origin": "", "total_weight": 0, "price": 0, "location": "", "date_arrived": 0, "expiry_date":0, "timestamp": 0}`
		db := &Db{mgTable}
		env := &Env{db}
		err := env.DbTest.AddInventory([]byte(req))
		Expect(err).ToNot(HaveOccurred())
	})

	It("Should add inventory in db", func() {
		// req := fmt.Sprintf("{\"end_date\": %d}", unixTime)

		req := fmt.Sprintf(`{"upc": 222222222222, "sku": 22222222, "name": "Apple", "origin": "Canada", "total_weight": 12, "price": 34, "location": "M201", "date_arrived": %d, "expiry_date":%d, "timestamp":%d}`, unixTime, unixTime, unixTime)
		db := &Db{mgTable}
		env := &Env{db}
		err := env.DbTest.AddInventory([]byte(req))
		Expect(err).ToNot(HaveOccurred())
	})

	It("Should not update inventory in db", func() {
		req := `{"upc": 0, "sku": 0, "name": "", "origin": "", "total_weight": 0, "price": 0, "location": "", "date_arrived": 0, "expiry_date":0, "timestamp": 0}`
		db := &Db{mgTable}
		env := &Env{db}
		_, err := env.DbTest.UpdateInventory([]byte(req))
		Expect(err).To(HaveOccurred())
	})

	// //CHECK WITH THIS---- SHOULD IT GIVE AN ERROR
	// It("Should not give an error if upc has a string", func() {
	// 	req := `{"upc": "1540011330"}`
	// 	inventory := Inventory{}
	// 	err := json.Unmarshal([]byte(req), &inventory)
	// 	Expect(err).To(BeNil())
	// })

	It("Should NOT give an error if sku has a string", func() {
		req := `{"sku": "15400113"}`
		inventory := Inventory{}
		err := json.Unmarshal([]byte(req), &inventory)
		Expect(err).To(BeNil())
	})

	It("Should not give an error if date arrived has a string", func() {
		req := `{"date_arrived": "15400113"}`
		inventory := Inventory{}
		err := json.Unmarshal([]byte(req), &inventory)
		Expect(err).To(BeNil())
	})

	It("Should not give an error if price has a string", func() {
		req := `{"price": "15400113"}`
		inventory := Inventory{}
		err := json.Unmarshal([]byte(req), &inventory)
		Expect(err).To(BeNil())
	})

	It("Should not give error if total weight has a string", func() {
		req := `{"total weight": "1540"}`
		inventory := Inventory{}
		err := json.Unmarshal([]byte(req), &inventory)
		Expect(err).To(BeNil())
	})

	It("should give an error when UPC is 0", func() {
		req := `{"upc": 0}`
		inventory := Inventory{}
		err := json.Unmarshal([]byte(req), &inventory)
		if inventory.UPC == 0 {
			err = errors.New("UPC is 0")
		}
		Expect(inventory.UPC).To(Equal(int64(0)))
		Expect(err).To(HaveOccurred())
	})

	It("should add inventory in db", func() {
		req := fmt.Sprintf(`{"upc": 222222222232, "sku": 22222211, "name": "Blah", "origin": "Canada", "total_weight": 12, "price": 34, "location": "M201", "date_arrived": %d, "expiry_date":%d, "timestamp":%d}`, unixTime, unixTime, unixTime)
		db := &Db{mgTable}
		env := &Env{db}
		err := env.DbTest.AddInventory([]byte(req))
		Expect(err).ToNot(HaveOccurred())

		searchAdd := `[{"search_key": "name", "search_val": "Blah"}]`
		results, err := env.DbTest.SearchDb([]byte(searchAdd))
		Expect(err).ToNot(HaveOccurred())

		log.Println(string(results))

		//NEED TO ASK ABOUT THIS-=-===================

		// inventory := []Inventory{}
		// err = json.Unmarshal(results, inventory)
		// // err = json.Unmarshal(results, &inventory)
		// Expect(err).To(HaveOccurred())
	})

	It("Should update inventory in db", func() {
		//Add first
		dUuid, err := uuuid.NewV4()
		inventory := Inventory{
			DeviceID: dUuid,
		}
		_, err = json.Marshal(&inventory)
		Expect(err).ToNot(HaveOccurred())

		req := fmt.Sprintf(`{"upc": 222222222232, "sku": 22222211, "name": "Blah", "origin": "Canada", "total_weight": 12, "price": 34, "location": "M201", "date_arrived": %d, "expiry_date":%d, "timestamp":%d, "device_id": "%s"}`, unixTime, unixTime, unixTime, inventory.DeviceID)
		db := &Db{mgTable}
		env := &Env{db}
		err = env.DbTest.AddInventory([]byte(req))
		Expect(err).ToNot(HaveOccurred())

		req = fmt.Sprintf(`{"upc": 222222222232, "sku": 22222211, "name": "Blah", "origin": "NOT Canada", "total_weight": 12, "price": 34, "location": "M201", "date_arrived": %d, "expiry_date":%d, "timestamp":%d, "device_id": "%s"}`, unixTime, unixTime, unixTime, inventory.DeviceID)

		upResult, err := env.DbTest.UpdateInventory([]byte(req))
		Expect(err).ToNot(HaveOccurred())
		Expect(upResult.(*mgo.UpdateResult).ModifiedCount).To(Equal(int64(1)))
	})

	It("Should unmarshal successfully", func() {
		//Add first
		dUuid, err := uuuid.NewV4()
		inventory := Inventory{
			DeviceID: dUuid,
		}
		_, err = json.Marshal(&inventory)
		Expect(err).ToNot(HaveOccurred())

		req := fmt.Sprintf(`[{"upc": 222222222232, "sku": 22222211, "name": "Blah", "origin": "Canada", "total_weight": 12, "price": 34, "location": "M201", "date_arrived": %d, "expiry_date":%d, "timestamp":%d, "device_id": "%s"}]`, unixTime, unixTime, unixTime, inventory.DeviceID)

		invDel := []Inventory{}
		//Convert body of type []byte into type []model.Inventory{}
		err = json.Unmarshal([]byte(req), &invDel)
		Expect(err).ToNot(HaveOccurred())
	})

	It("Should not delete successfully if upc and sku doesn't match", func() {
		//Add first
		dUuid, err := uuuid.NewV4()
		inventory := Inventory{
			DeviceID: dUuid,
		}
		_, err = json.Marshal(&inventory)
		Expect(err).ToNot(HaveOccurred())

		req := fmt.Sprintf(`{"upc": 222222222232, "sku": 22222211, "name": "Blah", "origin": "Canada", "total_weight": 12, "price": 34, "location": "M201", "date_arrived": %d, "expiry_date":%d, "timestamp":%d, "device_id": "%s"}`, unixTime, unixTime, unixTime, inventory.DeviceID)
		db := &Db{mgTable}
		env := &Env{db}
		err = env.DbTest.AddInventory([]byte(req))
		Expect(err).ToNot(HaveOccurred())

		delBody := `[{"upc": 22222222, "sku": 22222211}]`

		deleteResult, err := env.DbTest.DeleteInventory([]byte(delBody))
		Expect(err).ToNot(HaveOccurred())
		Expect(string(deleteResult)).To(Equal("0"))
	})

	It("Should delete successfully", func() {
		//Add first
		dUuid, err := uuuid.NewV4()
		inventory := Inventory{
			DeviceID: dUuid,
		}
		_, err = json.Marshal(&inventory)
		Expect(err).ToNot(HaveOccurred())

		req := fmt.Sprintf(`{"upc": 222222222232, "sku": 22222211, "name": "Blah", "origin": "Canada", "total_weight": 12, "price": 34, "location": "M201", "date_arrived": %d, "expiry_date":%d, "timestamp":%d, "device_id": "%s"}`, unixTime, unixTime, unixTime, inventory.DeviceID)
		db := &Db{mgTable}
		env := &Env{db}
		err = env.DbTest.AddInventory([]byte(req))
		Expect(err).ToNot(HaveOccurred())

		delBody := `[{"upc": 222222222232, "sku": 22222211}]`

		deleteResult, err := env.DbTest.DeleteInventory([]byte(delBody))
		Expect(err).ToNot(HaveOccurred())
		// delNum := strconv.Itoa(1)
		Expect(string(deleteResult)).To(Equal("1"))
	})

	// It("Should successfully unmarshal after adding ", func() {
	// 	//Add first
	// 	dUuid, err := uuuid.NewV4()
	// 	inventory := Inventory{
	// 		DeviceID: dUuid,
	// 	}
	// 	_, err = json.Marshal(&inventory)
	// 	Expect(err).ToNot(HaveOccurred())

	// 	req := fmt.Sprintf(`{"upc": 222222222232, "sku": 22222211, "name": "Blah", "origin": "Canada", "total_weight": 12, "price": 34, "location": "M201", "date_arrived": %d, "expiry_date":%d, "timestamp":%d, "device_id": "%s", "sold_weight": 8, "waste_weight": 2, "donate_weight": 2}`, unixTime, unixTime, unixTime, inventory.DeviceID)
	// 	db := &Db{mgTable}
	// 	env := &Env{db}
	// 	err = env.DbTest.AddInventory([]byte(req))
	// 	Expect(err).ToNot(HaveOccurred())

	// 	req = fmt.Sprintf(`[{"start_date":%d, "end_date": %d}]`, int64(0), unixTime)
	// 	invDel := []Inventory{}
	// 	//Convert body of type []byte into type []model.Inventory{}
	// 	err = json.Unmarshal([]byte(req), &invDel)
	// 	if err != nil {
	// 		err = errors.Wrap(err, "Unable to unmarshal - DeleteInventory")
	// 		log.Println(err)
	// 	}
	// 	Expect(err).ToNot(HaveOccurred())
	// })

	It("Should unmarshal successfully", func() {
		// var totalWeight float64
		// var soldWeight float64
		// var wasteWeight float64
		// var donateWeight float64

		// var tweight []float64
		// var sweight []float64
		// var wweight []float64
		// var dweight []float64

		endDate := time.Now().AddDate(0, 0, 20).Unix()
		startDate := time.Now().AddDate(0, 0, -15).Unix()

		//Add first
		dUuid, err := uuuid.NewV4()
		inventory := Inventory{
			DeviceID: dUuid,
		}
		_, err = json.Marshal(&inventory)
		Expect(err).ToNot(HaveOccurred())

		req := fmt.Sprintf(`{"upc": 222222222232, "sku": 22222211, "name": "Blah", "origin": "Canada", "total_weight": 12, "price": 34, "location": "M201", "date_arrived": %d, "expiry_date":%d, "timestamp":%d, "device_id": "%s", "sold_weight": 8, "waste_weight": 2, "donate_weight": 2}`, unixTime, unixTime, unixTime, inventory.DeviceID)
		db := &Db{mgTable}
		env := &Env{db}
		err = env.DbTest.AddInventory([]byte(req))
		Expect(err).ToNot(HaveOccurred())

		newReq := fmt.Sprintf(`[{"start_date":%d, "end_date": %d}]`, startDate, endDate)

		searchInv, err := env.DbTest.SearchDb([]byte(newReq))
		Expect(err).ToNot(HaveOccurred())
		log.Println(string(searchInv))

		invForGraph := []Inventory{}
		// Convert body of type []byte into type []model.Inventory{}
		err = json.Unmarshal(searchInv, &invForGraph)
		Expect(err).ToNot(HaveOccurred())

		// for _, v := range invForGraph {
		// 	totalWeight = v.TotalWeight + totalWeight
		// 	tweight = append(tweight, totalWeight)

		// 	soldWeight = v.SoldWeight + soldWeight
		// 	sweight = append(sweight, soldWeight)

		// 	wasteWeight = v.WasteWeight + wasteWeight
		// 	wweight = append(wweight, wasteWeight)

		// 	donateWeight = v.DonateWeight + donateWeight
		// 	dweight = append(dweight, donateWeight)
		// }

		// log.Println(totalWeight)

		// invSearch := []InvSearch{}
		// err = json.Unmarshal(searchInv, &invSearch)
		// Expect(err).ToNot(HaveOccurred())

		// log.Println(invSearch)
		// // var totalResult []byte

		// dash := make(map[int]InvDashboard)
		// // dash := []InvDashboard{}
		// for i, v := range invSearch {
		// 	dash[i] = InvDashboard{
		// 		TotalWeight:  tweight[i],
		// 		SoldWeight:   sweight[i],
		// 		WasteWeight:  wweight[i],
		// 		DonateWeight: dweight[i],
		// 		Dates:        v.StartDate,
		// 	}
		// }
		// totalResult, err := json.Marshal(dash)
		// Expect(err).ToNot(HaveOccurred())
		// log.Println(totalResult)

	})
})
