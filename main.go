package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bhushan-aruto/smart_parking_mqtt_message_processor/cache"
	"github.com/bhushan-aruto/smart_parking_mqtt_message_processor/db"
	"github.com/bhushan-aruto/smart_parking_mqtt_message_processor/handler"
	"github.com/bhushan-aruto/smart_parking_mqtt_message_processor/repository"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load the env variables from .env file\n")
	}

	db, err := db.Connect()

	if err != nil {
		log.Fatalf("error while connecting to database, Error -> %v\n", err.Error())
	}

	log.Printf("connected to database\n")

	cache, err := cache.Connect()

	if err != nil {
		log.Fatalf("error while connecting to redis, Error -> %v\n", err.Error())
	}

	log.Printf("connected to redis\n")

	postgresRepo := repository.NewPostgresRepo(db)

	redisRepo := repository.NewRedisRepository(cache)

	brokerHost := os.Getenv("BROKER_HOST")

	if brokerHost == "" {
		log.Fatalf("missing or empty env BROKER_HOST\n")
	}

	brokerPort := os.Getenv("BROKER_PORT")

	if brokerPort == "" {
		log.Fatalf("missing or empty env BROKER_PORT\n")
	}

	clientId := os.Getenv("CLIENT_ID")

	if clientId == "" {
		log.Fatalf("missing or empty env CLIENT_ID")
	}

	brokerAddress := fmt.Sprintf("tcp://%v:%v", brokerHost, brokerPort)

	fmt.Println(brokerAddress)

	opts := mqtt.NewClientOptions()

	opts.AddBroker(brokerAddress)

	opts.SetClientID(clientId)

	opts.OnConnect = func(c mqtt.Client) {
		log.Printf("connected to broker")
	}

	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		log.Printf("disconnected from broker, Error -> %v\n", err.Error())
	}

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Printf("error while connecting to broker, Error -> %v\n", err.Error())
		return
	}

	if client.IsConnected() {
		client.Subscribe("smart/mqtt/processor/1", 1, handler.MessageProcessHandler(redisRepo))
		client.Subscribe("smart/mqtt/processor/2", 1, handler.Gate1ControlHandler(redisRepo))
		client.Subscribe("smart/mqtt/processor/3", 1, handler.Gate2ControlHandler(redisRepo))
		client.Subscribe("smart/mqtt/processor/4", 1, handler.OpenBookedGateHandler(postgresRepo))
	}

	for {
		if !client.IsConnected() {
			if token := client.Connect(); token.Wait() && token.Error() != nil {
				log.Printf("error occurred while connecting to broker, Error -> %v\n", err.Error())
				continue
			}

			client.Unsubscribe("smart/mqtt/processor/1")
			client.Unsubscribe("smart/mqtt/processor/2")
			client.Unsubscribe("smart/mqtt/processor/3")
			client.Unsubscribe("smart/mqtt/processor/4")

			client.Subscribe("smart/mqtt/processor/1", 1, handler.MessageProcessHandler(redisRepo))
			client.Subscribe("smart/mqtt/processor/2", 1, handler.Gate1ControlHandler(redisRepo))
			client.Subscribe("smart/mqtt/processor/3", 1, handler.Gate2ControlHandler(redisRepo))
			client.Subscribe("smart/mqtt/processor/4", 1, handler.OpenBookedGateHandler(postgresRepo))
			log.Printf("reconnected to broker")
		}

		time.Sleep(time.Second * 1)
	}

}
