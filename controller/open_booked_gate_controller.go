package controller

import (
	"encoding/json"
	"log"

	"github.com/bhushan-aruto/smart_parking_mqtt_message_processor/model"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func OpenBookedGateController(c mqtt.Client, m mqtt.Message, dbRepo model.DatabaseRepository) {
	var mqttMessage model.OpenBookedGateRequestMessage

	if err := json.Unmarshal(m.Payload(), &mqttMessage); err != nil {
		log.Printf("error occured while deserialzing json, err -> %v\n", err.Error())
		return
	}

	userId, err := dbRepo.GetUserIdByEmail(mqttMessage.Email)

	if err != nil {
		log.Printf("error occurred with database while getting user id , err -> %v\n", err.Error())
		return
	}

	bookingExists, err := dbRepo.CheckUserBookingExists(userId)

	if err != nil {
		log.Printf("error occurred with database while checking booking exists, err -> %v\n", err.Error())
		return
	}

	if !bookingExists {
		return
	}

	if err := dbRepo.DeleteUserBooking(userId); err != nil {
		log.Printf("error occured while deleting the user booking, err -> %v", err.Error())
		return
	}

	response := &model.OpenBookedGateResponse{
		MessageType: 4,
	}

	responseJsonMessage, err := json.Marshal(response)

	if err != nil {
		log.Printf("error occured while serializing to json, err -> %v\n", err.Error())
		return
	}

	c.Publish("smart/controller/1", 1, false, responseJsonMessage)
}
