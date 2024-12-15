package handler

import (
	"github.com/bhushan-aruto/smart_parking_mqtt_message_processor/controller"
	"github.com/bhushan-aruto/smart_parking_mqtt_message_processor/model"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func OpenBookedGateHandler(dbRepo model.DatabaseRepository) mqtt.MessageHandler {
	return func(c mqtt.Client, m mqtt.Message) {
		go controller.OpenBookedGateController(c, m, dbRepo)
	}
}
