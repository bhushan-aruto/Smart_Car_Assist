package controller

import (
	"encoding/json"
	"log"
	"time"

	"github.com/bhushan-aruto/smart_parking_mqtt_message_processor/model"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func Gate1ControlController(c mqtt.Client, m mqtt.Message, cacheRepo model.CacheRepository) {

	slots, err := cacheRepo.GetSlotsStatus("s1", "s2", "s3", "s4")

	if err != nil {
		log.Printf("error occurred with redis while getting the slots status, err -> %v", err)
		return
	}

	for slotId, slotStatus := range slots {
		if slotStatus == 0 {
			if err := cacheRepo.OfflineBooking(slotId); err != nil {
				log.Printf("error occurred with redis while offline booking the slot, err -> %v\n", err.Error())
				return
			}

			go func() {
				time.Sleep(time.Second * 20)

				slotStatus, err := cacheRepo.GetSlotStatus(slotId)

				if err != nil {
					log.Printf("error occured with redis while getting the slot status, err -> %v\n", err.Error())
					return
				}

				if slotStatus == 2 {
					if err := cacheRepo.CancelOfflineBooking(slotId); err != nil {
						log.Printf("error occured with redis while canceling the offline booking , err -> %v\n", err.Error())
						return
					}
				}

			}()
			response := &model.GateOpenControlResponse{
				MessageType:    2,
				GateOpenStatus: 1,
				SlotId:         slotId,
			}

			responseJson, err := json.Marshal(response)

			if err != nil {
				log.Printf("failed to encode to json, err -> %v\n", err)
				return
			}

			c.Publish("smart/controller/1", 1, false, responseJson)
			return
		}
	}

	response := &model.GateOpenControlResponse{
		MessageType:    2,
		GateOpenStatus: 0,
		SlotId:         "",
	}

	responseJson, err := json.Marshal(response)

	if err != nil {
		log.Printf("failed to encode to json, err -> %v\n", err)
		return
	}

	c.Publish("smart/controller/1", 1, false, responseJson)
}
