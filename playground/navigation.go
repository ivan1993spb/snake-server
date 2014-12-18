package playground

type errNavigation struct {
	err error
}

func (e *errNavigation) Error() string {
	return "Navigation error: " + e.err.Error()
}

// CalculateDistance calculates distance between two dots
func CalculateDistance(from, to *Dot) (res uint16) {
	if !from.Equals(to) {
		if from.x > to.x {
			res = uint16(from.x - to.x)
		} else {
			res = uint16(to.x - from.x)
		}

		if from.y > to.y {
			res += uint16(from.y - to.y)
		} else {
			res += uint16(to.y - from.y)
		}
	}

	return
}

// CalculateDirection calculates direction by two passed dots
func CalculateDirection(from, to *Dot) Direction {
	if !from.Equals(to) {
		var diffX, diffY uint8

		if from.x > to.x {
			diffX = from.x - to.x
		} else {
			diffX = to.x - from.x
		}
		if from.y > to.y {
			diffY = from.y - to.y
		} else {
			diffY = to.y - from.y
		}

		if diffX > diffY {
			if to.x > from.x {
				return DIR_EAST
			}
			return DIR_WEST
		}

		if diffY > diffX {
			if to.y > from.y {
				return DIR_SOUTH
			}
			return DIR_NORTH
		}
	}

	return RandomDirection()
}

// ReverseDirection reverses passed direction
func ReverseDirection(dir Direction) (Direction, error) {
	switch dir {
	case DIR_NORTH:
		return DIR_SOUTH, nil
	case DIR_EAST:
		return DIR_WEST, nil
	case DIR_SOUTH:
		return DIR_NORTH, nil
	case DIR_WEST:
		return DIR_EAST, nil
	}

	return 0, &errNavigation{ErrInvalidDirection}
}

// Navigate calculates and returns dot placed on distance dis dots
// from passed dot in direction dir
func (pg *Playground) Navigate(dot *Dot, dir Direction, dis uint8,
) (*Dot, error) {
	// If distance is zero return passed dot
	if dis == 0 {
		return dot, nil
	}
	// Playground must contain passed dot
	if !pg.Contains(dot) {
		return nil, &errNavigation{ErrPGNotContainsDot}
	}

	switch dir {
	case DIR_NORTH, DIR_SOUTH:
		if dis > pg.height {
			dis = dis % pg.height
		}

		// North
		if dir == DIR_NORTH {
			if dis > dot.y {
				return &Dot{dot.x, pg.height - dis + dot.y}, nil
			}
			return &Dot{dot.x, dot.y - dis}, nil
		}

		// South
		if dot.y+dis+1 > pg.height {
			return &Dot{dot.x, dis - pg.height + dot.y}, nil
		}
		return &Dot{dot.x, dot.y + dis}, nil

	case DIR_WEST, DIR_EAST:
		if dis > pg.width {
			dis = dis % pg.width
		}

		// East
		if dir == DIR_EAST {
			if pg.width > dot.x+dis {
				return &Dot{dot.x + dis, dot.y}, nil
			}
			return &Dot{dis - pg.width + dot.x, dot.y}, nil
		}

		// West
		if dis > dot.x {
			return &Dot{pg.width - dis + dot.x, dot.y}, nil
		}
		return &Dot{dot.x - dis, dot.y}, nil
	}

	return nil, &errNavigation{ErrInvalidDirection}
}
