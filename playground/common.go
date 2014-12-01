package playground

import (
	"math/rand"
	"reflect"
	"strconv"
	"time"
)

// Each object must implements Object interface
type Object interface {
	DotCount() uint16  // DotCount must return dot count
	Dot(i uint16) *Dot // Dot returns dot by index
	Pack() string      // Pack converts object data to string
}

// Moving and shifting objects must implement Shifting interface.
// For objects which implements Object and Shifting interfaces method
// Pack returns all object data and method PackChanges returns only
// last updates.
type Shifting interface {
	PackChanges() string // PackChanges packs last changes
	Updated() time.Time  // Updated returns last updating time
}

// oid represents object identifier in playground
type oid uint16

func (i oid) String() string {
	return strconv.Itoa(int(i))
}

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

// objectType returns type name of passed object
func objectType(object Object) string {
	return reflect.Indirect(reflect.ValueOf(object)).Type().Name()
}
