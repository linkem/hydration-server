package models

import "time"

type Id struct {
	H   int32 `json:"h" bson:"h"`
	DoY int32 `json:"doy" bson:"doy"`
}

//Hydration type
type HydrationGroup struct {
	// ID             Id        `json:"_id,omitempty" bson:"_id,omitempty"`
	Soil           float64   `bson:"soil" json:"soil" validation:"required"`
	Hum            float64   `bson:"hum" json:"hum" validation:"required"`
	Temp           float64   `bson:"temp" json:"temp" validation:"required"`
	Samples        int32     `bson:"samples" json:"samples"`
	CreatedDateUtc time.Time `bson:"createdDateUtc" json:"createdDateUtc"`
}
