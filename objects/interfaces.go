package objects

import "github.com/ivan1993spb/snake-server/engine"

type Food interface {
	NutritionalValue(dot engine.Dot) uint16
}
