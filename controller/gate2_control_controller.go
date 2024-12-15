package controller

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/bhushan-aruto/smart_parking_mqtt_message_processor/model"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func Gate2ControlController(c mqtt.Client, m mqtt.Message, cacheRepo model.CacheRepository) {
	var mqttMessage model.RfidRequestMessage

	if err := json.Unmarshal(m.Payload(), &mqttMessage); err != nil {
		log.Printf("error occurred while deserializing the json, err -> %v\n", err.Error())
	}

	slotId, err := cacheRepo.GetSlotIdByRfid(mqttMessage.Rfid)

	if err != nil {
		log.Printf("error occurred with redis while getting the slotd id , err -> %v\n", err.Error())
		return
	}

	inTime, outTime, err := cacheRepo.GetSlotTimings(slotId)

	if err != nil {
		log.Printf("error occured with redis while getting the slot in and out time, err -> %v\n", err.Error())
		return
	}

	outTime = strings.Split(outTime, " m=")[0]
	inTime = strings.Split(inTime, " m=")[0]

	layout := "2006-01-02 15:04:05.999999999 -0700 MST"
	time1, err1 := time.Parse(layout, outTime)
	time2, err2 := time.Parse(layout, inTime)

	if err1 != nil {
		log.Printf("error occurred while parsing the time, err -> %v\n", err1.Error())
		return
	}

	if err2 != nil {
		log.Printf("error occurred while parsing the time, err -> %v\n", err2.Error())
		return
	}

	duration := time1.Sub(time2)

	slotUsageTime := duration.Seconds()

	cost := int32(slotUsageTime) * 1

	response := model.SlotUsageResponse{
		MessageType: 3,
		SlotId:      slotId,
		InTime:      inTime,
		OutTime:     outTime,
		Cost:        cost,
	}

	responseJson, err := json.Marshal(response)

	if err != nil {
		log.Printf("error occured while deserializing the json, err -> %v\n", responseJson)
		return
	}

	c.Publish("smart/controller/1", 1, false, responseJson)
	c.Publish("app", 1, false, responseJson)
}
