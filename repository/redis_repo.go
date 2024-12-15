package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	cache *redis.Client
}

func NewRedisRepository(cache *redis.Client) *RedisRepository {
	return &RedisRepository{
		cache,
	}
}

func (repo *RedisRepository) GetSlotsStatus(slotsId ...string) (map[string]int32, error) {
	slots := make(map[string]int32, 4)

	for _, slotId := range slotsId {

		slotStatusKey := fmt.Sprintf("%v_status", slotId)

		slotStatusStr, err := repo.cache.Get(context.Background(), slotStatusKey).Result()

		if err != nil {
			return nil, err
		}

		slotStatusInt, err := strconv.Atoi(slotStatusStr)

		if err != nil {
			return nil, err
		}

		slotStatusInt32 := int32(slotStatusInt)

		slots[slotId] = slotStatusInt32
	}
	return slots, nil
}

func (repo *RedisRepository) UpdateSlotUsageStartStatus(slotId string, inTime string) error {
	slotStatusKey := fmt.Sprintf("%v_status", slotId)

	slotInTimeKey := fmt.Sprintf("%v_in_time", slotId)

	if _, err := repo.cache.Set(context.Background(), slotStatusKey, 1, time.Duration(0)).Result(); err != nil {
		return err
	}

	if _, err := repo.cache.Set(context.Background(), slotInTimeKey, inTime, time.Duration(0)).Result(); err != nil {
		return err
	}

	return nil
}

func (repo *RedisRepository) UpdateSlotUsageStopStatus(slotId string, outTime string) error {
	slotStatusKey := fmt.Sprintf("%v_status", slotId)

	slotOutTimeKey := fmt.Sprintf("%v_out_time", slotId)

	if _, err := repo.cache.Set(context.Background(), slotStatusKey, 0, time.Duration(0)).Result(); err != nil {
		return err
	}

	if _, err := repo.cache.Set(context.Background(), slotOutTimeKey, outTime, time.Duration(0)).Result(); err != nil {
		return err
	}

	return nil
}

func (repo *RedisRepository) OfflineBooking(slotId string) error {
	slotStatusKey := fmt.Sprintf("%v_status", slotId)
	_, err := repo.cache.Set(context.Background(), slotStatusKey, 2, time.Duration(0)).Result()
	return err
}

func (repo *RedisRepository) GetSlotStatus(slotId string) (int32, error) {
	slotStatusKey := fmt.Sprintf("%v_status", slotId)

	slotStatusStr, err := repo.cache.Get(context.Background(), slotStatusKey).Result()

	if err != nil {
		return -1, err
	}

	slotStatus, err := strconv.Atoi(slotStatusStr)

	if err != nil {
		return -1, err
	}

	return int32(slotStatus), nil

}

func (repo *RedisRepository) CancelOfflineBooking(slotId string) error {
	slotStatusKey := fmt.Sprintf("%v_status", slotId)
	_, err := repo.cache.Set(context.Background(), slotStatusKey, 0, time.Duration(0)).Result()
	return err
}

func (repo *RedisRepository) GetSlotIdByRfid(rfid string) (string, error) {
	slotIdKey := fmt.Sprintf("%v", rfid)

	res, err := repo.cache.Get(context.Background(), slotIdKey).Result()

	if err != nil {
		return "", err
	}

	return res, nil
}

func (repo *RedisRepository) GetSlotTimings(slotId string) (string, string, error) {
	slotInTimeKey := fmt.Sprintf("%v_in_time", slotId)
	slotOutTimeKey := fmt.Sprintf("%v_out_time", slotId)

	slotInTime, err := repo.cache.Get(context.Background(), slotInTimeKey).Result()

	if err != nil {
		return "", "", err
	}

	slotOutTime, err := repo.cache.Get(context.Background(), slotOutTimeKey).Result()

	if err != nil {
		return "", "", err
	}

	return slotInTime, slotOutTime, nil
}
