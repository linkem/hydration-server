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
	GetHydrations(ctx context.Context, from time.Time, to time.Time, page int, pageSize int) (*[]models.HydrationGroup, error)
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
func (h *hydrations) GetHydrations(ctx context.Context, from time.Time, to time.Time, page int, pageSize int) (*[]models.HydrationGroup, error) {
	if page <= 0 {
		page = defaultPage
	}
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	pageSizeInt64 := int64(pageSize)
	pageSkipItemsInt64 := (int64(page) - 1) * pageSizeInt64

	// filter := bson.D{
	// 	// {Key: "createdDateUtc", Value: bson.D{
	// 	// 	{Key: "$gte", Value: from},
	// 	// 	{Key: "$lte", Value: to},
	// 	// }},
	// }
	// from = time.Date(2020, 5, 25, 23, 29, 0, 0, time.UTC)
	// to = time.Date(2020, 5, 25, 23, 35, 0, 0, time.UTC)
	ctxWithTimeout, _ := context.WithTimeout(ctx, 30*time.Second)
	//v2
	aggregation := []bson.M{}
	aggregation = append(aggregation, bson.M{"$match": bson.M{"temp": bson.M{"$nin": bson.A{0}}}})
	aggregation = append(aggregation, bson.M{"$match": bson.M{"hum": bson.M{"$nin": bson.A{0}}}})
	aggregation = append(aggregation, bson.M{"$match": bson.M{"soil": bson.M{"$nin": bson.A{0}}}})
	aggregation = append(aggregation, bson.M{"$group": bson.M{
		"_id": bson.M{
			"h":   bson.M{"$hour": "$createdDateUtc"},
			"doy": bson.M{"$dayOfYear": "$createdDateUtc"}},
		"soil":    bson.M{"$avg": "$soil"},
		"temp":    bson.M{"$avg": "$temp"},
		"hum":     bson.M{"$avg": "$hum"},
		"samples": bson.M{"$sum": 1},
		"date":    bson.M{"$min": "$createdDateUtc"},
	}})
	aggregation = append(aggregation, bson.M{"$project": bson.M{
		"soil":    bson.M{"$divide": bson.A{bson.M{"$trunc": bson.M{"$multiply": bson.A{"$soil", 10}}}, 10}},
		"temp":    bson.M{"$divide": bson.A{bson.M{"$trunc": bson.M{"$multiply": bson.A{"$temp", 10}}}, 10}},
		"hum":     bson.M{"$divide": bson.A{bson.M{"$trunc": bson.M{"$multiply": bson.A{"$hum", 10}}}, 10}},
		"samples": 1,
		"createdDateUtc": bson.M{
			"$dateToString": bson.M{
				"format": "%Y-%m-%dT%H:00:00.000Z",
				"date":   "$date"}}}})

	aggregation = append(aggregation, bson.M{"$sort": bson.D{{Key: "_id.doy", Value: -1}, {Key: "_id.h", Value: -1}}})
	aggregation = append(aggregation, bson.M{"$limit": pageSizeInt64})
	aggregation = append(aggregation, bson.M{"$skip": pageSkipItemsInt64})
	cur, err := h.hydrationCollection.Aggregate(ctxWithTimeout, aggregation)

	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctxWithTimeout)
	hydrationGroup := make([]models.HydrationGroup, 0, pageSize)
	cur.Next(ctxWithTimeout)
	h.l.Println(cur.Current)
	if err := cur.All(ctxWithTimeout, &hydrationGroup); err != nil {
		h.l.Panic(err)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	h.l.Println(hydrationGroup)
	return &hydrationGroup, nil
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
