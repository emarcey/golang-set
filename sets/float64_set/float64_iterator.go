package mapsetfloat64

// Iterator defines an iterator over a Set, its C channel can be used to range over the Set's
// elements.
type Float64Iterator struct {
	C    <-chan float64
	stop chan struct{}
}

// Stop stops the Float64Iterator, no further elements will be received on C, C will be closed.
func (i *Float64Iterator) Stop() {
	// Allows for Stop() to be called multiple times
	// (close() panics when called on already closed channel)
	defer func() {
		recover()
	}()

	close(i.stop)

	// Exhaust any remaining elements.
	for range i.C {
	}
}

// newFloat64Iterator returns a new Float64Iterator instance together with its item and stop channels.
func newFloat64Iterator() (*Float64Iterator, chan<- float64, <-chan struct{}) {
	itemChan := make(chan float64)
	stopChan := make(chan struct{})
	return &Float64Iterator{
		C:    itemChan,
		stop: stopChan,
	}, itemChan, stopChan
}
