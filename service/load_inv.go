package service

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/bhupeshbhatia/go-agg-inventory-v2/model"
	"github.com/pkg/errors"
)

func GetDataFromFile() {
	data, err := ioutil.ReadFile("mockdata/MOCK_DATA.json")
	if err != nil {
		err = errors.Wrap(err, "Unable to read the request body")
		log.Println(err)
	}

	inventory := []model.Inventory{}
	err = json.Unmarshal(data, &inventory)
	if err != nil {
		err = errors.Wrap(err, "Unable to unmarshal product into Inventory struct")
		log.Println(err)
	}

	for i := 1; i <= 100; i++ {
		if i <= 10 {
			inventory[i].Name = "Apple"
		}
		if i > 10 && i <= 20 {
			inventory[i].Name = "Banana"
		}
		if i > 20 && i <= 30 {
			inventory[i].Name = "Orange"
		}
		if i > 30 && i <= 40 {
			inventory[i].Name = "Mango"
		}
		if i > 40 && i <= 50 {
			inventory[i].Name = "Strawberry"
		}
		if i > 50 && i <= 60 {
			inventory[i].Name = "Tomato"
		}
		if i > 60 && i <= 70 {
			inventory[i].Name = "Lettuce"
		}
		if i > 70 && i <= 80 {
			inventory[i].Name = "Pear"
		}
		if i > 80 && i <= 90 {
			inventory[i].Name = "Grapes"
		}
		if i > 90 && i <= 100 {
			inventory[i].Name = "Sweet Peppers"
		}

		log.Println(inventory[i].Name)
	}
}
