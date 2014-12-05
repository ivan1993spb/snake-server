package playground

// Direction indicates movement direction of a object
type Direction uint8

const (
	DIR_NORTH Direction = iota
	DIR_EAST
	DIR_SOUTH
	DIR_WEST
)

// Count of directions
const _DIR_COUNT = 4

// RandomDirection returns random direction
func RandomDirection() Direction {
	return Direction(random.Intn(_DIR_COUNT))
}

// CalculateDirection calculates direction by two passed dots
func CalculateDirection(from, to *Dot) Direction {
	if !from.Equals(to) {
		var (
			fromX, fromY = from.Position()
			toX, toY     = to.Position()
			diffX, diffY uint8
		)

		if fromX > toX {
			diffX = fromX - toX
		} else {
			diffX = toX - fromX
		}
		if fromY > toY {
			diffY = fromY - toY
		} else {
			diffY = toY - fromY
		}

		if diffX > diffY {
			if toX > fromX {
				return DIR_EAST
			}
			return DIR_WEST
		}
		if diffY > diffX {
			if toY > fromY {
				return DIR_SOUTH
			}
			return DIR_NORTH
		}
	}

	return RandomDirection()
}

// ReverseDirection reverses passed direction dir
func ReverseDirection(dir Direction) Direction {
	switch dir {
	case DIR_NORTH:
		return DIR_SOUTH
	case DIR_EAST:
		return DIR_WEST
	case DIR_SOUTH:
		return DIR_NORTH
	case DIR_WEST:
		return DIR_EAST
	}
	return RandomDirection()
}

// ValidDirection returns true if passed direction is valid
func ValidDirection(dir Direction) bool {
	switch dir {
	case DIR_NORTH, DIR_EAST, DIR_SOUTH, DIR_WEST:
		return true
	}
	return false
}

// Pack packs direction
func (dir Direction) Pack() string {
	switch dir {
	case DIR_NORTH:
		return "!n"
	case DIR_SOUTH:
		return "!s"
	case DIR_EAST:
		return "!e"
	case DIR_WEST:
		return "!w"
	}
	return ""
}
