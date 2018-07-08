package objects

import "github.com/ivan1993spb/snake-server/engine"

// Food interface describes methods that must be implemented all edible objects
type Food interface {
	// Bites object on dot dot and returns nutritional value nv, success flag -
	// true if dot free and error err if error occurred
	Bite(dot engine.Dot) (nv uint16, success bool, err error)
}

// Alive interface describes methods that must be implemented all alive objects
type Alive interface {
	// Hits object on dot dot with force force and returns success flag - true
	// if dot free and error err if error occurred
	Hit(dot engine.Dot, force float32) (success bool, err error)
}

// Object interface describes methods that must be implemented all not alive objects
type Object interface {
	// Breaks object on dot dot with force force, returns success flag true if dot
	// free and error err if error occurred
	Break(dot engine.Dot, force float32) (success bool, err error)
}
