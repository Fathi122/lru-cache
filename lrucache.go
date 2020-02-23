package lrucache

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

// ListNode structure
type ListNode struct {
	Val  int
	Next *ListNode
}

// LRUCache structure
type LRUCache struct {
	hashMap     map[int]interface{}
	capacity    int
	lruListHead *ListNode
	lruListTail *ListNode
	lock        *sync.Mutex
}

func init() {
	// switch to DebugLevel for verbose
	log.SetLevel(log.InfoLevel)
}

// constructor LRUcache constructor
func constructor(capacity int) LRUCache {
	return LRUCache{hashMap: make(map[int]interface{}),
		lruListHead: nil,
		lruListTail: nil,
		capacity:    capacity,
		lock:        &sync.Mutex{}}
}

// addToBack add to back in the list
func (e *LRUCache) addToBack(key int) {
	var prev *ListNode
	prev = nil

	newNode := new(ListNode)
	newNode.Val = key
	newNode.Next = nil

	for tmp := e.lruListHead; tmp != nil; tmp = tmp.Next {
		prev = tmp
	}

	if prev == nil {
		e.lruListHead = newNode
		e.lruListTail = newNode
	} else {
		prev.Next = newNode
		e.lruListTail = newNode
	}
}

// pushToBack push to back in the list
func (e *LRUCache) pushToBack(key int) {
	var prev *ListNode
	prev = nil

	if e.lruListTail != nil && e.lruListTail.Val == key {
		log.Debug("Key ", key, " already at back nothing to do")
		return
	}

	for tmp := e.lruListHead; tmp != nil; tmp = tmp.Next {
		if tmp.Val == key {
			if prev != nil {
				prev.Next = tmp.Next
			} else {
				e.lruListHead = tmp.Next
			}
			tmp.Next = nil
			if e.lruListTail != nil {
				e.lruListTail.Next = tmp
			}
			e.lruListTail = tmp
			break
		}
		// update previous
		prev = tmp
	}
}

// updateHead update head in the list
func (e *LRUCache) updateHead() {
	if e.lruListHead != nil {
		tmp := e.lruListHead.Next
		e.lruListHead = tmp
	}
}

// Get get value from cache
func (e *LRUCache) Get(key int) int {
	e.lock.Lock()
	defer e.lock.Unlock()

	if _, keyFound := e.hashMap[key]; !keyFound {
		log.Debug(" Key ", key, " doesn't exists")
		return -1
	}
	log.Debug("Key ", key, " found")

	// pushing it to back
	e.pushToBack(key)
	return e.hashMap[key].(int)
}

// Put put element in cache
func (e *LRUCache) Put(key int, value int) {
	e.lock.Lock()
	defer e.lock.Unlock()

	keyFound := false
	if _, keyFound = e.hashMap[key]; keyFound {
		log.Debug("Key ", key, " already exists")
	}

	// add new element and update list
	e.hashMap[key] = value
	if !keyFound {
		log.Debug("Key ", key, " create entry")
		e.addToBack(key)
	} else {
		e.pushToBack(key)
	}

	// evict current head (least recently used element)
	if len(e.hashMap) > e.capacity {
		log.Debug("Evicting Key ", e.lruListHead.Val)
		delete(e.hashMap, e.lruListHead.Val)
		e.updateHead()
	}
}
