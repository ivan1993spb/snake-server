package world

import (
	"fmt"
	"math"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewIdentifierRegistry_CreatesIdentifierRegistry(t *testing.T) {
	ir := NewIdentifierRegistry()

	require.True(t, indexStart == ir.index)

	require.Len(t, ir.obtainedIdentifiers, 0)
	require.Equal(t, identifiersBufferSize, cap(ir.obtainedIdentifiers))

	require.NotNil(t, ir.mux)
}

func Test_IdentifierRegistry_Release_ReleasesIdentifierFromStartCorrectly(t *testing.T) {
	identifiers := []Identifier{1, 2, 3, 4, 5, 6}
	index := indexStart + indexDelta*uint32(len(identifiers))

	ir := &IdentifierRegistry{
		obtainedIdentifiers: identifiers,
		mux:                 &sync.Mutex{},
		index:               index,
	}

	ir.Release(Identifier(1))
	require.Equal(t, []Identifier{2, 3, 4, 5, 6}, ir.obtainedIdentifiers)
}

func Test_IdentifierRegistry_Release_ReleasesIdentifierFromMiddleCorrectly(t *testing.T) {
	identifiers := []Identifier{1, 2, 3, 4, 5, 6}
	index := indexStart + indexDelta*uint32(len(identifiers))

	ir := &IdentifierRegistry{
		obtainedIdentifiers: identifiers,
		mux:                 &sync.Mutex{},
		index:               index,
	}

	ir.Release(Identifier(4))
	require.Equal(t, []Identifier{1, 2, 3, 5, 6}, ir.obtainedIdentifiers)
}

func Test_IdentifierRegistry_Release_ReleasesIdentifierFromEndCorrectly(t *testing.T) {
	identifiers := []Identifier{1, 2, 3, 4, 5, 6}
	index := indexStart + indexDelta*uint32(len(identifiers))

	ir := &IdentifierRegistry{
		obtainedIdentifiers: identifiers,
		mux:                 &sync.Mutex{},
		index:               index,
	}

	ir.Release(Identifier(6))
	require.Equal(t, []Identifier{1, 2, 3, 4, 5}, ir.obtainedIdentifiers)
}

func Test_IdentifierRegistry_unsafeIsObtainedIdentifier_worksCorrectly(t *testing.T) {
	tests := []struct {
		expected    bool
		identifier  Identifier
		identifiers []Identifier
	}{
		{false, 1, []Identifier{}},
		{false, 1, []Identifier{3}},
		{true, 1, []Identifier{1, 2, 3, 4}},
		{false, 5, []Identifier{2, 3, 4}},
		{true, 5, []Identifier{2, 3, 4, 5}},
		{true, 3, []Identifier{2, 3, 4, 5}},
	}

	for number, test := range tests {
		index := indexStart + indexDelta*uint32(len(test.identifiers))
		ir := &IdentifierRegistry{
			obtainedIdentifiers: test.identifiers,
			mux:                 &sync.Mutex{},
			index:               index,
		}
		msg := fmt.Sprintf("error test case: %d", number)
		require.Equal(t, test.expected, ir.unsafeIsObtainedIdentifier(test.identifier), msg)
	}
}

func Test_IdentifierRegistry_incrementIndex_worksCorrectly(t *testing.T) {
	tests := []struct {
		expectedIndex       uint32
		currentIndex        uint32
		obtainedIdentifiers []Identifier
	}{
		{3, 2, []Identifier{}},
		{11, 10, []Identifier{}},
		{indexStart + indexDelta, math.MaxUint32, []Identifier{}},
		{indexStart + indexDelta*4, math.MaxUint32, []Identifier{
			Identifier(indexStart + indexDelta*1),
			Identifier(indexStart + indexDelta*2),
			Identifier(indexStart + indexDelta*3),
			Identifier(indexStart + indexDelta*5),
			Identifier(indexStart + indexDelta*6),
		}},
	}

	for number, test := range tests {
		ir := &IdentifierRegistry{
			obtainedIdentifiers: test.obtainedIdentifiers,
			mux:                 &sync.Mutex{},
			index:               test.currentIndex,
		}

		msg := fmt.Sprintf("error test case: %d", number)

		ir.incrementIndex()

		require.Equal(t, test.expectedIndex, ir.index, msg)
	}
}

func TestIdentifierRegistry_Obtain_ObtainsIdentifierCorrectly(t *testing.T) {
	tests := []struct {
		expectedIdentifier          Identifier
		currentIndex                uint32
		currentObtainedIdentifiers  []Identifier
		expectedObtainedIdentifiers []Identifier
	}{
		{
			3,
			2,
			[]Identifier{},
			[]Identifier{3},
		},
		{
			11,
			10,
			[]Identifier{},
			[]Identifier{11},
		},
		{
			Identifier(indexStart + indexDelta),
			math.MaxUint32,
			[]Identifier{},
			[]Identifier{Identifier(indexStart + indexDelta)},
		},
		{
			Identifier(indexStart + indexDelta*4),
			math.MaxUint32,
			[]Identifier{
				Identifier(indexStart + indexDelta*1),
				Identifier(indexStart + indexDelta*2),
				Identifier(indexStart + indexDelta*3),
				Identifier(indexStart + indexDelta*5),
				Identifier(indexStart + indexDelta*6),
			},
			[]Identifier{
				Identifier(indexStart + indexDelta*1),
				Identifier(indexStart + indexDelta*2),
				Identifier(indexStart + indexDelta*3),
				Identifier(indexStart + indexDelta*5),
				Identifier(indexStart + indexDelta*6),
				Identifier(indexStart + indexDelta*4),
			},
		},
		{
			Identifier(math.MaxUint32),
			math.MaxUint32 - indexDelta,
			[]Identifier{},
			[]Identifier{Identifier(math.MaxUint32)},
		},
	}

	for number, test := range tests {
		ir := &IdentifierRegistry{
			obtainedIdentifiers: test.currentObtainedIdentifiers,
			mux:                 &sync.Mutex{},
			index:               test.currentIndex,
		}

		msg := fmt.Sprintf("error test case: %d", number)

		actualIdentifier := ir.Obtain()

		require.Equal(t, test.expectedIdentifier, actualIdentifier, msg)
	}
}
