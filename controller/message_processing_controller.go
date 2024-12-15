package controller

import (
	"encoding/json"
	"log"
	"time"

	"github.com/bhushan-aruto/smart_parking_mqtt_message_processor/model"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func MessageProcessingController(c mqtt.Client, m mqtt.Message, cacheRepo model.CacheRepository) {

	var mqttMessage model.ProcessRequestMessage

	if err := json.Unmarshal(m.Payload(), &mqttMessage); err != nil {
		log.Printf("error while decoding json,Error -> %v\n", err.Error())
		return
	}

	slots, err := cacheRepo.GetSlotsStatus("s1", "s2", "s3", "s4")

	if err != nil {
		log.Printf("error occurred with database while getting the slots status, Error -> %v\n", err.Error())
		return
	}

	var slot1Status int32 = slots["s1"]
	var slot2Status int32 = slots["s2"]
	var slot3Status int32 = slots["s3"]
	var slot4Status int32 = slots["s4"]

	slot1Channel := make(chan int32)
	slot2Channel := make(chan int32)
	slot3Channel := make(chan int32)
	slot4Channel := make(chan int32)

	currentTime := time.Now().String()

	go func() {
		if slot1Status != mqttMessage.Slot1 {
			if mqttMessage.Slot1 == 1 {
				if err := cacheRepo.UpdateSlotUsageStartStatus("s1", currentTime); err != nil {
					log.Printf("error occurred with database while updating the slot status, Error -> %v\n", err.Error())
					slot1Channel <- mqttMessage.Slot1
					return
				}
			} else {
				if slot1Status != 2 {
					if err := cacheRepo.UpdateSlotUsageStopStatus("s1", currentTime); err != nil {
						log.Printf("error occurred with database while updating the slot status, Error -> %v\n", err.Error())
						slot1Channel <- mqttMessage.Slot1
						return
					}
				} else {
					slot1Channel <- 2
				}

			}
		}

		slot1Channel <- mqttMessage.Slot1
	}()

	go func() {

		if slot2Status != mqttMessage.Slot2 {
			if mqttMessage.Slot2 == 1 {
				if err := cacheRepo.UpdateSlotUsageStartStatus("s2", currentTime); err != nil {
					log.Printf("error occurred with database while updating the slot status, Error -> %v\n", err.Error())
					slot2Channel <- mqttMessage.Slot2
					return
				}
			} else {
				if slot2Status != 2 {
					if err := cacheRepo.UpdateSlotUsageStopStatus("s2", currentTime); err != nil {
						log.Printf("error occurred with database while updating the slot status, Error -> %v\n", err.Error())
						slot2Channel <- mqttMessage.Slot2
						return
					}
				} else {
					slot2Channel <- 2
				}

			}
		}

		slot2Channel <- mqttMessage.Slot2

	}()

	go func() {

		if slot3Status != mqttMessage.Slot3 {
			if mqttMessage.Slot3 == 1 {
				if err := cacheRepo.UpdateSlotUsageStartStatus("s3", currentTime); err != nil {
					log.Printf("error occurred with database while updating the slot status, Error -> %v\n", err.Error())
					slot3Channel <- mqttMessage.Slot3
					return
				}

			} else {
				if slot3Status != 2 {
					if err := cacheRepo.UpdateSlotUsageStopStatus("s3", currentTime); err != nil {
						log.Printf("error occurred with database while updating the slot status, Error -> %v\n", err.Error())
						slot3Channel <- mqttMessage.Slot3
						return
					}

				} else {
					slot3Channel <- 2
				}

			}
		}

		slot3Channel <- mqttMessage.Slot3

	}()

	go func() {
		if slot4Status != mqttMessage.Slot4 {
			if mqttMessage.Slot4 == 1 {
				if err := cacheRepo.UpdateSlotUsageStartStatus("s4", currentTime); err != nil {
					log.Printf("error occurred with database while updating the slot status, Error -> %v\n", err.Error())
					slot4Channel <- mqttMessage.Slot4
					return
				}

			} else {
				if slot4Status != 2 {
					if err := cacheRepo.UpdateSlotUsageStopStatus("s4", currentTime); err != nil {
						log.Printf("error occurred with database while updating the slot status, Error -> %v\n", err.Error())
						slot4Channel <- mqttMessage.Slot4
						return
					}

				} else {
					slot4Channel <- 2
				}

			}
		}

		slot4Channel <- mqttMessage.Slot4
	}()

	slot1Status = <-slot1Channel
	slot2Status = <-slot2Channel
	slot3Status = <-slot3Channel
	slot4Status = <-slot4Channel

	response := &model.ProcessResponseMessage{
		MessageType: 1,
		Slot1:       slot1Status,
		Slot2:       slot2Status,
		Slot3:       slot3Status,
		Slot4:       slot4Status,
	}

	publishMessage, err := json.Marshal(response)

	if err != nil {
		log.Printf("error occurred while publishing the message, Error -> %v\n", err.Error())
		return
	}
	c.Publish("app", 1, false, publishMessage)

	c.Publish("smart/controller/1", 1, false, publishMessage)

}
