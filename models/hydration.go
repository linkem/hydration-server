package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	validator "gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

// Validate Hydration struct
func (h *Hydration) Validate() bool {
	validate = validator.New()
	err := validate.Struct(h)
	isValid := err != nil
	return isValid
}

//Hydration type
type Hydration struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SensorID       string             `bson:"sensorId" json:"sensorId" validation:"required"`
	SensorName     string             `bson:"sensorName" json:"sensorName" validation:"required"`
	Hum            float32            `bson:"hum" json:"hum" validation:"required"`
	Temp           float32            `bson:"temp" json:"temp" validation:"required"`
	Soil           float32            `bson:"soil" json:"soil" validation:"required"`
	CreatedDateUtc time.Time          `bson:"createdDateUtc" json:"createdDateUtc" validation:"required"`
}
