package main

import (
	"context"
	"fmt"
	"log"
	httpServer "mongoDbTest/http"
	config "mongoDbTest/models"
	mqttServer "mongoDbTest/mqtt"
	"mongoDbTest/services"
	"os"
	"os/signal"
	"strconv"
)

var c config.Config = config.Config{
	MongoDb: struct {
		ConnectionString    string `json:"connectionString"`
		ServerName          string `json:"serverName"`
		HydrationCollection string `json:"hydration"`
	}{
		ConnectionString: "mongodb://mlinke:Ny1y5SuGBeNAh8duSBz3x8HSMlVrSvTFcjLXjMuMMV8jxN4G2bMDZfCH0SGpZuN67OZ3AQjXs4CRWISrxqUcKw==@mlinke.mongo.cosmos.azure.com:10255/?ssl=true&replicaSet=globaldb&retrywrites=false&maxIdleTimeMS=120000&appName=@mlinke@",
		// ConnectionString: "mongodb://192.168.1.100:27017",
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
var progLog *log.Logger = log.New(os.Stdout, "App ", log.LstdFlags)

func initConfig() {
	// MONGO
	dbConnectionString := os.Getenv("MONGO__CONNECTIONSTRING")
	if dbConnectionString != "" {
		c.MongoDb.ConnectionString = dbConnectionString
	}

	//MQTT
	mqttEnabled := os.Getenv("MQTT__ENABLED")
	if mqttEnabled != "" {
		mqttEnabledBool, err := strconv.ParseBool(mqttEnabled)
		if err != nil {
			panic(fmt.Sprint("MQTT__ENABLED - variable had error;" + err.Error()))
		}
		c.Mqtt.Enabled = mqttEnabledBool
	}
	if c.Mqtt.Enabled {
		mqttConnectionString := os.Getenv("MQTT__CONNECTIONSTRING")
		if mqttConnectionString != "" {
			c.Mqtt.ConnectionString = mqttConnectionString
		}
	}

	//SERVER
	httpServerEnabled := os.Getenv("HTTPSERVER__ENABLED")
	if httpServerEnabled != "" {
		httpServerEnabledBool, err := strconv.ParseBool(httpServerEnabled)
		if err != nil {
			panic(fmt.Sprint("HTTPSERVER__ENABLED - variable had error;", err))
		}
		c.Server.Enabled = httpServerEnabledBool
	}
}
func main() {
	initConfig()
	progLog.Println("MQTT ENV value: " + os.Getenv("MQTT__ENABLED"))
	progLog.Println("MQTT Conf value: " + strconv.FormatBool(c.Mqtt.Enabled))

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
	mqttServer.InitMqtt(c, hydrationService)
	httpServer.InitHTTPServer(c, hydrationService)
	handleCancelation(ctx, cancel)
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
