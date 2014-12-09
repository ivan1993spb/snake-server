package objects

type errCreateObject struct {
	err error
}

func (e *errCreateObject) Error() string {
	return "Cannot create object: " + e.err.Error()
}
