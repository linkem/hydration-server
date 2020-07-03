package services

import (
	"context"
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
	defaultPageSize int = 50
)

//Hydrations service interface
type Hydrations interface {
	GetHydrations(ctx context.Context, from time.Time, to time.Time, page int, pageSize int) (*[]models.Hydration, error)
	CreateHydration(ctx context.Context, hydration *models.Hydration) *models.Hydration
}

//hydrations service
type hydrations struct {
	hydrationCollection *mongo.Collection
	l                   *log.Logger
}

//NewHydration return new Hydrations object
func NewHydration(logger *log.Logger, databaseAddress string) Hydrations {
	h := &hydrations{l: logger}
	db := initMongoDb(databaseAddress, logger)
	collectionName := "hydration"
	h.hydrationCollection = db.Collection(collectionName)
	return h
}

//GetHydrations get all hydrations
func (h *hydrations) GetHydrations(ctx context.Context, from time.Time, to time.Time, page int, pageSize int) (*[]models.Hydration, error) {
	if page <= 0 {
		page = defaultPage
	}
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	hydration := make([]models.Hydration, 0, pageSize)
	pageSizeInt64 := int64(pageSize)
	pageSkipInt64 := int64(page) - 1
	opt := options.FindOptions{
		Limit: &pageSizeInt64,
		Skip:  &pageSkipInt64,
		Sort:  bson.D{{Key: "createdDateUtc", Value: -1}},
	}
	filter := bson.D{
		// {Key: "createdDateUtc", Value: bson.D{
		// 	{Key: "$gte", Value: from},
		// 	{Key: "$lte", Value: to},
		// }},
	}
	// from = time.Date(2020, 5, 25, 23, 29, 0, 0, time.UTC)
	// to = time.Date(2020, 5, 25, 23, 35, 0, 0, time.UTC)
	ctxWithTimeout, _ := context.WithTimeout(ctx, 30*time.Second)
	cur, err := h.hydrationCollection.Find(ctxWithTimeout, filter, &opt)
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctxWithTimeout)
	cur.All(ctxWithTimeout, &hydration)
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	return &hydration, nil
}

//CreateHydration creates Hydration and insert into database
func (h *hydrations) CreateHydration(ctx context.Context, model *models.Hydration) *models.Hydration {
	h.l.Printf("[Start] Insert")
	ctxWithTimeout, _ := context.WithTimeout(ctx, 30*time.Second)
	result, err := h.hydrationCollection.InsertOne(ctxWithTimeout, model)
	if err != nil {
		log.Fatal(err.Error())
	}
	model.ID = result.InsertedID.(primitive.ObjectID)
	h.l.Printf("Inserted. Id: %s", model.ID.String())
	h.l.Printf("[Finish] Insert")
	return model
}

func initMongoDb(databaseAddress string, l *log.Logger) *mongo.Database {
	// Set client options
	clientOptions := options.Client().ApplyURI(databaseAddress)
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		l.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		l.Fatal(err)
	}
	l.Print("Connected to MongoDB!")
	return client.Database("smart-home")
}
