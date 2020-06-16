package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mongoDbTest/models"
	"mongoDbTest/services"
	"net/url"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

//#region mqtt
func connect(clientID string, uri *url.URL) mqtt.Client {
	opts := createClientOptions(clientID, uri)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
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
		fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
		var hydr models.Hydration
		json.Unmarshal(msg.Payload(), &hydr)
		hydr.CreatedDateUtc = time.Now().UTC()
		fmt.Print("Struct: ")
		fmt.Println(hydr)
		hs.CreateHydration(&hydr)
	}
}
func Listen(clientID string, uri *url.URL, topics map[string]mqtt.MessageHandler) {
	client := connect(clientID, uri)
	for topic, handler := range topics {
		client.Subscribe(topic, 0, handler)
	}
}
