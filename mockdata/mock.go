package mockdata

import (
	"time"

	"github.com/bhupeshbhatia/go-agg-inventory-v2/model"
)

func JsonForGetJSONString() string {
	var FoodProduct = `{
		"fruit_id": 1,
		"name": "Granny Smith Apples",
		"origin":"ON, Canada",
		"sale_price": 1.12,
		"original_weight": 700,
		"device_id": 1111
	  }`

	return FoodProduct
}

// mockInventory := model.Inventory{
// 	FruitID:      1,
// 	RsCustomerID: "2",
// 	// Name:         "Test",
// 	Origin:           "ON, Canada",
// 	DateArrived:      time.Now(),
// 	DateSold:         time.Now().Add(2),
// 	DeviceID:         1111,
// 	SalePrice:        3.00,
// 	OriginalWeight:   1.00,
// 	SalesWeight:      0.75,
// 	WasteWeight:      0,
// 	DonateWeight:     0,
// 	AggregateVersion: 8,
// 	AggregateID:      1,
// }

func JsonForProductBeginning() []string {
	const data = map[string]interface{}{
		"item_id": 1,
		"name": "This is bulllllllllll",
		"origin": "PEI, Canada",
		"device_id": 233332,
		"date_arrived": 15374686249,
		"expiry_date": 15374686252,
		"price": 106.00,
		"total_weight": 1000.00,
		"location": "A1",
	},
	map[string]interface{}{
		"item_id": 1,
		"name": "This is bulllllllllll",
		"origin": "PEI, Canada",
		"device_id": 233332,
		"date_arrived": 15374686249,
		"expiry_date": 15374686252,
		"price": 106.00,
		"total_weight": 1000.00,
		"location": "A1"
	},
	map[string]interface{}{
		"item_id": 1,
		"name": "This is bulllllllllll",
		"origin": "PEI, Canada",
		"device_id": 233332,
		"date_arrived": 15374686249,
		"expiry_date": 15374686252,
		"price": 106.00,
		"total_weight": 1000.00,
		"location": "A1",
	}]
	return data
}

func JsonForAddProduct() string {
	var FoodProduct = `{
		"fruit_id": 1,
		"rs_customer_id": "1",
		"origin": "ON, Canada",
		"device_id": 1111,
		"date_arrived": "2018-09-13T00:32:23.534Z",
		"sale_price": 3.00,
		"original_weight": 1.00,
		"sales_weight": 0.75,
		"waste_weight": 0,
		"donate_weight": 0,
		"aggregate_version": 8,
		"aggregate_id": 1
		}`
	return FoodProduct
}

func JsonAddWithoutID() string {
	var FoodProduct = `{
		"fruit_id": 0,
		"rs_customer_id": "",
		"origin": "",
		"device_id": 1111,
		"sale_price": 3.00,
		"original_weight": 1.00
		}`
	return FoodProduct
}

func JsonForUpdateProduct() string {
	var FoodProduct = `{
		"fruit_id": 1,
		"rs_customer_id": "1",
		"origin": "ON, Canada",
		"device_id": 1111,
		"sale_price": 3000.00,
		"original_weight": 10.00
		}`
	return FoodProduct
	// "original_weight": 1000.00,

}

func JsonEmptyUpdateProduct() string {
	var FoodProduct = `{
		"fruit_id": 0,
		"rs_customer_id": "",
		"origin": "",
		"device_id": 0,
		"sale_price": 0,
		"original_weight": 0
		}`
	return FoodProduct
	// "original_weight": 1000.00,

}

func JsonDeleteProduct() string {
	var FoodProduct = `{
		"fruit_id": 1,
		"origin": "ON, Canada",
		"device_id": 1111,
		"sale_price": 3.00,
		"original_weight": 1.00
		}`
	return FoodProduct
	// "original_weight": 1000.00,

}

func JsonDelWithoutFruitID() string {
	var FoodProduct = `{
		"fruit_id": 0,
		"origin": "ON, Canada",
		"device_id": 1111,
		"sale_price": 3.00,
		"original_weight": 1.00
		}`
	return FoodProduct
	// "original_weight": 1000.00,

}

func InventoryMock() *model.Inventory {
	mock := &model.Inventory{
		FruitID:      1,
		RsCustomerID: "2",
		// Name:         "Test",
		Origin:      "ON, Canada",
		DateArrived: model.TestTime{time.Now()},
		// DateSold:         time.Now().Add(2),
		DeviceID:         1111,
		SalePrice:        3.00,
		OriginalWeight:   1.00,
		SalesWeight:      0.75,
		WasteWeight:      0,
		DonateWeight:     0,
		AggregateVersion: 8,
		AggregateID:      1,
	}
	return mock
}
