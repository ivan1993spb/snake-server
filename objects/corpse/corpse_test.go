package corpse

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ivan1993spb/snake-server/engine"
	"github.com/ivan1993spb/snake-server/world"
)

func Test_NewCorpse_CreatesCorpseAndLocatesObject(t *testing.T) {
	w, err := world.NewWorld(100, 100)
	require.Nil(t, err, "cannot initialize world")
	require.NotNil(t, w, "cannot initialize world")

	corpse, err := NewCorpse(w, engine.Location{
		engine.Dot{10, 0},
		engine.Dot{9, 0},
		engine.Dot{8, 0},
		engine.Dot{7, 0},
	})
	require.Nil(t, err)
	require.True(t, corpse.location.Equals(engine.Location{
		engine.Dot{10, 0},
		engine.Dot{9, 0},
		engine.Dot{8, 0},
		engine.Dot{7, 0},
	}))
}

func Test_Corpse_NutritionalValue_ReturnsValidNutritionalValue(t *testing.T) {
	w, err := world.NewWorld(100, 100)
	require.Nil(t, err, "cannot initialize world")
	require.NotNil(t, w, "cannot initialize world")

	corpse := &Corpse{
		world: w,
		location: engine.Location{
			engine.Dot{10, 0},
			engine.Dot{9, 0},
			engine.Dot{8, 0},
			engine.Dot{7, 0},
		},
		mux:  &sync.RWMutex{},
		stop: make(chan struct{}),
	}

	err = w.CreateObject(corpse, engine.Location{
		engine.Dot{10, 0},
		engine.Dot{9, 0},
		engine.Dot{8, 0},
		engine.Dot{7, 0},
	})
	require.Nil(t, err, "cannot create object")

	nutritionalValue := corpse.NutritionalValue(engine.Dot{10, 0})
	require.Equal(t, corpseNutritionalValue, nutritionalValue)
	require.True(t, corpse.location.Equals(engine.Location{
		engine.Dot{9, 0},
		engine.Dot{8, 0},
		engine.Dot{7, 0},
	}))
}

func Test_Corpse_NutritionalValue_ReturnsZeroForInvalidDot(t *testing.T) {
	w, err := world.NewWorld(100, 100)
	require.Nil(t, err, "cannot initialize world")
	require.NotNil(t, w, "cannot initialize world")

	corpse := &Corpse{
		world: w,
		location: engine.Location{
			engine.Dot{10, 0},
			engine.Dot{9, 0},
			engine.Dot{8, 0},
			engine.Dot{7, 0},
		},
		mux:  &sync.RWMutex{},
		stop: make(chan struct{}),
	}

	err = w.CreateObject(corpse, engine.Location{
		engine.Dot{10, 0},
		engine.Dot{9, 0},
		engine.Dot{8, 0},
		engine.Dot{7, 0},
	})
	require.Nil(t, err, "cannot create object")

	nutritionalValue := corpse.NutritionalValue(engine.Dot{10, 10})
	require.Equal(t, uint16(0), nutritionalValue)
	require.Equal(t, engine.Location{
		engine.Dot{10, 0},
		engine.Dot{9, 0},
		engine.Dot{8, 0},
		engine.Dot{7, 0},
	}, corpse.location)
}
