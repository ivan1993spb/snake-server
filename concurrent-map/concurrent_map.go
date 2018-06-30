package cmap

import (
	"encoding/json"
	"errors"
	"sync"
)

const defaultShardCount = 32

const bufferSize = 64

// A "thread" safe map of type uint16:Anything.
// To avoid lock bottlenecks this map is dived to several (shardCount) map shards.
type ConcurrentMap struct {
	shards []*ConcurrentMapShared
	count  int
}

// A "thread" safe uint16 to anything map.
type ConcurrentMapShared struct {
	items map[uint16]interface{}
	mux   *sync.RWMutex
}

// Creates a new concurrent map.
func New(shardCount int) (*ConcurrentMap, error) {
	if shardCount < 1 {
		return nil, errors.New("invalid shard count: less than 1")
	}

	shards := make([]*ConcurrentMapShared, shardCount)

	for i := 0; i < shardCount; i++ {
		shards[i] = &ConcurrentMapShared{
			items: make(map[uint16]interface{}),
			mux:   &sync.RWMutex{},
		}
	}

	return &ConcurrentMap{
		shards: shards,
		count:  shardCount,
	}, nil
}

func NewDefault() *ConcurrentMap {
	m, err := New(defaultShardCount)
	if err != nil {
		panic(err)
	}
	return m
}

// Returns shard under given key
func (m *ConcurrentMap) getShard(key uint16) *ConcurrentMapShared {
	return m.shards[m.getShardIndex(key)]
}

func (m *ConcurrentMap) getShardIndex(key uint16) uint16 {
	return key % uint16(m.count)
}

func (m *ConcurrentMap) sortShardsTuples(data map[uint16]interface{}) map[uint16][]Tuple {
	shardsTuples := map[uint16][]Tuple{}

	for key, value := range data {
		shardIndex := m.getShardIndex(key)

		if _, ok := shardsTuples[shardIndex]; !ok {
			shardsTuples[shardIndex] = make([]Tuple, 0, bufferSize)
		}

		shardsTuples[shardIndex] = append(shardsTuples[shardIndex], Tuple{
			Key: key,
			Val: value,
		})
	}

	return shardsTuples
}

func (m *ConcurrentMap) sortShardsKeys(keys []uint16) map[uint16][]uint16 {
	shardsKeys := map[uint16][]uint16{}

	for _, key := range keys {
		shardIndex := m.getShardIndex(key)

		if _, ok := shardsKeys[shardIndex]; !ok {
			shardsKeys[shardIndex] = make([]uint16, 0, bufferSize)
		}

		shardsKeys[shardIndex] = append(shardsKeys[shardIndex], key)
	}

	return shardsKeys
}

func (m *ConcurrentMap) MSet(data map[uint16]interface{}) {
	shardsTuples := m.sortShardsTuples(data)

	for shardIndex, tuples := range shardsTuples {
		shard := m.shards[shardIndex]
		shard.mux.Lock()
		for _, tuple := range tuples {
			shard.items[tuple.Key] = tuple.Val
		}
		shard.mux.Unlock()
	}
}

func (m *ConcurrentMap) MSetIfAllAbsent(data map[uint16]interface{}) bool {
	shardsTuples := m.sortShardsTuples(data)
	rollbackKeys := make([]uint16, 0, len(data))

	for shardIndex, tuples := range shardsTuples {
		shard := m.shards[shardIndex]
		shard.mux.Lock()
		for _, tuple := range tuples {
			if _, ok := shard.items[tuple.Key]; !ok {
				shard.items[tuple.Key] = tuple.Val
				rollbackKeys = append(rollbackKeys, tuple.Key)
			} else {
				shard.mux.Unlock()

				// Rollback
				m.MRemove(rollbackKeys)

				return false
			}
		}
		shard.mux.Unlock()
	}

	return true
}

// MSetIfAbsent sets the given value only if key has no value associated with it.
func (m *ConcurrentMap) MSetIfAbsent(data map[uint16]interface{}) []uint16 {
	shardsTuples := m.sortShardsTuples(data)
	keys := make([]uint16, 0, len(data))

	for shardIndex, tuples := range shardsTuples {
		shard := m.shards[shardIndex]
		shard.mux.Lock()
		for _, tuple := range tuples {
			if _, ok := shard.items[tuple.Key]; !ok {
				shard.items[tuple.Key] = tuple.Val
				keys = append(keys, tuple.Key)
			}
		}
		shard.mux.Unlock()
	}

	return keys
}

// Sets the given value under the specified key.
func (m *ConcurrentMap) Set(key uint16, value interface{}) {
	// Get map shard.
	shard := m.getShard(key)
	shard.mux.Lock()
	shard.items[key] = value
	shard.mux.Unlock()
}

// Callback to return new element to be inserted into the map
// It is called while lock is held, therefore it MUST NOT
// try to access other keys in same map, as it can lead to deadlock since
// Go sync.RWLock is not reentrant
type UpsertCb func(exist bool, valueInMap interface{}, newValue interface{}) interface{}

// Insert or Update - updates existing element or inserts a new one using UpsertCb
func (m *ConcurrentMap) Upsert(key uint16, value interface{}, cb UpsertCb) (res interface{}) {
	shard := m.getShard(key)
	shard.mux.Lock()
	v, ok := shard.items[key]
	res = cb(ok, v, value)
	shard.items[key] = res
	shard.mux.Unlock()
	return res
}

// Sets the given value under the specified key if no value was associated with it.
func (m *ConcurrentMap) SetIfAbsent(key uint16, value interface{}) bool {
	// Get map shard.
	shard := m.getShard(key)
	shard.mux.Lock()
	_, ok := shard.items[key]
	if !ok {
		shard.items[key] = value
	}
	shard.mux.Unlock()
	return !ok
}

// Retrieves an element from map under given key.
func (m *ConcurrentMap) Get(key uint16) (interface{}, bool) {
	// Get shard
	shard := m.getShard(key)
	shard.mux.RLock()
	// Get item from shard.
	val, ok := shard.items[key]
	shard.mux.RUnlock()
	return val, ok
}

func (m *ConcurrentMap) MGet(keys []uint16) map[uint16]interface{} {
	shardsKeys := m.sortShardsKeys(keys)
	items := make(map[uint16]interface{})

	for shardIndex, keys := range shardsKeys {
		shard := m.shards[shardIndex]
		shard.mux.RLock()

		for _, key := range keys {
			if value, ok := shard.items[key]; ok {
				items[key] = value
			}
		}

		shard.mux.RUnlock()
	}

	return items
}

// Returns the number of elements within the map.
func (m *ConcurrentMap) Count() int {
	count := 0
	for i := 0; i < m.count; i++ {
		shard := m.shards[i]
		shard.mux.RLock()
		count += len(shard.items)
		shard.mux.RUnlock()
	}
	return count
}

// Looks up an item under specified key
func (m *ConcurrentMap) Has(key uint16) bool {
	// Get shard
	shard := m.getShard(key)
	shard.mux.RLock()
	// See if element is within shard.
	_, ok := shard.items[key]
	shard.mux.RUnlock()
	return ok
}

func (m *ConcurrentMap) HasAny(keys []uint16) bool {
	flagHasAny := false

	shardsKeys := m.sortShardsKeys(keys)

	for shardIndex, keys := range shardsKeys {
		shard := m.shards[shardIndex]
		shard.mux.RLock()

		for _, key := range keys {
			if _, ok := shard.items[key]; ok {
				flagHasAny = true
				break
			}
		}

		shard.mux.RUnlock()

		if flagHasAny {
			break
		}
	}

	return flagHasAny
}

func (m *ConcurrentMap) HasAll(keys []uint16) bool {
	flagHasAll := true

	shardsKeys := m.sortShardsKeys(keys)

	for shardIndex, keys := range shardsKeys {
		shard := m.shards[shardIndex]
		shard.mux.RLock()

		for _, key := range keys {
			if _, ok := shard.items[key]; !ok {
				flagHasAll = false
				break
			}
		}

		shard.mux.RUnlock()

		if !flagHasAll {
			break
		}
	}

	return flagHasAll
}

// Removes an element from the map.
func (m *ConcurrentMap) Remove(key uint16) {
	// Try to get shard.
	shard := m.getShard(key)
	shard.mux.Lock()
	delete(shard.items, key)
	shard.mux.Unlock()
}

func (m *ConcurrentMap) MRemove(keys []uint16) {
	shardsKeys := m.sortShardsKeys(keys)

	for shardIndex, keys := range shardsKeys {
		shard := m.shards[shardIndex]
		shard.mux.Lock()
		for _, key := range keys {
			delete(shard.items, key)
		}
		shard.mux.Unlock()
	}
}

// RemoveCb is a callback executed in a map.RemoveCb() call, while Lock is held
// If returns true, the element will be removed from the map
type RemoveCb func(key uint16, v interface{}, exists bool) bool

// RemoveCb locks the shard containing the key, retrieves its current value and calls the callback with those params
// If callback returns true and element exists, it will remove it from the map
// Returns the value returned by the callback (even if element was not present in the map)
func (m *ConcurrentMap) RemoveCb(key uint16, cb RemoveCb) bool {
	// Try to get shard.
	shard := m.getShard(key)
	shard.mux.Lock()
	v, ok := shard.items[key]
	remove := cb(key, v, ok)
	if remove && ok {
		delete(shard.items, key)
	}
	shard.mux.Unlock()
	return remove
}

// Removes an element from the map and returns it
func (m *ConcurrentMap) Pop(key uint16) (v interface{}, exists bool) {
	// Try to get shard.
	shard := m.getShard(key)
	shard.mux.Lock()
	v, exists = shard.items[key]
	delete(shard.items, key)
	shard.mux.Unlock()
	return v, exists
}

// Checks if map is empty.
func (m *ConcurrentMap) IsEmpty() bool {
	return m.Count() == 0
}

// Used by the Iter & IterBuffered functions to wrap two variables together over a channel,
type Tuple struct {
	Key uint16
	Val interface{}
}

// Returns an iterator which could be used in a for range loop.
//
// Deprecated: using IterBuffered() will get a better performence
func (m *ConcurrentMap) Iter() <-chan Tuple {
	chans := snapshot(m)
	ch := make(chan Tuple)
	go fanIn(chans, ch)
	return ch
}

// Returns a buffered iterator which could be used in a for range loop.
func (m *ConcurrentMap) IterBuffered() <-chan Tuple {
	chans := snapshot(m)
	total := 0
	for _, c := range chans {
		total += cap(c)
	}
	ch := make(chan Tuple, total)
	go fanIn(chans, ch)
	return ch
}

// Returns a array of channels that contains elements in each shard,
// which likely takes a snapshot of `m`.
// It returns once the size of each buffered channel is determined,
// before all the channels are populated using goroutines.
func snapshot(m *ConcurrentMap) (chans []chan Tuple) {
	chans = make([]chan Tuple, m.count)
	wg := sync.WaitGroup{}
	wg.Add(m.count)
	// Foreach shard.
	for index, shard := range m.shards {
		go func(index int, shard *ConcurrentMapShared) {
			// Foreach key, value pair.
			shard.mux.RLock()
			chans[index] = make(chan Tuple, len(shard.items))
			wg.Done()
			for key, val := range shard.items {
				chans[index] <- Tuple{key, val}
			}
			shard.mux.RUnlock()
			close(chans[index])
		}(index, shard)
	}
	wg.Wait()
	return chans
}

// fanIn reads elements from channels `chans` into channel `out`
func fanIn(chans []chan Tuple, out chan Tuple) {
	wg := sync.WaitGroup{}
	wg.Add(len(chans))
	for _, ch := range chans {
		go func(ch chan Tuple) {
			for t := range ch {
				out <- t
			}
			wg.Done()
		}(ch)
	}
	wg.Wait()
	close(out)
}

// Returns all items as map[uint16]interface{}
func (m *ConcurrentMap) Items() map[uint16]interface{} {
	tmp := make(map[uint16]interface{})

	// Insert items to temporary map.
	for item := range m.IterBuffered() {
		tmp[item.Key] = item.Val
	}

	return tmp
}

// Iterator callback,called for every key,value found in
// maps. RLock is held for all calls for a given shard
// therefore callback sess consistent view of a shard,
// but not across the shards
type IterCb func(key uint16, v interface{})

// Callback based iterator, cheapest way to read
// all elements in a map.
func (m *ConcurrentMap) IterCb(fn IterCb) {
	for idx := range m.shards {
		shard := m.shards[idx]
		shard.mux.RLock()
		for key, value := range shard.items {
			fn(key, value)
		}
		shard.mux.RUnlock()
	}
}

// Return all keys as []uint16
func (m *ConcurrentMap) Keys() []uint16 {
	count := m.Count()
	ch := make(chan uint16, count)
	go func() {
		// Foreach shard.
		wg := sync.WaitGroup{}
		wg.Add(m.count)
		for _, shard := range m.shards {
			go func(shard *ConcurrentMapShared) {
				// Foreach key, value pair.
				shard.mux.RLock()
				for key := range shard.items {
					ch <- key
				}
				shard.mux.RUnlock()
				wg.Done()
			}(shard)
		}
		wg.Wait()
		close(ch)
	}()

	// Generate keys
	keys := make([]uint16, 0, count)
	for k := range ch {
		keys = append(keys, k)
	}
	return keys
}

//Reviles ConcurrentMap "private" variables to json marshal.
func (m *ConcurrentMap) MarshalJSON() ([]byte, error) {
	// Create a temporary map, which will hold all item spread across shards.
	tmp := make(map[uint16]interface{})

	// Insert items to temporary map.
	for item := range m.IterBuffered() {
		tmp[item.Key] = item.Val
	}
	return json.Marshal(tmp)
}

// Concurrent map uses Interface{} as its value, therefor JSON Unmarshal
// will probably won't know which to type to unmarshal into, in such case
// we'll end up with a value of type map[uint16]interface{}, In most cases this isn't
// out value type, this is why we've decided to remove this functionality.

// func (m *ConcurrentMap) UnmarshalJSON(b []byte) (err error) {
// 	// Reverse process of Marshal.

// 	tmp := make(map[uint16]interface{})

// 	// Unmarshal into a single map.
// 	if err := json.Unmarshal(b, &tmp); err != nil {
// 		return nil
// 	}

// 	// foreach key,value pair in temporary map insert into our concurrent map.
// 	for key, val := range tmp {
// 		m.Set(key, val)
// 	}
// 	return nil
// }
