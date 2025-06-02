package main

import (
	"sync"
	"time"
)

/*
Необходимо написать in-memory кэш, который будет по ключу (uuid пользователя) возвращать профиль и список его заказов.

1. У кэша должен быть TTL (2 сек)
2. Кэшем может пользоваться функция(-и), которая работает с заказами (добавляет/обновляет/удаляет). Если TTL истек, то возвращается nil. При апдейте TTL снова устанавливается 2 сек. Методы должны быть потокобезопасными
3. Должны быть написаны тестовые сценарии использования данного кэша
(базовые структуры не менять)

Доп задание: автоматическая очистка истекших записей

type Profile struct {
	UUID   string
	Name   string
	Orders []*Order
}

type Order struct {
	UUID string
	Value any
	CreatedAt time. Time
	UpdatedAt time.Time
}
*/

const TtlSec = 2

type CacheItem struct {
	ProfileValue *Profile
	ExpiresAt    time.Time
}

type CustomCache struct {
	syncedMap sync.Map
}

func NewCustomCache() *CustomCache {
	customCache := &CustomCache{
		syncedMap: sync.Map{},
	}
	customCache.AutomaticCleanUp()
	return customCache
}

func (cc *CustomCache) GetProfile(uuid string) *Profile {
	value, ok := cc.syncedMap.Load(uuid)
	if !ok {
		return nil
	}
	item, okConversion := value.(CacheItem)
	if !okConversion {
		return nil
	}
	if time.Now().Before(item.ExpiresAt) {
		return item.ProfileValue
	}
	return nil
}

func (cc *CustomCache) SetUpdateProfile(uuid string, profile *Profile) {
	cc.syncedMap.Store(uuid, CacheItem{
		ProfileValue: profile,
		ExpiresAt:    time.Now().Add(TtlSec * time.Second),
	})
}

func (cc *CustomCache) DeleteProfile(uuid string) {
	cc.syncedMap.Delete(uuid)
}

func (cc *CustomCache) AutomaticCleanUp() {
	go func() {
		for {
			time.Sleep(TtlSec * time.Second)
			cc.syncedMap.Range(func(key, value any) bool {
				item, okConversion := value.(CacheItem)
				if !okConversion {
					return true
				}
				if time.Now().After(item.ExpiresAt) {
					cc.syncedMap.Delete(key)
				}
				return true
			})
		}
	}()
}
