package playground

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ivan1993spb/snake-server/engine"
)

func Test_Playground_CreateObject(t *testing.T) {
	object := &struct{}{}
	location := engine.Location{engine.Dot{0, 0}}

	pg, err := NewPlayground(100, 100)
	require.Nil(t, err)

	err = pg.CreateObject(object, location)
	require.Nil(t, err)
}

func Test_Playground_CreateObjectRandomRect(t *testing.T) {
	object := &struct{}{}

	pg, err := NewPlayground(100, 100)
	require.Nil(t, err)

	location, err := pg.CreateObjectRandomRect(object, 10, 10)
	require.Nil(t, err)
	require.Len(t, location, 100)
}

func Test_Playground_UpdateObject(t *testing.T) {
	t.SkipNow()
}

func Benchmark_Playground_UpdateObject(b *testing.B) {
	// TODO: Implement benchmark.
}

func Test_Playground_DeleteObject(t *testing.T) {
	t.SkipNow()
	// TODO: Implement test.
}

func Test_Playground_CreateObjectAvailableDots_EmptyScene(t *testing.T) {
	t.SkipNow()
	// TODO: Implement test.
}

func Test_Playground_CreateObjectAvailableDots_LocationNotAvailable(t *testing.T) {
	t.SkipNow()
	// TODO: Implement test.
}

func Test_Playground_CreateObjectAvailableDots_LocationsIntersects(t *testing.T) {
	t.SkipNow()
	// TODO: Implement test.
}

func Test_Playground_UpdateObjectAvailableDots_SuccessfullyUpdates(t *testing.T) {
	t.SkipNow()
	// TODO: Implement test.
}
