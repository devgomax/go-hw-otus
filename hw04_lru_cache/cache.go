package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mutex    *sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		mutex:    new(sync.Mutex),
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

type cachedValue struct {
	key   Key
	value any
}

func (lc *lruCache) Set(key Key, value any) bool {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	cv := cachedValue{key: key, value: value}

	item, ok := lc.items[key]
	if ok {
		item.Value = cv
		lc.queue.MoveToFront(item)
		return true
	}

	lc.items[key] = lc.queue.PushFront(cv)

	if lc.queue.Len() > lc.capacity {
		last := lc.queue.Back()
		lc.queue.Remove(last)
		delete(lc.items, last.Value.(cachedValue).key)
	}

	return false
}

func (lc *lruCache) Get(key Key) (any, bool) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	if item, ok := lc.items[key]; ok {
		lc.queue.MoveToFront(item)
		return item.Value.(cachedValue).value, true
	}

	return nil, false
}

func (lc *lruCache) Clear() {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	lc.items = make(map[Key]*ListItem, lc.capacity)
	lc.queue = NewList()
}
