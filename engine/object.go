package engine

type Object struct {
	value interface{}
}

func NewObject(value interface{}) *Object {
	return &Object{
		value: value,
	}
}

func (o *Object) Value() interface{} {
	return o.value
}
