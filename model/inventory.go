package model

import (
	"github.com/TerrexTech/uuuid"
)

//Inventory represents inventory collection
type Inventory struct {
	ItemID           uuuid.UUID `bson:"item_id,omitempty" json:"item_id,omitempty"`
	Name             string     `bson:"name,omitempty" json:"name,omitempty"`
	Origin           string     `bson:"origin,omitempty" json:"origin,omitempty"`
	DeviceID         int64      `bson:"device_id,omitempty" json:"device_id,omitempty"`
	TotalWeight      float64    `bson:"total_weight,omitempty" json:"total_weight,omitempty"`
	Price            float64    `bson:"price,omitempty" json:"price,omitempty"`
	Location         string     `bson:"location,omitempty" json:"location,omitempty"`
	DateArrived      int64      `bson:"date_arrived,omitempty" json:"date_arrived,omitempty"`
	ExpiryDate       int64      `bson:"expiry_date,omitempty" json:"expiry_date,omitempty"`
	Timestamp        int64      `bson:"timestamp,omitempty" json:"timestamp,omitempty"`
	RsCustomerID     string     `bson:"rs_customer_id,omitempty" json:"rs_customer_id,omitempty"`
	WasteWeight      float64    `bson:"waste_weight,omitempty" json:"waste_weight,omitempty"`
	DonateWeight     float64    `bson:"donate_weight,omitempty" json:"donate_weight,omitempty"`
	AggregateVersion int64      `bson:"aggregate_version,omitempty" json:"aggregate_version,omitempty"`
	AggregateID      int64      `bson:"aggregate_id,omitempty" json:"aggregate_id,omitempty"`
	DateSold         int64      `bson:"date_sold,omitempty" json:"date_sold,omitempty"`
	SalePrice        float64    `bson:"sale_price,omitempty" json:"sale_price,omitempty"`
	OriginalWeight   float64    `bson:"original_weight,omitempty" json:"original_weight,omitempty"`
}

type marshalInventory struct {
	ItemID string `bson:"item_id,omitempty" json:"item_id,omitempty"`
	// RsCustomerID string  `bson:"rs_customer_id,omitempty" json:"rs_customer_id,omitempty"`
	Name        string  `bson:"name,omitempty" json:"name,omitempty"`
	Origin      string  `bson:"origin,omitempty" json:"origin,omitempty"`
	DateArrived string  `bson:"date_arrived,omitempty" json:"date_arrived,omitempty"`
	ExpiryDate  string  `bson:"expiry_date,omitempty" json:"expiry_date,omitempty"`
	DeviceID    int64   `bson:"device_id,omitempty" json:"device_id,omitempty"`
	TotalWeight float64 `bson:"total_weight,omitempty" json:"total_weight,omitempty"`
	Price       float64 `bson:"price,omitempty" json:"price,omitempty"`
	// WasteWeight      float64 `bson:"waste_weight,omitempty" json:"waste_weight,omitempty"`
	// DonateWeight     float64 `bson:"donate_weight,omitempty" json:"donate_weight,omitempty"`
	Location string `bson:"location,omitempty" json:"location,omitempty"`
	// AggregateVersion int64  `bson:"aggregate_version,omitempty" json:"aggregate_version,omitempty"`
	// AggregateID      int64  `bson:"aggregate_id,omitempty" json:"aggregate_id,omitempty"`
	// DateSold         string  `bson:"date_sold,omitempty" json:"date_sold,omitempty"`
	// SalePrice        float64   `bson:"sale_price,omitempty" json:"sale_price,omitempty"`
	// SoldWeight float64 `bson:"sold_weight,omitempty" json:"sold_weight,omitempty"`
	Timestamp string `bson:"timestamp,omitempty" json:"timestamp,omitempty"`
}

// func (i *Inventory) MarshalJSON() ([]byte, error) {
// 	in := &marshalInventory{
// 		Name:        i.Name,
// 		Origin:      i.Origin,
// 		DeviceID:    i.DeviceID,
// 		TotalWeight: i.TotalWeight,
// 		Price:       i.Price,
// 		// WasteWeight:      i.WasteWeight,
// 		// DonateWeight:     i.DonateWeight,
// 		Location: i.Location,
// 		// AggregateVersion: i.AggregateVersion,
// 		// AggregateID:      i.AggregateID,
// 		// SalePrice: i.SalePrice,
// 		// OriginalWeight: i.OriginalWeight,
// 	}

// 	if i.ItemID.String() != (uuid.UUID{}).String() {
// 		in.ItemID = i.ItemID.String()
// 	}

// 	if i.ExpiryDate.String() != "0001-01-01 00:00:00 +0000 UTC" {
// 		in.ExpiryDate = i.ExpiryDate.Format(time.RFC3339Nano)
// 	}
// 	if i.DateArrived.String() != "0001-01-01 00:00:00 +0000 UTC" {
// 		in.DateArrived = i.DateArrived.Format(time.RFC3339Nano)
// 	}
// 	// if i.DateSold.String() != "0001-01-01 00:00:00 +0000 UTC" {
// 	// 	in.DateSold = i.DateSold.Format(time.RFC3339Nano)
// 	// }
// 	if i.Timestamp.String() != "0001-01-01 00:00:00 +0000 UTC" {
// 		in.Timestamp = i.Timestamp.Format(time.RFC3339Nano)
// 	}

// 	return json.Marshal(in)
// }

// func (i *Inventory) UnmarshalJSON(in []byte) error {
// 	m := make(map[string]interface{})
// 	err := json.Unmarshal(in, &m)
// 	if err != nil {
// 		err = errors.Wrap(err, "Unmarshal Error")
// 		return err
// 	}

// 	if m["name"] != nil {
// 		i.Name = m["name"].(string)
// 		if err != nil {
// 			err = errors.Wrap(err, "Unmarshal Error: Error parsing name")
// 			return err
// 		}
// 	}

// 	if m["origin"] != nil {
// 		i.Origin = m["origin"].(string)
// 		if err != nil {
// 			err = errors.Wrap(err, "Unmarshal Error: Error parsing origin")
// 			return err
// 		}
// 	}

// 	if m["device_id"] != nil {
// 		i.DeviceID = int64((m["device_id"]).(float64))
// 		if err != nil {
// 			err = errors.Wrap(err, "Unmarshal Error: Error parsing DeviceID")
// 			return err
// 		}
// 	}

// 	if m["total_weight"] != nil {
// 		i.TotalWeight = m["total_weight"].(float64)
// 		if err != nil {
// 			err = errors.Wrap(err, "Unmarshal Error: Error parsing total_weight")
// 			return err
// 		}
// 	}

// 	if m["price"] != nil {
// 		i.Price = m["price"].(float64)
// 		if err != nil {
// 			err = errors.Wrap(err, "Unmarshal Error: Error parsing price")
// 			return err
// 		}
// 	}

// 	if m["location"] != nil {
// 		i.Location = m["location"].(string)
// 		if err != nil {
// 			err = errors.Wrap(err, "Unmarshal Error: Error parsing location")
// 			return err
// 		}
// 	}

// 	// i.DonateWeight = m["donate_weight"].(float64)

// 	// i.AggregateVersion = int64(m["aggregate_version"].(float64))
// 	// i.AggregateID = int64(m["aggregate_id"].(float64))

// 	log.Println(m["item_id"])

// 	if m["item_id"] != nil {
// 		i.ItemID, err = uuid.FromString(m["item_id"].(string))
// 		if err != nil {
// 			err = errors.Wrap(err, "Unmarshal Error: Error parsing user _id")
// 			return err
// 		}
// 	}

// 	if m["date_arrived"] != nil {
// 		arriveTime, err := time.Parse(time.RFC3339Nano, m["date_arrived"].(string))
// 		i.DateArrived = arriveTime.Unix()
// 		// i.DateArrived, err = m["date_arrived"].(string)
// 		log.Println("date_arrived", m["date_arrived"])
// 		if err != nil {
// 			err = errors.Wrap(err, "Error parsing time while Unmarshalling Inventory")
// 			return err
// 		}
// 	}

// 	// if m["date_sold"] != nil {
// 	// 	i.DateSold, err = time.Parse(time.RFC3339Nano, m["date_sold"].(string))
// 	// 	if err != nil {
// 	// 		err = errors.Wrap(err, "Error parsing time while Unmarshalling Inventory- date_sold")
// 	// 		return err
// 	// 	}
// 	// }

// 	if m["expiry_date"] != nil {
// 		parsedTime, err := time.Parse(time.RFC3339Nano, m["expiry_date"].(string))
// 		i.ExpiryDate = parsedTime.Unix()
// 		if err != nil {
// 			err = errors.Wrap(err, "Error parsing time while Unmarshalling Inventory - expiry_date")
// 			return err
// 		}
// 	}

// 	if m["timestamp"] != nil {
// 		timestrampTime, err := time.Parse(time.RFC3339Nano, m["timestamp"].(string))
// 		i.Timestamp = timestrampTime.Unix()
// 		log.Println("timestamp", m["timestamp"])
// 		if err != nil {
// 			err = errors.Wrap(err, "Error parsing time while Unmarshalling Inventory - expiry_date")
// 			return err
// 		}
// 	}
// 	return nil

// }
