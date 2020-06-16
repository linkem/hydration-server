package services

import (
	"context"
	"fmt"
	"log"
	"mongoDbTest/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	defaultPage     int = 1
	defaultPageSize int = 10
)

//Hydrations service interface
type Hydrations interface {
	GetHydrations(from time.Time, to time.Time, page int, pageSize int) (*[]models.Hydration, error)
	CreateHydration(*models.Hydration) *models.Hydration
}

//hydrations service
type hydrations struct {
	ctx                 context.Context
	hydrationCollection *mongo.Collection
	l                   *log.Logger
}

//NewHydration return new Hydrations object
func NewHydration(ctx context.Context, logger *log.Logger, databaseAddress string) Hydrations {
	h := &hydrations{l: logger}
	db := initMongoDb(databaseAddress)
	collectionName := "hydration"
	ctxWithTimeout, _ := context.WithTimeout(h.ctx, 30*time.Second)

	h.ctx = ctxWithTimeout
	h.hydrationCollection = db.Collection(collectionName)
	return h
}

//GetHydrations get all hydrations
func (h *hydrations) GetHydrations(from time.Time, to time.Time, page int, pageSize int) (*[]models.Hydration, error) {
	if page <= 0 {
		page = defaultPage
	}
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	hydration := make([]models.Hydration, 0, pageSize)
	opt := options.FindOptions{
		// Limit: &pageSize,
		Sort: bson.D{{Key: "createdDateUtc", Value: -1}},
	}

	// from = time.Date(2020, 5, 25, 23, 29, 0, 0, time.UTC)
	// to = time.Date(2020, 5, 25, 23, 35, 0, 0, time.UTC)
	cur, err := h.hydrationCollection.Find(h.ctx, bson.D{
		{Key: "createdDateUtc", Value: bson.D{
			{Key: "$gte", Value: from},
			{Key: "$lte", Value: to},
		}},
	}, &opt)
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(h.ctx)
	cur.All(h.ctx, &hydration)
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	return &hydration, nil
}

//CreateHydration creates Hydration and insert into database
func (h *hydrations) CreateHydration(model *models.Hydration) *models.Hydration {
	h.l.Printf("[Start] Insert")
	result, err := h.hydrationCollection.InsertOne(h.ctx, model)
	if err != nil {
		log.Fatal(err.Error())
	}
	model.ID = result.InsertedID.(primitive.ObjectID)
	h.l.Printf("Inserted. Id: %s", model.ID.String())
	h.l.Printf("[Finish] Insert")
	return model
}

func initMongoDb(databaseAddress string) *mongo.Database {
	// Set client options
	clientOptions := options.Client().ApplyURI(databaseAddress)
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	return client.Database("smart-home")
}
