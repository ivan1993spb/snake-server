package playground

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrPGNotContainsDot = errors.New("Playground doesn't contain dot")
	ErrNilPlayground    = errors.New("Playground is nil")
	ErrInvalid_W_or_H   = errors.New("Invalid width or height")
	ErrPassedNilObject  = errors.New("Passed nil object")
)

// Playground object contains all objects on map
type Playground struct {
	width, height uint8             // Width and height of map
	objects       map[oid]Object    // All objects on map
	lastUpdates   map[oid]time.Time // Last updates of objects
}

// NewPlayground returns new empty playground
func NewPlayground(width, height uint8) (*Playground, error) {
	if width*height == 0 {
		return nil, fmt.Errorf("Cannot create playground: %s",
			ErrInvalid_W_or_H)
	}

	return &Playground{width, height, make(map[oid]Object),
		make(map[oid]time.Time)}, nil
}

// GetSize returns width and height of playground
func (pg *Playground) GetSize() (uint8, uint8) {
	if pg == nil {
		return 1, 1
	}
	return pg.width, pg.height
}

// GetArea returns playground area
func (pg *Playground) GetArea() uint16 {
	var width, height = pg.GetSize()
	return uint16(width) * uint16(height)
}

// RandomDot generates random dot located on playground
func (pg *Playground) RandomDot() *Dot {
	var width, height = pg.GetSize()
	return NewRandomDotOnSquare(0, 0, width, height)
}

// Occupied returns true if passed dot already used by any object
// located on playground
func (pg *Playground) Occupied(dot *Dot) bool {
	return pg.GetObjectByDot(dot) != nil
}

// GetObjectByDot returns object which contains passed dot
func (pg *Playground) GetObjectByDot(dot *Dot) Object {
	if pg.Contains(dot) {
		for _, object := range pg.objects {
			for i := uint16(0); i < object.DotCount(); i++ {
				if object.Dot(i).Equals(dot) {
					return object
				}
			}
		}
	}
	return nil
}

// Locate tries to create object to playground
func (pg *Playground) Locate(object Object) error {
	if object == nil {
		return fmt.Errorf("Cannot locate: %s", ErrPassedNilObject)
	}
	// Return error if object is already located on playground
	if pg.Located(object) {
		return errors.New("Object is already located")
	}
	// Check each dot of passed object
	for i := uint16(0); i < object.DotCount(); i++ {
		var dot = object.Dot(i)
		// Return error if any dot is occupied or invalid

		if !pg.Contains(dot) {
			return ErrPGNotContainsDot
		}
		if pg.Occupied(dot) {
			return errors.New("Dot is occupied")
		}
	}

	// Object count can't be more than playground area
	var maxId = oid(pg.GetArea())

	// Add to object list of playground
	for id := oid(0); id < maxId; id++ {
		if _, ok := pg.objects[id]; !ok {
			pg.objects[id] = object
			return nil
		}
	}

	return errors.New("Playground is full")
}

// Located returns true if passed object is located on playground
func (pg *Playground) Located(object Object) bool {
	if object != nil {
		for i := range pg.objects {
			if pg.objects[i] == object {
				return true
			}
		}
	}

	return false
}

// Contains return true if playground contains passed dot
func (pg *Playground) Contains(dot *Dot) bool {
	var (
		x, y          = dot.Position()
		width, height = pg.GetSize()
	)
	return width > x && height > y
}

// Delete deletes passed object from playground and returns error if
// there is a problem
func (pg *Playground) Delete(object Object) error {
	if object == nil {
		return fmt.Errorf("Cannot delocate: %s", ErrPassedNilObject)
	}

	if pg.Located(object) {
		for id := range pg.objects {
			if pg.objects[id] == object {
				// Delete object from object storage
				delete(pg.objects, id)
				// Delete information about object last update
				delete(pg.lastUpdates, id)

				return nil
			}
		}
	}

	return errors.New("Passed object isn't located")
}

// Pack packs playground in accordance with standard ST_1
func (pg *Playground) Pack() (output string) {

	// Updates which was detected now
	var currentUpdates = make(map[oid]time.Time)

	for id, object := range pg.objects {
		if shifting, ok := object.(Shifting); ok {
			currentUpdates[id] = shifting.Updated()

			if _, ok := pg.lastUpdates[id]; ok {
				if currentUpdates[id] == pg.lastUpdates[id] {
					// Object was not updated => add onlu ID
					output += "," + id.String()
				} else {
					// Object was added earlier and was updated =>
					// pack only last changes
					output +=
						"," +
							"'" + objectType(object) + "'" +
							id.String() +
							"[" + shifting.PackChanges() + "]"
				}
			} else {
				// Object just was added => pack full data
				output +=
					"," +
						"'" + objectType(object) + "'" +
						id.String() +
						"[" + object.Pack() + "]"
			}
		} else {
			if _, ok := pg.lastUpdates[id]; ok {
				// Object was not removed => add only ID
				output += "," + id.String()
				currentUpdates[id] = pg.lastUpdates[id]
			} else {
				// Object was just added => pack full data
				output +=
					"," +
						"'" + objectType(object) + "'" +
						id.String() +
						"[" + object.Pack() + "]"
				currentUpdates[id] = time.Unix(0, 0)
			}
		}
	}

	pg.lastUpdates = currentUpdates

	if len(output) > 0 {
		output = output[1:]
	}
	return
}

// Updated returns true if any object on playground was updated or
// just added or deleted
func (pg *Playground) Updated() bool {
	// If an object was created
	if len(pg.objects) > len(pg.lastUpdates) {
		return true
	}

	// Check each object
	for i, object := range pg.objects {
		// Is object already created?
		if _, ok := pg.lastUpdates[i]; ok {
			if shifting, ok := object.(Shifting); ok {
				// Is shifting object updated?
				if shifting.Updated() != pg.lastUpdates[i] {
					// Existian object was updated
					return true
				}
			}
		} else {
			// An object was created
			return true
		}
	}

	// Updates was not found
	return false
}

// PackObjects packs objects to string in accordance with standard
// ST_1 by passed object identifiers
func (pg *Playground) PackObjects(ids []int) (output string) {
	if len(ids) > 0 {
		for i := range ids {
			var id = oid(ids[i])

			if object, ok := pg.objects[id]; ok {
				output +=
					// Delimiter
					"," +
						// Type name of object
						"'" + objectType(object) + "'" +
						// Object ID on playground
						id.String() +
						// Object data
						"[" + object.Pack() + "]"
			}
		}

		if len(output) > 0 {
			output = output[1:]
		}
	}

	return
}

type errNavigation struct {
	err error
}

func (e *errNavigation) Error() string {
	return "Cannot navigate: " + e.err.Error()
}

// Navigate calculates and returns dot placed on distance dis dots
// from passed dot in direction dir
func (pg *Playground) Navigate(dot *Dot, dir Direction, dis int16,
) (*Dot, error) {
	// Check direction
	if !ValidDirection(dir) {
		return nil, &errNavigation{ErrInvalidDirection}
	}
	// If distance is zero return passed dot
	if dis == 0 {
		return dot, nil
	}
	// Playground must contain passed dot
	if !pg.Contains(dot) {
		return nil, &errNavigation{ErrPGNotContainsDot}
	}

	// Get default dot if passed nil dot
	if dot == nil {
		dot = NewDefaultDot()
	}

	var distance uint8
	if dis > 0 {
		distance = uint8(dis)
	} else if dis < 0 {
		// reverse direction if passed negative distance
		distance = uint8(-1 * dis)
		dir = ReverseDirection(dir)
	}

	var (
		pgWidth, pgHeight = pg.GetSize()
		dotX, dotY        = dot.Position()
	)

	// North and south

	if dir == DIR_NORTH || dir == DIR_SOUTH {
		if distance > pgHeight {
			distance = distance % pgHeight
		}

		// North
		if dir == DIR_NORTH {
			if distance > dotY {
				return &Dot{dotX, pgHeight - distance + dotY}, nil
			}
			return &Dot{dotX, dotY - distance}, nil
		}

		// South
		if dotY+distance+1 > pgHeight {
			return &Dot{dotX, distance - pgHeight + dotY}, nil
		}
		return &Dot{dotX, dotY + distance}, nil

	}

	// East and west

	if distance > pgWidth {
		distance = distance % pgWidth
	}

	// East
	if dir == DIR_EAST {
		if pgWidth > dotX+distance {
			return &Dot{dotX + distance, dotY}, nil
		}
		return &Dot{distance - pgWidth + dotX, dotY}, nil
	}

	// West
	if distance > dotX {
		return &Dot{pgWidth - distance + dotX, dotY}, nil
	}
	return &Dot{dotX - distance, dotY}, nil
}
