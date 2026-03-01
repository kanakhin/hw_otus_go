package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mu       sync.Mutex
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.mu.Lock()
	defer func() { l.mu.Unlock() }()

	item := cacheItem{value: value, key: key}
	listItem, ok := l.items[key]
	if ok {
		listItem.Value = item
		l.queue.MoveToFront(listItem)

		return true
	}

	listItem = l.queue.PushFront(item)
	l.items[key] = listItem

	if l.queue.Len() > l.capacity {
		lastItem := l.queue.Back()
		lastItemValue := lastItem.Value
		keyToDelete := lastItemValue.(cacheItem).key

		delete(l.items, keyToDelete)
		l.queue.Remove(lastItem)
	}

	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mu.Lock()
	defer func() { l.mu.Unlock() }()

	listItem, ok := l.items[key]

	if !ok {
		return nil, false
	}

	l.queue.MoveToFront(listItem)
	value := listItem.Value.(cacheItem).value

	return value, true
}

func (l *lruCache) Clear() {
	l.mu.Lock()
	defer func() { l.mu.Unlock() }()

	l.items = make(map[Key]*ListItem, l.capacity)
	l.queue = NewList()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
