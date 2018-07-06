package objects

import "github.com/ivan1993spb/snake-server/engine"

// Food interface describes methods that must be implemented all edible objects
type Food interface {
	Bite(dot engine.Dot) (nv uint16, success bool, err error)
}

// Alive interface describes methods that must be implemented all alive objects
type Alive interface {
	Hit(dot engine.Dot, force float32) (success bool, err error)
}

// Object interface describes methods that must be implemented all not alive objects
type Object interface {
	Break(dot engine.Dot, force float32) (success bool, err error)
}
