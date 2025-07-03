package custom_cache

import (
	"fmt"
	"sync"
)

type CustomCache struct {
	valueMap *sync.Map
}

func NewCustomCache() *CustomCache {
	return &CustomCache{
		valueMap: &sync.Map{},
	}
}

func (c *CustomCache) Put(key string, value string) {
	c.valueMap.Store(key, value)
}

func (c *CustomCache) Get(key string) (string, error) {
	value, ok := c.valueMap.Load(key)
	if !ok {
		return "", fmt.Errorf("error getting value from map, key: %s not found", key)
	}
	res, okConversion := value.(string)
	if !okConversion {
		return "", fmt.Errorf("error converting value to string")
	}
	return res, nil
}
