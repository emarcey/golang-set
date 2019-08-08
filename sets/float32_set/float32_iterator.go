package mapsetfloat32

// Iterator defines an iterator over a Set, its C channel can be used to range over the Set's
// elements.
type Float32Iterator struct {
	C    <-chan float32
	stop chan struct{}
}

// Stop stops the Float32Iterator, no further elements will be received on C, C will be closed.
func (i *Float32Iterator) Stop() {
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

// newFloat32Iterator returns a new Float32Iterator instance together with its item and stop channels.
func newFloat32Iterator() (*Float32Iterator, chan<- float32, <-chan struct{}) {
	itemChan := make(chan float32)
	stopChan := make(chan struct{})
	return &Float32Iterator{
		C:    itemChan,
		stop: stopChan,
	}, itemChan, stopChan
}
