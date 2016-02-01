package objects

// Object interfaces
type (
	Object interface {
		Pack() interface{}
	}

	Shifting interface {
		PackShifts() interface{}
	}
)

// objects use GameProcessor to notifying about its states
type GameProcessor interface {
	OccurredError(object interface{}, err error)
	OccurredCreating(object interface{})
	OccurredDeleting(object interface{})
	OccurredUpdating(object interface{})
}
