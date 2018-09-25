package model

import (
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
	SoldWeight       float64           `bson:"sale_weight,omitempty" json:"sale_weight,omitempty"`
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
	SoldWeight       float64           `bson:"sale_weight,omitempty" json:"sale_weight,omitempty"`
}

func (i *Inventory) MarshalBSON() ([]byte, error) {
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

func (i *Inventory) UnmarshalJSON(in []byte) error {
	m := make(map[string]interface{})
	err := bson.Unmarshal(in, m)
	if err != nil {
		err = errors.Wrap(err, "Unmarshal Error")
		return err
	}

	i.ID = m["_id"].(objectid.ObjectID)

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
	i.Name = m["name"].(string)
	i.Origin = m["origin"].(string)
	i.TotalWeight = m["total_weight"].(float64)
	i.Price = m["price"].(float64)
	i.Location = m["location"].(string)
	i.DateArrived = m["date_arrived"].(int64)
	i.ExpiryDate = m["expiry_date"].(int64)
	i.Timestamp = m["timestamp"].(int64)
	i.WasteWeight = m["waste_weight"].(float64)
	i.DonateWeight = m["donate_weight"].(float64)
	i.AggregateVersion = m["aggregate_version"].(int64)
	i.AggregateID = m["aggregate_id"].(int8)
	i.DateSold = m["date_sold"].(int64)
	i.SalePrice = m["sale_price"].(float64)
	i.SoldWeight = m["sold_weight"].(float64)

	return nil
}
