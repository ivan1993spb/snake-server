package objects

import "github.com/ivan1993spb/snake-server/engine"

// Food interface describes methods which must be implemented by all edible
// objects
type Food interface {
	// Bite bites an object at the passed dot and returns the nutritional
	// value nv, success flag true if the dot has been released or an error err
	// if one occurred
	Bite(dot engine.Dot) (nv uint16, success bool, err error)
}

// Alive interface describes methods which must be implemented by all living
// objects
type Alive interface {
	// Hit hits an object at the passed dot with the given force force and
	// returns success flag true if the dot has been released or an error err
	// if one occurred
	Hit(dot engine.Dot, force float64) (success bool, err error)
}

// Breakable interface describes methods that must be implemented by all
// objects which could be broken
type Breakable interface {
	// Break breaks an object at the passed dot with the given force and
	// returns success flag true if the dot has been released or an error err
	// if one occurred
	Break(dot engine.Dot, force float64) (success bool, err error)
}
