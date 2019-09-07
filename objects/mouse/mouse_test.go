package mouse

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ivan1993spb/snake-server/engine"
)

func Test_Mouse_MarshalJSON(t *testing.T) {
	tests := []struct {
		mouse *Mouse
		json  []byte
		err   error
	}{
		{
			mouse: &Mouse{
				id: 22,

				dot: engine.Dot{
					X: 2,
					Y: 3,
				},
				direction: engine.DirectionEast,

				world: nil,
				mux:   &sync.RWMutex{},
			},
			json: []byte(`{"type":"mouse","id":22,"dot":[2,3],"direction":"east"}`),
			err:  nil,
		},
		{
			mouse: &Mouse{
				id: 4294967295,

				dot: engine.Dot{
					X: 255,
					Y: 255,
				},
				direction: engine.DirectionSouth,

				world: nil,
				mux:   &sync.RWMutex{},
			},
			json: []byte(`{"type":"mouse","id":4294967295,"dot":[255,255],"direction":"south"}`),
			err:  nil,
		},
		{
			mouse: &Mouse{
				id: 5492,

				dot: engine.Dot{
					X: 0,
					Y: 0,
				},
				direction: engine.DirectionWest,

				world: nil,
				mux:   &sync.RWMutex{},
			},
			json: []byte(`{"type":"mouse","id":5492,"dot":[0,0],"direction":"west"}`),
			err:  nil,
		},
	}

	for i, test := range tests {
		json, err := test.mouse.MarshalJSON()
		require.Equal(t, test.json, json, fmt.Sprintf("json number: %d", i))
		require.Equal(t, test.err, err, fmt.Sprintf("error number: %d", i))
	}
}
