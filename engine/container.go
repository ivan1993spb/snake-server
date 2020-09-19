package engine

// Object is a game object
type Object interface{}

// Container contains an object placed on a map
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
