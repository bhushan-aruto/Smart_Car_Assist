package handler

import (
	"github.com/bhushan-aruto/smart_parking_mqtt_message_processor/controller"
	"github.com/bhushan-aruto/smart_parking_mqtt_message_processor/model"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func Gate2ControlHandler(cacheRepo model.CacheRepository) mqtt.MessageHandler {
	return func(c mqtt.Client, m mqtt.Message) {
		go controller.Gate2ControlController(c, m, cacheRepo)
	}
}
