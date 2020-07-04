package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mongoDbTest/handlers"
	"mongoDbTest/middleware"
	"mongoDbTest/models"
	"mongoDbTest/services"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/mux"
)

type config struct {
	MongoDb struct {
		ConnectionString    string `json:"connectionString"`
		ServerName          string `json:"serverName"`
		HydrationCollection string `json:"hydration"`
	} `json:"mongoDb"`
	Mqtt struct {
		Enabled          bool   `json:"enabled"`
		ClientID         string `json:"clientId"`
		ConnectionString string `json:"connectionString"`
		HydrationTopic   string `json:"hydrationTopic"`
	} `json:"mqtt"`
	Server struct {
		Enabled bool   `json:"enabled"`
		Port    string `json:"port"`
	}
}

var c config = config{
	MongoDb: struct {
		ConnectionString    string `json:"connectionString"`
		ServerName          string `json:"serverName"`
		HydrationCollection string `json:"hydration"`
	}{
		ConnectionString: "mongodb://mlinke:Ny1y5SuGBeNAh8duSBz3x8HSMlVrSvTFcjLXjMuMMV8jxN4G2bMDZfCH0SGpZuN67OZ3AQjXs4CRWISrxqUcKw==@mlinke.mongo.cosmos.azure.com:10255/?ssl=true&replicaSet=globaldb&retrywrites=false&maxIdleTimeMS=120000&appName=@mlinke@",
	},
	Mqtt: struct {
		Enabled          bool   `json:"enabled"`
		ClientID         string `json:"clientId"`
		ConnectionString string `json:"connectionString"`
		HydrationTopic   string `json:"hydrationTopic"`
	}{
		Enabled:          false,
		ClientID:         "raspberry-serv-1",
		ConnectionString: "mqtt://192.168.1.100:1883",
		HydrationTopic:   "home/hydration/livingroom/avocado",
	},
	Server: struct {
		Enabled bool   `json:"enabled"`
		Port    string `json:"port"`
	}{
		Enabled: true,
		Port:    "8080",
	},
}
var progLog *log.Logger

//#endregion mqtt
func main() {
	progLog = log.New(os.Stdout, "App ", log.LstdFlags)
	progLog.Println("Create Context")
	// trap Ctrl+C and call cancel on the context
	backgroundCtx := context.Background()
	ctx, cancel := context.WithCancel(backgroundCtx)

	progLog.Println("Start program")
	hydrationService := func() services.Hydrations {
		progLog.Println("Creating hydration")

		return services.NewHydration(
			log.New(os.Stdout, "hydration-service ", log.LstdFlags),
			c.MongoDb.ConnectionString)
	}
	initMqtt(c, hydrationService)
	initHTTPServer(c, hydrationService)
	handleCancelation(ctx, cancel)
}

func initMqtt(c config, hydrationService func() services.Hydrations) {
	if c.Mqtt.Enabled {
		progLog.Println("Start MQTT")
		m := make(map[string]mqtt.MessageHandler, 1)
		m[c.Mqtt.HydrationTopic] = HydrationHandler(hydrationService())
		uri, err := url.Parse(c.Mqtt.ConnectionString)
		if err != nil {
			log.Fatal(err)
		}
		go Listen(c.Mqtt.ClientID, uri, m)
	}
}
func initHTTPServer(c config, hydrationService func() services.Hydrations) {
	if c.Server.Enabled {
		progLog.Printf("Start HTTP at port %s\n", c.Server.Port)
		router := routes(hydrationService())
		go func() {
			log.Fatal(http.ListenAndServe(":"+c.Server.Port, router))
		}()
	}
}
func handleCancelation(ctx context.Context, cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		progLog.Println("In Defered Funs")
		signal.Stop(c)
		cancel()
	}()
	select {
	case <-c:
		progLog.Println("in C channel")
		cancel()
	case <-ctx.Done():
		progLog.Println("in Done channel")
	}
	progLog.Println("Finish")
}

func routes(hydrationService services.Hydrations) *mux.Router {
	router := mux.NewRouter()
	controller := handlers.NewHydrationController(
		log.New(os.Stdout, "hydration-api ", log.LstdFlags),
		hydrationService)

	router.Use(middleware.LoggingMiddleware, middleware.HeadersMiddleware)

	router.HandleFunc("/time", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, time.Now().UTC())
	})
	router.HandleFunc("/hydration", controller.GetHydrations()).Methods("GET")
	return router
}

//#region mqtt
var mqttLog *log.Logger

func connect(clientID string, uri *url.URL) mqtt.Client {
	mqttLog.Printf("Connecting MQTT to %s", uri.String())
	opts := createClientOptions(clientID, uri)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	// for !token.WaitTimeout(3 * time.Second) {
	// 	mqttLog.Println("MQTT Waiting for connection")
	// }
	for !token.Wait() {
		mqttLog.Println("MQTT Waiting for connection")
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	mqttLog.Printf("MQTT Connected to %s", uri.String())
	return client
}

func createClientOptions(clientID string, uri *url.URL) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))
	// opts.SetUsername(uri.User.Username())
	// password, _ := uri.User.Password()
	// opts.SetPassword(password)
	opts.SetClientID(clientID)
	return opts
}

func HydrationHandler(hs services.Hydrations) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		mqttLog.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))

		var hydr models.Hydration
		json.Unmarshal(msg.Payload(), &hydr)
		hydr.CreatedDateUtc = time.Now().UTC()
		mqttLog.Print("Struct: ")
		mqttLog.Println(hydr)
		hs.CreateHydration(context.Background(), &hydr)
	}
}
func Listen(clientID string, uri *url.URL, topics map[string]mqtt.MessageHandler) {
	mqttLog = log.New(os.Stdout, "MQTT ", log.LstdFlags)
	client := connect(clientID, uri)
	mqttLog.Printf("Connected MQTT to %s\n", uri.String())
	for topic, handler := range topics {
		mqttLog.Printf("Subscribing to topic: %s", topic)
		client.Subscribe(topic, 0, handler)
		mqttLog.Printf("Subscribed to topic: %s", topic)
	}
}
