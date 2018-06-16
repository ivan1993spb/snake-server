package cmap

import (
	"sort"
	"strconv"
	"sync"
	"testing"
)

const (
	keyElephant = iota
	keyMonkey
	keyMoney
	keyMarine
	keyPredator
	keyNone
	keyHorse
)

const defaultShardCount = 32

type Animal struct {
	name string
}

func TestMapCreation(t *testing.T) {
	m := New(defaultShardCount)
	if m == nil {
		t.Error("map is null.")
	}

	if m.Count() != 0 {
		t.Error("new map should be empty.")
	}
}

func TestInsert(t *testing.T) {
	m := New(defaultShardCount)
	elephant := Animal{"elephant"}
	monkey := Animal{"monkey"}

	m.Set(keyElephant, elephant)
	m.Set(keyMonkey, monkey)

	if m.Count() != 2 {
		t.Error("map should contain exactly two elements.")
	}
}

func TestInsertAbsent(t *testing.T) {
	m := New(defaultShardCount)
	elephant := Animal{"elephant"}
	monkey := Animal{"monkey"}

	m.SetIfAbsent(keyElephant, elephant)
	if ok := m.SetIfAbsent(keyElephant, monkey); ok {
		t.Error("map set a new value even the entry is already present")
	}
}

func TestGet(t *testing.T) {
	m := New(defaultShardCount)

	// Get a missing element.
	val, ok := m.Get(keyMoney)

	if ok == true {
		t.Error("ok should be false when item is missing from map.")
	}

	if val != nil {
		t.Error("Missing values should return as null.")
	}

	elephant := Animal{"elephant"}
	m.Set(keyElephant, elephant)

	// Retrieve inserted element.

	tmp, ok := m.Get(keyElephant)
	elephant = tmp.(Animal) // Type assertion.

	if ok == false {
		t.Error("ok should be true for item stored within the map.")
	}

	if &elephant == nil {
		t.Error("expecting an element, not null.")
	}

	if elephant.name != "elephant" {
		t.Error("item was modified.")
	}
}

func TestHas(t *testing.T) {
	m := New(defaultShardCount)

	// Get a missing element.
	if m.Has(keyMoney) == true {
		t.Error("element shouldn't exists")
	}

	elephant := Animal{"elephant"}
	m.Set(keyElephant, elephant)

	if m.Has(keyElephant) == false {
		t.Error("element exists, expecting Has to return True.")
	}
}

func TestRemove(t *testing.T) {
	m := New(defaultShardCount)

	monkey := Animal{"monkey"}
	m.Set(keyMonkey, monkey)

	m.Remove(keyMonkey)

	if m.Count() != 0 {
		t.Error("Expecting count to be zero once item was removed.")
	}

	temp, ok := m.Get(keyMonkey)

	if ok != false {
		t.Error("Expecting ok to be false for missing items.")
	}

	if temp != nil {
		t.Error("Expecting item to be nil after its removal.")
	}

	// Remove a none existing element.
	m.Remove(keyNone)
}

func TestRemoveCb(t *testing.T) {
	m := New(defaultShardCount)

	monkey := Animal{"monkey"}
	m.Set(keyMonkey, monkey)
	elephant := Animal{"elephant"}
	m.Set(keyElephant, elephant)

	var (
		mapKey   uint16
		mapVal   interface{}
		wasFound bool
	)
	cb := func(key uint16, val interface{}, exists bool) bool {
		mapKey = key
		mapVal = val
		wasFound = exists

		if animal, ok := val.(Animal); ok {
			return animal.name == "monkey"
		}
		return false
	}

	// Monkey should be removed
	result := m.RemoveCb(keyMonkey, cb)
	if !result {
		t.Errorf("Result was not true")
	}

	if mapKey != keyMonkey {
		t.Error("Wrong key was provided to the callback")
	}

	if mapVal != monkey {
		t.Errorf("Wrong value was provided to the value")
	}

	if !wasFound {
		t.Errorf("Key was not found")
	}

	if m.Has(keyMonkey) {
		t.Errorf("Key was not removed")
	}

	// Elephant should not be removed
	result = m.RemoveCb(keyElephant, cb)
	if result {
		t.Errorf("Result was true")
	}

	if mapKey != keyElephant {
		t.Error("Wrong key was provided to the callback")
	}

	if mapVal != elephant {
		t.Errorf("Wrong value was provided to the value")
	}

	if !wasFound {
		t.Errorf("Key was not found")
	}

	if !m.Has(keyElephant) {
		t.Errorf("Key was removed")
	}

	// Unset key should remain unset
	result = m.RemoveCb(keyHorse, cb)
	if result {
		t.Errorf("Result was true")
	}

	if mapKey != keyHorse {
		t.Error("Wrong key was provided to the callback")
	}

	if mapVal != nil {
		t.Errorf("Wrong value was provided to the value")
	}

	if wasFound {
		t.Errorf("Key was found")
	}

	if m.Has(keyHorse) {
		t.Errorf("Key was created")
	}
}

func TestPop(t *testing.T) {
	m := New(defaultShardCount)

	monkey := Animal{"monkey"}
	m.Set(keyMonkey, monkey)

	v, exists := m.Pop(keyMonkey)

	if !exists {
		t.Error("Pop didn't find a monkey.")
	}

	m1, ok := v.(Animal)

	if !ok || m1 != monkey {
		t.Error("Pop found something else, but monkey.")
	}

	v2, exists2 := m.Pop(keyMonkey)
	m1, ok = v2.(Animal)

	if exists2 || ok || m1 == monkey {
		t.Error("Pop keeps finding monkey")
	}

	if m.Count() != 0 {
		t.Error("Expecting count to be zero once item was Pop'ed.")
	}

	temp, ok := m.Get(keyMonkey)

	if ok != false {
		t.Error("Expecting ok to be false for missing items.")
	}

	if temp != nil {
		t.Error("Expecting item to be nil after its removal.")
	}
}

func TestCount(t *testing.T) {
	m := New(defaultShardCount)
	for i := 0; i < 100; i++ {
		m.Set(uint16(i), Animal{strconv.Itoa(i)})
	}

	if m.Count() != 100 {
		t.Error("Expecting 100 element within map.")
	}
}

func TestIsEmpty(t *testing.T) {
	m := New(defaultShardCount)

	if m.IsEmpty() == false {
		t.Error("new map should be empty")
	}

	m.Set(keyElephant, Animal{"elephant"})

	if m.IsEmpty() != false {
		t.Error("map shouldn't be empty.")
	}
}

func TestIterator(t *testing.T) {
	m := New(defaultShardCount)

	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Set(uint16(i), Animal{strconv.Itoa(i)})
	}

	counter := 0
	// Iterate over elements.
	for item := range m.Iter() {
		val := item.Val

		if val == nil {
			t.Error("Expecting an object.")
		}
		counter++
	}

	if counter != 100 {
		t.Error("We should have counted 100 elements.")
	}
}

func TestBufferedIterator(t *testing.T) {
	m := New(defaultShardCount)

	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Set(uint16(i), Animal{strconv.Itoa(i)})
	}

	counter := 0
	// Iterate over elements.
	for item := range m.IterBuffered() {
		val := item.Val

		if val == nil {
			t.Error("Expecting an object.")
		}
		counter++
	}

	if counter != 100 {
		t.Error("We should have counted 100 elements.")
	}
}

func TestIterCb(t *testing.T) {
	m := New(defaultShardCount)

	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Set(uint16(i), Animal{strconv.Itoa(i)})
	}

	counter := 0
	// Iterate over elements.
	m.IterCb(func(key uint16, v interface{}) {
		_, ok := v.(Animal)
		if !ok {
			t.Error("Expecting an animal object")
		}

		counter++
	})
	if counter != 100 {
		t.Error("We should have counted 100 elements.")
	}
}

func TestItems(t *testing.T) {
	m := New(defaultShardCount)

	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Set(uint16(i), Animal{strconv.Itoa(i)})
	}

	items := m.Items()

	if len(items) != 100 {
		t.Error("We should have counted 100 elements.")
	}
}

func TestConcurrent(t *testing.T) {
	m := New(defaultShardCount)
	ch := make(chan int)
	const iterations = 1000
	var a [iterations]int

	// Using go routines insert 1000 ints into our map.
	go func() {
		for i := 0; i < iterations/2; i++ {
			// Add item to map.
			m.Set(uint16(i), i)

			// Retrieve item from map.
			val, _ := m.Get(uint16(i))

			// Write to channel inserted value.
			ch <- val.(int)
		} // Call go routine with current index.
	}()

	go func() {
		for i := iterations / 2; i < iterations; i++ {
			// Add item to map.
			m.Set(uint16(i), i)

			// Retrieve item from map.
			val, _ := m.Get(uint16(i))

			// Write to channel inserted value.
			ch <- val.(int)
		} // Call go routine with current index.
	}()

	// Wait for all go routines to finish.
	counter := 0
	for elem := range ch {
		a[counter] = elem
		counter++
		if counter == iterations {
			break
		}
	}

	// Sorts array, will make is simpler to verify all inserted values we're returned.
	sort.Ints(a[0:iterations])

	// Make sure map contains 1000 elements.
	if m.Count() != iterations {
		t.Error("Expecting 1000 elements.")
	}

	// Make sure all inserted values we're fetched from map.
	for i := 0; i < iterations; i++ {
		if i != a[i] {
			t.Error("missing value", i)
		}
	}
}

func TestKeys(t *testing.T) {
	m := New(defaultShardCount)

	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Set(uint16(i), Animal{strconv.Itoa(i)})
	}

	keys := m.Keys()
	if len(keys) != 100 {
		t.Error("We should have counted 100 elements.")
	}
}

func TestMInsert(t *testing.T) {
	animals := map[uint16]interface{}{
		keyElephant: Animal{"elephant"},
		keyMonkey:   Animal{"monkey"},
	}
	m := New(defaultShardCount)
	m.MSet(animals)

	if m.Count() != 2 {
		t.Error("map should contain exactly two elements.")
	}
}

func TestUpsert(t *testing.T) {
	dolphin := Animal{"dolphin"}
	whale := Animal{"whale"}
	tiger := Animal{"tiger"}
	lion := Animal{"lion"}

	cb := func(exists bool, valueInMap interface{}, newValue interface{}) interface{} {
		nv := newValue.(Animal)
		if !exists {
			return []Animal{nv}
		}
		res := valueInMap.([]Animal)
		return append(res, nv)
	}

	m := New(defaultShardCount)
	m.Set(keyMarine, []Animal{dolphin})
	m.Upsert(keyMarine, whale, cb)
	m.Upsert(keyPredator, tiger, cb)
	m.Upsert(keyPredator, lion, cb)

	if m.Count() != 2 {
		t.Error("map should contain exactly two elements.")
	}

	compare := func(a, b []Animal) bool {
		if a == nil || b == nil {
			return false
		}

		if len(a) != len(b) {
			return false
		}

		for i, v := range a {
			if v != b[i] {
				return false
			}
		}
		return true
	}

	marineAnimals, ok := m.Get(keyMarine)
	if !ok || !compare(marineAnimals.([]Animal), []Animal{dolphin, whale}) {
		t.Error("Set, then Upsert failed")
	}

	predators, ok := m.Get(keyPredator)
	if !ok || !compare(predators.([]Animal), []Animal{tiger, lion}) {
		t.Error("Upsert, then Upsert failed")
	}
}

func TestKeysWhenRemoving(t *testing.T) {
	m := New(defaultShardCount)

	// Insert 100 elements.
	Total := 100
	for i := 0; i < Total; i++ {
		m.Set(uint16(i), Animal{strconv.Itoa(i)})
	}

	wg := sync.WaitGroup{}

	// Remove 10 elements concurrently.
	Num := 10
	wg.Add(Num)

	for i := 0; i < Num; i++ {
		go func(c *ConcurrentMap, n uint16) {
			c.Remove(n)
			wg.Done()
		}(m, uint16(i))
	}

	wg.Wait()

	keys := m.Keys()

	if len(keys) != Total-Num {
		t.Error("Invalid keys count")
	}
}

//
func TestUnDrainedIter(t *testing.T) {
	m := New(defaultShardCount)
	// Insert 100 elements.
	Total := 100
	for i := 0; i < Total; i++ {
		m.Set(uint16(i), Animal{strconv.Itoa(i)})
	}
	counter := 0
	// Iterate over elements.
	ch := m.Iter()
	for item := range ch {
		val := item.Val

		if val == nil {
			t.Error("Expecting an object.")
		}
		counter++
		if counter == 42 {
			break
		}
	}
	for i := Total; i < 2*Total; i++ {
		m.Set(uint16(i), Animal{strconv.Itoa(i)})
	}
	for item := range ch {
		val := item.Val

		if val == nil {
			t.Error("Expecting an object.")
		}
		counter++
	}

	if counter != 100 {
		t.Error("We should have been right where we stopped")
	}

	counter = 0
	for item := range m.IterBuffered() {
		val := item.Val

		if val == nil {
			t.Error("Expecting an object.")
		}
		counter++
	}

	if counter != 200 {
		t.Error("We should have counted 200 elements.")
	}
}

func TestUnDrainedIterBuffered(t *testing.T) {
	m := New(defaultShardCount)
	// Insert 100 elements.
	Total := 100
	for i := 0; i < Total; i++ {
		m.Set(uint16(i), Animal{strconv.Itoa(i)})
	}
	counter := 0
	// Iterate over elements.
	ch := m.IterBuffered()
	for item := range ch {
		val := item.Val

		if val == nil {
			t.Error("Expecting an object.")
		}
		counter++
		if counter == 42 {
			break
		}
	}
	for i := Total; i < 2*Total; i++ {
		m.Set(uint16(i), Animal{strconv.Itoa(i)})
	}
	for item := range ch {
		val := item.Val

		if val == nil {
			t.Error("Expecting an object.")
		}
		counter++
	}

	if counter != 100 {
		t.Error("We should have been right where we stopped")
	}

	counter = 0
	for item := range m.IterBuffered() {
		val := item.Val

		if val == nil {
			t.Error("Expecting an object.")
		}
		counter++
	}

	if counter != 200 {
		t.Error("We should have counted 200 elements.")
	}
}
