package models

type Config struct {
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
