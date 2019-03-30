package observers

type Observer interface {
	Observe(stop <-chan struct{})
}
