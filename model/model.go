package model

type DatabaseRepository interface {
	GetUserIdByEmail(email string) (string, error)
	CheckUserBookingExists(userId string) (bool, error)
	DeleteUserBooking(userId string) error
}

type CacheRepository interface {
	GetSlotsStatus(slotsId ...string) (map[string]int32, error)
	UpdateSlotUsageStartStatus(slotId string, inTime string) error
	UpdateSlotUsageStopStatus(slotId string, outTime string) error
	OfflineBooking(slotId string) error
	GetSlotStatus(slotId string) (int32, error)
	GetSlotIdByRfid(rfid string) (string, error)
	GetSlotTimings(slotId string) (string, string, error)
	CancelOfflineBooking(slotId string) error
}

type ProcessRequestMessage struct {
	Slot1 int32 `json:"s1"`
	Slot2 int32 `json:"s2"`
	Slot3 int32 `json:"s3"`
	Slot4 int32 `json:"s4"`
}

type ProcessResponseMessage struct {
	MessageType int32 `json:"mty"`
	Slot1       int32 `json:"s1"`
	Slot2       int32 `json:"s2"`
	Slot3       int32 `json:"s3"`
	Slot4       int32 `json:"s4"`
}

type GateOpenControlResponse struct {
	MessageType    int32  `json:"mty"`
	GateOpenStatus int32  `json:"gos"`
	SlotId         string `json:"sid"`
}

type RfidRequestMessage struct {
	Rfid string `json:"rfid"`
}

type SlotUsageResponse struct {
	MessageType int32  `json:"mty"`
	SlotId      string `json:"sid"`
	InTime      string `json:"itm"`
	OutTime     string `json:"otm"`
	Cost        int32  `json:"cost"`
}

type OpenBookedGateRequestMessage struct {
	Email string `json:"email"`
}

type OpenBookedGateResponse struct {
	MessageType int32 `json:"mty"`
}
