package playground

import (
	"encoding/json"
	"math/rand"
	"time"
)

type Entity interface {
	Dot(i uint16) *Dot // Dot returns dot by index
	DotCount() uint16  // DotCount must return dot count
}

type Json interface {
	PackJson() (json.RawMessage, error)
}

type Object interface {
	Entity
}

var random = rand.New(rand.NewSource(time.Now().UnixNano()))
