package world

import (
	"math"
	"sync"
)

const identifiersBufferSize = 2 << 12

const (
	indexStart uint32 = 0
	indexDelta uint32 = 1
)

type Identifier uint32

type IdentifierRegistry struct {
	index               uint32
	obtainedIdentifiers []Identifier
	mux                 *sync.Mutex
}

func NewIdentifierRegistry() *IdentifierRegistry {
	return &IdentifierRegistry{
		index:               indexStart,
		obtainedIdentifiers: make([]Identifier, 0, identifiersBufferSize),
		mux:                 &sync.Mutex{},
	}
}

func (ir *IdentifierRegistry) Obtain() Identifier {
	ir.mux.Lock()
	defer ir.mux.Unlock()

	ir.incrementIndex()
	id := Identifier(ir.index)
	ir.unsafeRegisterIdentifier(id)

	return id
}

func (ir *IdentifierRegistry) incrementIndex() {
	for {
		if ir.index == math.MaxUint32 {
			ir.index = indexStart
		}

		ir.index += indexDelta

		if !ir.unsafeIsObtainedIdentifier(Identifier(ir.index)) {
			break
		}
	}
}

func (ir *IdentifierRegistry) unsafeRegisterIdentifier(id Identifier) {
	ir.obtainedIdentifiers = append(ir.obtainedIdentifiers, id)
}

func (ir *IdentifierRegistry) unsafeIsObtainedIdentifier(id Identifier) bool {
	for _, obtainedIdentifier := range ir.obtainedIdentifiers {
		if obtainedIdentifier == id {
			return true
		}
	}
	return false
}

func (ir *IdentifierRegistry) Release(id Identifier) {
	ir.mux.Lock()
	defer ir.mux.Unlock()
	ir.unsafeRelease(id)
}

func (ir *IdentifierRegistry) unsafeRelease(id Identifier) {
	for i, obtainedIdentifier := range ir.obtainedIdentifiers {
		if obtainedIdentifier == id {
			ir.obtainedIdentifiers = append(ir.obtainedIdentifiers[:i], ir.obtainedIdentifiers[i+1:]...)
			return
		}
	}
}
