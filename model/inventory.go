package model

import (
	"encoding/json"
	"log"

	"github.com/TerrexTech/uuuid"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/pkg/errors"
)

//Inventory represents inventory collection
type Inventory struct {
	ID               objectid.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ItemID           uuuid.UUID        `bson:"item_id,omitempty" json:"item_id,omitempty"`
	Name             string            `bson:"name,omitempty" json:"name,omitempty"`
	Origin           string            `bson:"origin,omitempty" json:"origin,omitempty"`
	DeviceID         uuuid.UUID        `bson:"device_id,omitempty" json:"device_id,omitempty"`
	TotalWeight      float64           `bson:"total_weight,omitempty" json:"total_weight,omitempty"`
	Price            float64           `bson:"price,omitempty" json:"price,omitempty"`
	Location         string            `bson:"location,omitempty" json:"location,omitempty"`
	DateArrived      int64             `bson:"date_arrived,omitempty" json:"date_arrived,omitempty"`
	ExpiryDate       int64             `bson:"expiry_date,omitempty" json:"expiry_date,omitempty"`
	Timestamp        int64             `bson:"timestamp,omitempty" json:"timestamp,omitempty"`
	RsCustomerID     uuuid.UUID        `bson:"rs_customer_id,omitempty" json:"rs_customer_id,omitempty"`
	WasteWeight      float64           `bson:"waste_weight,omitempty" json:"waste_weight,omitempty"`
	DonateWeight     float64           `bson:"donate_weight,omitempty" json:"donate_weight,omitempty"`
	AggregateVersion int64             `bson:"aggregate_version,omitempty" json:"aggregate_version,omitempty"`
	AggregateID      int8              `bson:"aggregate_id,omitempty" json:"aggregate_id,omitempty"`
	DateSold         int64             `bson:"date_sold,omitempty" json:"date_sold,omitempty"`
	SalePrice        float64           `bson:"sale_price,omitempty" json:"sale_price,omitempty"`
	SoldWeight       float64           `bson:"sold_weight,omitempty" json:"sold_weight,omitempty"`
}

type marshalInventory struct {
	ID               objectid.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ItemID           string            `bson:"item_id,omitempty" json:"item_id,omitempty"`
	Name             string            `bson:"name,omitempty" json:"name,omitempty"`
	Origin           string            `bson:"origin,omitempty" json:"origin,omitempty"`
	DeviceID         string            `bson:"device_id,omitempty" json:"device_id,omitempty"`
	TotalWeight      float64           `bson:"total_weight,omitempty" json:"total_weight,omitempty"`
	Price            float64           `bson:"price,omitempty" json:"price,omitempty"`
	Location         string            `bson:"location,omitempty" json:"location,omitempty"`
	DateArrived      int64             `bson:"date_arrived,omitempty" json:"date_arrived,omitempty"`
	ExpiryDate       int64             `bson:"expiry_date,omitempty" json:"expiry_date,omitempty"`
	Timestamp        int64             `bson:"timestamp,omitempty" json:"timestamp,omitempty"`
	RsCustomerID     string            `bson:"rs_customer_id,omitempty" json:"rs_customer_id,omitempty"`
	WasteWeight      float64           `bson:"waste_weight,omitempty" json:"waste_weight,omitempty"`
	DonateWeight     float64           `bson:"donate_weight,omitempty" json:"donate_weight,omitempty"`
	AggregateVersion int64             `bson:"aggregate_version,omitempty" json:"aggregate_version,omitempty"`
	AggregateID      int8              `bson:"aggregate_id,omitempty" json:"aggregate_id,omitempty"`
	DateSold         int64             `bson:"date_sold,omitempty" json:"date_sold,omitempty"`
	SalePrice        float64           `bson:"sale_price,omitempty" json:"sale_price,omitempty"`
	SoldWeight       float64           `bson:"sold_weight,omitempty" json:"sold_weight,omitempty"`
}

func (i *Inventory) MarshalJSON() ([]byte, error) {
	log.Println("}}}}}}}}}}}}}}}}}}}}}}}}}}")
	log.Println(i)
	in := &marshalInventory{
		Name:             i.Name,
		Origin:           i.Origin,
		TotalWeight:      i.TotalWeight,
		Price:            i.Price,
		Location:         i.Location,
		DateArrived:      i.DateArrived,
		ExpiryDate:       i.ExpiryDate,
		Timestamp:        i.Timestamp,
		WasteWeight:      i.WasteWeight,
		DonateWeight:     i.DonateWeight,
		AggregateVersion: i.AggregateVersion,
		AggregateID:      i.AggregateID,
		DateSold:         i.DateSold,
		SalePrice:        i.SalePrice,
		SoldWeight:       i.SoldWeight,
	}

	if i.ItemID.String() != (uuuid.UUID{}).String() {
		in.ItemID = i.ItemID.String()
	}
	if i.DeviceID.String() != (uuuid.UUID{}).String() {
		in.DeviceID = i.DeviceID.String()
	}
	if i.RsCustomerID.String() != (uuuid.UUID{}).String() {
		in.RsCustomerID = i.RsCustomerID.String()
	}
	e, _ := json.Marshal(in)
	log.Println(string(e))

	return json.Marshal(in)
}

func (i Inventory) MarshalBSON() ([]byte, error) {
	log.Println("222222222222222222222222222222")
	in := &marshalInventory{
		Name:             i.Name,
		Origin:           i.Origin,
		TotalWeight:      i.TotalWeight,
		Price:            i.Price,
		Location:         i.Location,
		DateArrived:      i.DateArrived,
		ExpiryDate:       i.ExpiryDate,
		Timestamp:        i.Timestamp,
		WasteWeight:      i.WasteWeight,
		DonateWeight:     i.DonateWeight,
		AggregateVersion: i.AggregateVersion,
		AggregateID:      i.AggregateID,
		DateSold:         i.DateSold,
		SalePrice:        i.SalePrice,
		SoldWeight:       i.SoldWeight,
	}

	if i.ItemID.String() != (uuuid.UUID{}).String() {
		in.ItemID = i.ItemID.String()
	}
	if i.DeviceID.String() != (uuuid.UUID{}).String() {
		in.DeviceID = i.DeviceID.String()
	}
	if i.RsCustomerID.String() != (uuuid.UUID{}).String() {
		in.RsCustomerID = i.RsCustomerID.String()
	}

	return bson.Marshal(in)
}

func (i *Inventory) UnmarshalBSON(in []byte) error {
	m := make(map[string]interface{})
	err := bson.Unmarshal(in, m)
	if err != nil {
		err = errors.Wrap(err, "Unmarshal Error")
		return err
	}

	if m["_id"] != nil {
		i.ID = m["_id"].(objectid.ObjectID)
	}

	i.ItemID, err = uuuid.FromString(m["item_id"].(string))
	if err != nil {
		err = errors.Wrap(err, "Error parsing ItemID for inventory")
		return err
	}

	i.DeviceID, err = uuuid.FromString(m["device_id"].(string))
	if err != nil {
		err = errors.Wrap(err, "Error parsing DeviceID for inventory")
		return err
	}

	i.RsCustomerID, err = uuuid.FromString(m["rs_customer_id"].(string))
	if err != nil {
		err = errors.Wrap(err, "Error parsing DeviceID for inventory")
		return err
	}

	if m["name"] != nil {
		i.Name = m["name"].(string)
	}

	if m["origin"] != nil {
		i.Origin = m["origin"].(string)
	}

	if m["total_weight"] != nil {
		i.TotalWeight = m["total_weight"].(float64)
	}

	if m["price"] != nil {
		i.Price = m["price"].(float64)
	}

	if m["location"] != nil {
		i.Location = m["location"].(string)
	}

	if m["date_arrived"] != nil {
		i.DateArrived = m["date_arrived"].(int64)
	}

	if m["expiry_date"] != nil {
		i.ExpiryDate = m["expiry_date"].(int64)
	}

	if m["timestamp"] != nil {
		i.Timestamp = m["timestamp"].(int64)
	}

	if m["waste_weight"] != nil {
		i.WasteWeight = m["waste_weight"].(float64)
	}

	if m["donate_weight"] != nil {
		i.DonateWeight = m["donate_weight"].(float64)
	}

	if m["aggregate_version"] != nil {
		i.AggregateVersion = m["aggregate_version"].(int64)
	}

	if m["aggregate_id"] != nil {
		// i.AggregateID = m["aggregate_id"].(int8)
	}

	if m["date_sold"] != nil {
		i.DateSold = m["date_sold"].(int64)
	}

	if m["sale_price"] != nil {
		i.SalePrice = m["sale_price"].(float64)
	}

	if m["sold_weight"] != nil {
		i.SoldWeight = m["sold_weight"].(float64)
	}
	return nil
}

func (i *Inventory) UnmarshalJSON(in []byte) error {
	m := make(map[string]interface{})
	err := json.Unmarshal(in, &m)
	if err != nil {
		err = errors.Wrap(err, "Unmarshal Error")
		return err
	}

	if m["_id"] != nil {
		i.ID = m["_id"].(objectid.ObjectID)
	}

	if m["item_id"] != nil {
		i.ItemID, err = uuuid.FromString(m["item_id"].(string))
	}
	if err != nil {
		err = errors.Wrap(err, "Error parsing ItemID for inventory")
		return err
	}

	if m["device_id"] != nil {
		i.DeviceID, err = uuuid.FromString(m["device_id"].(string))
	}
	if err != nil {
		err = errors.Wrap(err, "Error parsing DeviceID for inventory")
		return err
	}

	if m["rs_customer_id"] != nil {
		i.RsCustomerID, err = uuuid.FromString(m["rs_customer_id"].(string))
	}
	if err != nil {
		err = errors.Wrap(err, "Error parsing DeviceID for inventory")
		return err
	}

	if m["name"] != nil {
		i.Name = m["name"].(string)
	}

	if m["origin"] != nil {
		i.Origin = m["origin"].(string)
	}

	if m["total_weight"] != nil {
		i.TotalWeight = m["total_weight"].(float64)
	}

	if m["price"] != nil {
		i.Price = m["price"].(float64)
	}

	if m["location"] != nil {
		i.Location = m["location"].(string)
	}

	if m["date_arrived"] != nil {
		i.DateArrived = m["date_arrived"].(int64)
	}

	if m["expiry_date"] != nil {
		i.ExpiryDate = m["expiry_date"].(int64)
	}

	if m["timestamp"] != nil {
		i.Timestamp = m["timestamp"].(int64)
	}

	if m["waste_weight"] != nil {
		i.WasteWeight = m["waste_weight"].(float64)
	}

	if m["donate_weight"] != nil {
		i.DonateWeight = m["donate_weight"].(float64)
	}

	if m["aggregate_version"] != nil {
		i.AggregateVersion = m["aggregate_version"].(int64)
	}

	if m["aggregate_id"] != nil {
		// i.AggregateID = m["aggregate_id"].(int8)
	}

	if m["date_sold"] != nil {
		i.DateSold = m["date_sold"].(int64)
	}

	if m["sale_price"] != nil {
		i.SalePrice = m["sale_price"].(float64)
	}

	if m["sold_weight"] != nil {
		i.SoldWeight = m["sold_weight"].(float64)
	}

	// var ok bool
	// i.Name = m["name"].(string)
	// i.Origin = m["origin"].(string)
	// i.TotalWeight = m["total_weight"].(float64)
	// i.Price, ok = m["price"].(float64)
	// if !ok {
	// 	log.Println("Error converting price to float64")
	// }
	// i.Location = m["location"].(string)
	// i.DateArrived = m["date_arrived"].(int64)
	// i.ExpiryDate = m["expiry_date"].(int64)
	// i.Timestamp = m["timestamp"].(int64)
	// i.WasteWeight = m["waste_weight"].(float64)
	// i.DonateWeight = m["donate_weight"].(float64)
	// if m["aggregate_version"] != nil {
	// 	i.AggregateVersion = m["aggregate_version"].(int64)
	// }
	// if m["aggregate_id"] != nil {
	// 	// i.AggregateID = m["aggregate_id"].(int8)
	// }
	// i.DateSold = m["date_sold"].(int64)
	// i.SalePrice = m["sale_price"].(float64)
	// if m["sold_weight"] != nil {
	// 	i.SoldWeight = m["sold_weight"].(float64)
	// }

	return nil
}
