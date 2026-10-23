package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

// AI generated one for quick demo purpose
type inMemoryCache struct {
	MemLocker
	sync.Map
	data map[string]any
	sets map[string]map[string]struct{}
	hash map[string]map[string]any
}

func NewInMemoryCache() *inMemoryCache {
	return &inMemoryCache{
		MemLocker: NewMemLocker(),
		data:      make(map[string]any),
		sets:      make(map[string]map[string]struct{}),
		hash:      make(map[string]map[string]any),
	}
}
func (c *inMemoryCache) GetObject(key string, dest any) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	val, ok := c.data[key]
	if !ok {
		return false, nil
	}
	bytes, err := json.Marshal(val)
	if err != nil {
		return false, err
	}
	err = json.Unmarshal(bytes, dest)
	return err == nil, err
}

func (c *inMemoryCache) GetValue(key string) (string, bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	val, ok := c.data[key]
	if !ok {
		return "", false, nil
	}
	strVal, ok := val.(string)
	if !ok {
		return "", false, errors.New("value is not a string")
	}
	return strVal, true, nil
}

func (c *inMemoryCache) SetObject(key string, obj any, _ time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = obj
	return nil
}

func (c *inMemoryCache) SetValue(key string, value string, _ time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
	return nil
}

func (c *inMemoryCache) RemoveKey(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
	return nil
}

func (c *inMemoryCache) RemoveKeyWithCount(key string) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.data[key]; ok {
		delete(c.data, key)
		return 1, nil
	}
	return 0, nil
}

func (c *inMemoryCache) RemoveKeysWithCount(keys []string) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	var count int64
	for _, key := range keys {
		if _, ok := c.data[key]; ok {
			delete(c.data, key)
			count++
		}
	}
	return count, nil
}

func (c *inMemoryCache) RemoveKeys(keys ...string) error {
	_, err := c.RemoveKeysWithCount(keys)
	return err
}

func (c *inMemoryCache) AddSet(setKey string, member string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.sets[setKey]; !ok {
		c.sets[setKey] = make(map[string]struct{})
	}
	c.sets[setKey][member] = struct{}{}
	return nil
}

func (c *inMemoryCache) GetSetMembers(setKey string) ([]string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	set, ok := c.sets[setKey]
	if !ok {
		return nil, nil
	}
	members := make([]string, 0, len(set))
	for member := range set {
		members = append(members, member)
	}
	return members, nil
}

func (c *inMemoryCache) RemoveSetMember(setKey string, member string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if set, ok := c.sets[setKey]; ok {
		delete(set, member)
		if len(set) == 0 {
			delete(c.sets, setKey)
		}
	}
	return nil
}

func (c *inMemoryCache) SetH(key string, values map[string]any, _ time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.hash[key]; !exists {
		c.hash[key] = make(map[string]any)
	}

	for field, value := range values {
		c.hash[key][field] = value
	}

	return nil
}

func (c *inMemoryCache) GetH(key string, field string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	fields, ok := c.hash[key]
	if !ok {
		return "", errors.New("key not found")
	}

	val, exists := fields[field]
	if !exists {
		return "", errors.New("field not found")
	}

	return fmt.Sprint(val), nil
}

func (c *inMemoryCache) Exists(key string) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.data[key]; ok {
		return 1, nil
	}
	if _, ok := c.sets[key]; ok {
		return 1, nil
	}
	if _, ok := c.hash[key]; ok {
		return 1, nil
	}
	return 0, nil
}

func (c *inMemoryCache) Incr(key string) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	val, ok := c.data[key]
	if !ok {
		c.data[key] = "1"
		return 1, nil
	}

	str, ok := val.(string)
	if !ok {
		return 0, errors.New("value is not an integer or out of range")
	}
	n, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, errors.New("value is not an integer or out of range")
	}
	n++
	c.data[key] = strconv.FormatInt(n, 10)
	return n, nil
}

type MemLocker struct {
	mu    sync.Mutex
	locks map[string]chan struct{}
}

func NewMemLocker() MemLocker {
	return MemLocker{
		locks: make(map[string]chan struct{}),
	}
}

func (l *MemLocker) Lock(key string) (func() error, error) {
	l.mu.Lock()

	ch, ok := l.locks[key]
	if !ok {
		ch = make(chan struct{}, 1)
		l.locks[key] = ch
	}

	l.mu.Unlock()

	ch <- struct{}{}

	var once sync.Once
	return func() error {
		unlocked := false

		once.Do(func() {
			<-ch
			unlocked = true
		})

		if !unlocked {
			return errors.New("lock already released")
		}

		return nil
	}, nil
}
