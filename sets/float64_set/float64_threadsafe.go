package mapsetfloat64

import (
	"sync"
)

type threadSafeFloat64Set struct {
	s threadUnsafeFloat64Set
	sync.RWMutex
}

func newThreadSafeFloat64Set() threadSafeFloat64Set {
	return threadSafeFloat64Set{s: newThreadUnsafeFloat64Set()}
}

func (set *threadSafeFloat64Set) Add(i float64) bool {
	set.Lock()
	ret := set.s.Add(i)
	set.Unlock()
	return ret
}

func (set *threadSafeFloat64Set) Contains(i ...float64) bool {
	set.RLock()
	ret := set.s.Contains(i...)
	set.RUnlock()
	return ret
}

func (set *threadSafeFloat64Set) IsSubset(other Float64Set) bool {
	o := other.(*threadSafeFloat64Set)

	set.RLock()
	o.RLock()

	ret := set.s.IsSubset(&o.s)
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeFloat64Set) IsProperSubset(other Float64Set) bool {
	o := other.(*threadSafeFloat64Set)

	set.RLock()
	defer set.RUnlock()
	o.RLock()
	defer o.RUnlock()

	return set.s.IsProperSubset(&o.s)
}

func (set *threadSafeFloat64Set) IsSuperset(other Float64Set) bool {
	return other.IsSubset(set)
}

func (set *threadSafeFloat64Set) IsProperSuperset(other Float64Set) bool {
	return other.IsProperSubset(set)
}

func (set *threadSafeFloat64Set) Union(other Float64Set) Float64Set {
	o := other.(*threadSafeFloat64Set)

	set.RLock()
	o.RLock()

	unsafeUnion := set.s.Union(&o.s).(*threadUnsafeFloat64Set)
	ret := &threadSafeFloat64Set{s: *unsafeUnion}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeFloat64Set) Intersect(other Float64Set) Float64Set {
	o := other.(*threadSafeFloat64Set)

	set.RLock()
	o.RLock()

	unsafeIntersection := set.s.Intersect(&o.s).(*threadUnsafeFloat64Set)
	ret := &threadSafeFloat64Set{s: *unsafeIntersection}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeFloat64Set) Difference(other Float64Set) Float64Set {
	o := other.(*threadSafeFloat64Set)

	set.RLock()
	o.RLock()

	unsafeDifference := set.s.Difference(&o.s).(*threadUnsafeFloat64Set)
	ret := &threadSafeFloat64Set{s: *unsafeDifference}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeFloat64Set) SymmetricDifference(other Float64Set) Float64Set {
	o := other.(*threadSafeFloat64Set)

	set.RLock()
	o.RLock()

	unsafeDifference := set.s.SymmetricDifference(&o.s).(*threadUnsafeFloat64Set)
	ret := &threadSafeFloat64Set{s: *unsafeDifference}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeFloat64Set) Clear() {
	set.Lock()
	set.s = newThreadUnsafeFloat64Set()
	set.Unlock()
}

func (set *threadSafeFloat64Set) Remove(i float64) {
	set.Lock()
	delete(set.s, i)
	set.Unlock()
}

func (set *threadSafeFloat64Set) Cardinality() int {
	set.RLock()
	defer set.RUnlock()
	return len(set.s)
}

func (set *threadSafeFloat64Set) Each(cb func(float64) bool) {
	set.RLock()
	for elem := range set.s {
		if cb(elem) {
			break
		}
	}
	set.RUnlock()
}

func (set *threadSafeFloat64Set) Iter() <-chan float64 {
	ch := make(chan float64)
	go func() {
		set.RLock()

		for elem := range set.s {
			ch <- elem
		}
		close(ch)
		set.RUnlock()
	}()

	return ch
}

func (set *threadSafeFloat64Set) Iterator() *Float64Iterator {
	iterator, ch, stopCh := newFloat64Iterator()

	go func() {
		set.RLock()
	L:
		for elem := range set.s {
			select {
			case <-stopCh:
				break L
			case ch <- elem:
			}
		}
		close(ch)
		set.RUnlock()
	}()

	return iterator
}

func (set *threadSafeFloat64Set) Equal(other Float64Set) bool {
	o := other.(*threadSafeFloat64Set)

	set.RLock()
	o.RLock()

	ret := set.s.Equal(&o.s)
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeFloat64Set) Clone() Float64Set {
	set.RLock()

	unsafeClone := set.s.Clone().(*threadUnsafeFloat64Set)
	ret := &threadSafeFloat64Set{s: *unsafeClone}
	set.RUnlock()
	return ret
}

func (set *threadSafeFloat64Set) String() string {
	set.RLock()
	ret := set.s.String()
	set.RUnlock()
	return ret
}

/*
// Not yet supported
func (set *threadSafeFloat64Set) PowerSet() Float64Set {
    set.RLock()
    unsafePowerSet := set.s.PowerSet().(*threadUnsafeFloat64Set)
    set.RUnlock()

    ret := &threadSafeFloat64Set{s: newThreadUnsafeFloat64Set()}
    for subset := range unsafePowerSet.Iter() {
        unsafeSubset := subset.(*threadUnsafeFloat64Set)
        ret.Add(&threadSafeFloat64Set{s: *unsafeSubset})
    }
    return ret
}
*/

func (set *threadSafeFloat64Set) Pop() float64 {
	set.Lock()
	defer set.Unlock()
	return set.s.Pop()
}

/*
// Not yet supported
func (set *threadSafeFloat64Set) CartesianProduct(other Float64Set) Float64Set {
    o := other.(*threadSafeFloat64Set)

    set.RLock()
    o.RLock()

    unsafeCartProduct := set.s.CartesianProduct(&o.s).(*threadUnsafeFloat64Set)
    ret := &threadSafeFloat64Set{s: *unsafeCartProduct}
    set.RUnlock()
    o.RUnlock()
    return ret
}
*/

func (set *threadSafeFloat64Set) ToSlice() []float64 {
	keys := make([]float64, 0, set.Cardinality())
	set.RLock()
	for elem := range set.s {
		keys = append(keys, elem)
	}
	set.RUnlock()
	return keys
}

func (set *threadSafeFloat64Set) MarshalJSON() ([]byte, error) {
	set.RLock()
	b, err := set.s.MarshalJSON()
	set.RUnlock()

	return b, err
}

func (set *threadSafeFloat64Set) UnmarshalJSON(p []byte) error {
	set.RLock()
	err := set.s.UnmarshalJSON(p)
	set.RUnlock()

	return err
}
