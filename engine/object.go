package engine

// Object is a wrapper for any object to be placed on a map
type Object struct {
	value interface{}
}

// NewObject creates new wrapped object
func NewObject(value interface{}) *Object {
	return &Object{
		value: value,
	}
}

// Value returns the value
func (o *Object) Value() interface{} {
	return o.value
}
