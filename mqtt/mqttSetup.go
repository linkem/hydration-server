package mqttServer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mongoDbTest/models"
	"mongoDbTest/services"
	"net/url"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var mqttLog *log.Logger = log.New(os.Stdout, "MQTT ", log.LstdFlags)

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

func createClientOptions(clientID string, uri *url.URL) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))
	// opts.SetUsername(uri.User.Username())
	// password, _ := uri.User.Password()
	// opts.SetPassword(password)
	opts.SetClientID(clientID)
	return opts
}

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

func InitMqtt(c models.Config, hydrationService func() services.Hydrations) {
	if c.Mqtt.Enabled {
		mqttLog.Println("Start MQTT")
		m := make(map[string]mqtt.MessageHandler, 1)
		m[c.Mqtt.HydrationTopic] = HydrationHandler(hydrationService())
		uri, err := url.Parse(c.Mqtt.ConnectionString)
		if err != nil {
			log.Fatal(err)
		}
		go Listen(c.Mqtt.ClientID, uri, m)
	}
}
