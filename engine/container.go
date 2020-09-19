package engine

// Object is a game object
type Object interface{}

// Container contains an object placed on a map.
//
// As objects are addressed with unsafe's pointers, it is neccessary
// to know the exact type of the stored object. However, as there are
// a lot of different types of objects, you never know what you find
// at a dot. Therefore, a container type is needed to wrap an unknown
// game object with a certain type you can convert from a pointer.
type Container struct {
	object Object
}

// NewContainer creates a new container with the passed object inside
func NewContainer(object Object) *Container {
	return &Container{
		object: object,
	}
}

// GetObject returns an object from the container
func (o *Container) GetObject() Object {
	return o.object
}
