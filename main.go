package main

import (
	"context"
	"log"
	"mongoDbTest/services"
	"net/http"
	"net/url"
	"os"
	"os/signal"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// func testMqttBasic(client mqtt.Client, msg mqtt.Message) {
// }

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
		Enabled:          true,
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
	ctx, cancel := context.WithCancel(context.Background())

	progLog.Println("Start program")
	progLog.Println("Creating hydration")
	hydrationService := func() services.Hydrations {
		return services.NewHydration(
			ctx,
			log.New(os.Stdout, "hydration-service ", log.LstdFlags),
			c.MongoDb.ConnectionString)
	}
	progLog.Println("Created Hydration")
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
		progLog.Println("Start HTTP")
		router := Routes(hydrationService())
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
