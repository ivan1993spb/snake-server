package objects

import "github.com/ivan1993spb/snake-server/engine"

// Food interface describes methods that must be implemented all edible objects
type Food interface {
	NutritionalValue(dot engine.Dot) uint16
}

type Alive interface {
	Kill(dot engine.Dot)
}

type Strong interface {
	Strength(dot engine.Dot)
}

// Если предмет съедобный - то он кусается Food - Bite
// Если предмет твердый - то он ломается Hard - Break
// Если предмет живой - Alive - Kill
