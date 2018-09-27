package service

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/TerrexTech/uuuid"
	"github.com/pkg/errors"

	"github.com/bhupeshbhatia/go-agg-inventory-v2/model"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type ModifyInvData struct {
	Inv        model.Inventory
	Datearr    int64
	Expirydate int64
	Timestamp  int64
	Randnum    int
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func generateRandomValue(num1, num2 int) int {
	// rand.Seed(time.Now().Unix())
	return random(num1, num2)
}

func generateNewUUID() uuuid.UUID {
	uuid, err := uuuid.NewV4()
	if err != nil {
		err = errors.Wrap(err, "Unable to generate UUID")
		log.Println(err)
	}
	return uuid
}

var productsName = []string{"Banana", "Orange", "Apple", "Mango", "Strawberry", "Tomato", "Lettuce", "Pear", "Grapes", "Sweet Pepper"}
var locationName = []string{"A101", "B201", "O301", "M401", "S501", "T601", "L701", "P801", "G901", "SW1001"}
var provinceNames = []string{"ON Canada", "BC Canada", "SK Canada", "MN Canada", "NS Canada", "PEI Canada", "QC Canada"}

func GenerateDataForInv() model.Inventory {

	randNameAndLocation := generateRandomValue(1, 10)
	randOrigin := generateRandomValue(1, 4)
	randDateArr := generateRandomValue(1, 7)                          //in hours
	randTimestamp := generateRandomValue(randDateArr, randDateArr+2)  //in hours
	randExpiry := generateRandomValue(int((randTimestamp/24)+1), 21)  //in days
	randDatesold := generateRandomValue(randTimestamp, randExpiry*24) //in hours
	randPrice := generateRandomValue(5000, 10000)
	randTotalWeight := generateRandomValue(100, 300)
	randWasteWeight := generateRandomValue(1, 80)

	inventory := model.Inventory{
		ItemID:       generateNewUUID(),
		RsCustomerID: generateNewUUID(),
		DeviceID:     generateNewUUID(),
		Name:         productsName[randNameAndLocation-1], //-1 because rand starts from 1
		Origin:       provinceNames[randOrigin-1],
		TotalWeight:  float64(randTotalWeight),
		Price:        float64(randPrice),
		Location:     locationName[randNameAndLocation],
		WasteWeight:  float64(randWasteWeight - 1),
		DonateWeight: float64(generateRandomValue(1, 21)),
		AggregateID:  2,
		DateArrived:  time.Now().Add(time.Duration(randDateArr) * time.Hour).Unix(),
		ExpiryDate:   time.Now().AddDate(0, 0, randExpiry).Unix(),
		Timestamp:    time.Now().Add(time.Duration(randTimestamp) * time.Hour).Unix(),
		DateSold:     time.Now().Add(time.Duration(randDatesold) * time.Hour).Unix(),
		SalePrice:    float64(generateRandomValue(2, 4)),
		SoldWeight:   float64(generateRandomValue(randWasteWeight, randTotalWeight)),
	}
	return inventory
}

func TestIfDataGenerated() {
	inventory := []model.Inventory{}
	for i := 0; i < 100; i++ {
		inventory = append(inventory, GenerateDataForInv())
	}

	jsonWithInvData, err := json.Marshal(&inventory)
	if err != nil {
		log.Println(err)
	}
	log.Println(jsonWithInvData)
}

// func GetProdAndTotalWeight(prodName string, inventory model.Inventory) (float64, float64) {
// 	// productWeight := 0
// 	// totalWeight := 0
// 	var productWeight, totalWeight float64
// 	// inventory := []model.Inventory{}
// 	// for i := 0; i < 100; i++ {
// 	// 	inventory = append(inventory, GenerateDataForInv())
// 	// }

// 	// var twApple, twBanana, twOrange, twMango, twStrawberry, twLettuce, twPear, twGrapes, twSweetPepper float64

// 		switch v.Name {
// 		case prodName:
// 			productWeight = v.TotalWeight + productWeight
// 		}
// 		totalWeight = v.TotalWeight + totalWeight
// 	return productWeight, totalWeight
// }

// switch v.Name {
// case "Apple":
// 	twApple = v.TotalWeight + twApple
// case "Banana":
// 	twBanana = v.TotalWeight + twBanana
// case "Mango":
// 	twMango = v.TotalWeight + twMango
// case "Strawberry":
// 	twStrawberry = v.TotalWeight + twStrawberry
// case "Lettuce":
// 	twLettuce = v.TotalWeight + twLettuce
// case "Pear":
// 	twPear = v.TotalWeight + twPear
// case "Grapes":
// 	twGrapes = v.TotalWeight + twGrapes
// case "Sweet Pepper":
// 	twSweetPepper = v.TotalWeight + twSweetPepper
// }

// func GetDataFromFile() {
// 	data, err := ioutil.ReadFile("mockdata/MOCK_DATA.json")
// 	if err != nil {
// 		err = errors.Wrap(err, "Unable to read the request body")
// 		log.Println(err)
// 	}

// 	inventory := []model.Inventory{}
// 	err = json.NewDecoder(strings.NewReader(string(data))).Decode(&inventory)
// 	// err = json.Unmarshal(data, &inventory)
// 	if err != nil {
// 		err = errors.Wrap(err, "Unable to unmarshal product into Inventory struct")
// 		log.Println(err)
// 	}

// 	//Generate random number

// 	// log.Println(myrand)

// 	// dateArr := time.Now().Unix() - 3600
// 	// expiry := time.Now().AddDate(0, 0, 14).Unix()
// 	// timestamp := time.Now().Unix()
// 	// datesold := time.Now().AddDate(0, 0, myrand).Unix()

// 	for i := 1; i <= 100; i++ {
// 		if i <= 10 {
// 			rand.Seed(time.Now().Unix())
// 			randBtw := random(3, 10)
// 			inventory[i] = ChangesInitialData(ModifyInvData{
// 				Inv:        inventory[i],
// 				Name:       "Apple",
// 				Datearr:    -3600,
// 				Expirydate: 14,
// 				Timestamp:  0,
// 				Randnum:    randBtw,
// 			})
// 		}
// 		if i > 10 && i <= 20 {
// 			rand.Seed(time.Now().Unix())
// 			randBtw := random(3, 10)

// 			inventory[i] = ChangesInitialData(ModifyInvData{
// 				Inv:        inventory[i],
// 				Name:       "Banana",
// 				Datearr:    3600,
// 				Expirydate: 7,
// 				Timestamp:  7200,
// 				Randnum:    randBtw,
// 			})
// 		}
// 		if i > 20 && i <= 30 {
// 			inventory[i].Name = "Orange"
// 			rand.Seed(time.Now().Unix())
// 			randBtw := random(3, 10)

// 			inventory[i] = ChangesInitialData(ModifyInvData{
// 				Inv:        inventory[i],
// 				Name:       "Orange",
// 				Datearr:    7200,
// 				Expirydate: 7,
// 				Timestamp:  7200,
// 				Randnum:    randBtw,
// 			})

// 			log.Println(inventory[i].ItemID)
// 		}
// 		if i > 30 && i <= 40 {
// 			inventory[i].Name = "Mango"
// 		}
// 		if i > 40 && i <= 50 {
// 			inventory[i].Name = "Strawberry"
// 		}
// 		if i > 50 && i <= 60 {
// 			inventory[i].Name = "Tomato"
// 		}
// 		if i > 60 && i <= 70 {
// 			inventory[i].Name = "Lettuce"
// 		}
// 		if i > 70 && i <= 80 {
// 			inventory[i].Name = "Pear"
// 		}
// 		if i > 80 && i <= 90 {
// 			inventory[i].Name = "Grapes"
// 		}
// 		if i > 90 && i <= 100 {
// 			inventory[i].Name = "Sweet Peppers"
// 		}

// 	}
// }

// func ChangesInitialData(modify ModifyInvData) model.Inventory {
// 	inventory.DateArrived = time.Now().Unix() + modify.Datearr
// 	inventory.ExpiryDate = time.Now().AddDate(0, 0, modify.Expirydate).Unix()
// 	inventory.Timestamp = time.Now().Unix() + modify.Timestamp
// 	inventory.DateSold = time.Now().AddDate(0, 0, modify.Randnum).Unix()

// 	return inventory
// }
