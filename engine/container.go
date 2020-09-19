package engine

// Container contains an object placed on a map
type Container struct {
	object interface{}
}

// NewContainer creates a new container with the passed object inside
func NewContainer(object interface{}) *Container {
	return &Container{
		object: object,
	}
}

// GetObject returns an object from the container
func (o *Container) GetObject() interface{} {
	return o.object
}
