package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"mongoDbTest/models"
	"net/url"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func connect(clientId string, uri *url.URL) mqtt.Client {
	opts := createClientOptions(clientId, uri)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
		progLog.Println("Try to connect")
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	return client
}

func createClientOptions(clientId string, uri *url.URL) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))
	opts.SetUsername(uri.User.Username())
	password, _ := uri.User.Password()
	opts.SetPassword(password)
	opts.SetClientID(clientId)
	return opts
}

func listen(uri *url.URL, topic string) {
	client := connect("sub", uri)
	client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		progLog.Println("Recived Message.")
		// progLog.Println("Subscribe Callback: ")
		// progLog.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
		// progLog.Println("********************")

	})
}

var progLog *log.Logger

func main() {
	progLog = log.New(os.Stdout, "SensorMock ", log.LstdFlags)
	progLog.Println("Start")
	uri, err := url.Parse("mqtt://192.168.1.100:1883")
	if err != nil {
		log.Fatal(err)
	}
	topic := "home/hydration/livingroom/avocado"
	go listen(uri, topic)

	progLog.Println("Connect")
	client := connect("Mock", uri)
	timer := time.NewTicker(1 * time.Second)
	for range timer.C {
		temp := rand.Float32() * 30
		mockSensor := &models.Hydration{
			SensorID:   "mock",
			SensorName: "mock",
			Hum:        float32(int(temp*3+10*100)) / 100,
			Soil:       float32(int(temp*3*100)) / 100,
			Temp:       float32(int(temp*100)) / 100,
		}
		b, _ := json.Marshal(mockSensor)
		mockSensorJSON := string(b)
		progLog.Println("Publish Message")
		progLog.Printf("%s\n", mockSensorJSON)
		client.Publish(topic, 0, false, mockSensorJSON)
	}
}
