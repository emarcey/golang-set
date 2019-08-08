package mapsetfloat32

import (
	"sync"
)

type threadSafeFloat32Set struct {
	s threadUnsafeFloat32Set
	sync.RWMutex
}

func newThreadSafeFloat32Set() threadSafeFloat32Set {
	return threadSafeFloat32Set{s: newThreadUnsafeFloat32Set()}
}

func (set *threadSafeFloat32Set) Add(i float32) bool {
	set.Lock()
	ret := set.s.Add(i)
	set.Unlock()
	return ret
}

func (set *threadSafeFloat32Set) Contains(i ...float32) bool {
	set.RLock()
	ret := set.s.Contains(i...)
	set.RUnlock()
	return ret
}

func (set *threadSafeFloat32Set) IsSubset(other Float32Set) bool {
	o := other.(*threadSafeFloat32Set)

	set.RLock()
	o.RLock()

	ret := set.s.IsSubset(&o.s)
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeFloat32Set) IsProperSubset(other Float32Set) bool {
	o := other.(*threadSafeFloat32Set)

	set.RLock()
	defer set.RUnlock()
	o.RLock()
	defer o.RUnlock()

	return set.s.IsProperSubset(&o.s)
}

func (set *threadSafeFloat32Set) IsSuperset(other Float32Set) bool {
	return other.IsSubset(set)
}

func (set *threadSafeFloat32Set) IsProperSuperset(other Float32Set) bool {
	return other.IsProperSubset(set)
}

func (set *threadSafeFloat32Set) Union(other Float32Set) Float32Set {
	o := other.(*threadSafeFloat32Set)

	set.RLock()
	o.RLock()

	unsafeUnion := set.s.Union(&o.s).(*threadUnsafeFloat32Set)
	ret := &threadSafeFloat32Set{s: *unsafeUnion}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeFloat32Set) Intersect(other Float32Set) Float32Set {
	o := other.(*threadSafeFloat32Set)

	set.RLock()
	o.RLock()

	unsafeIntersection := set.s.Intersect(&o.s).(*threadUnsafeFloat32Set)
	ret := &threadSafeFloat32Set{s: *unsafeIntersection}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeFloat32Set) Difference(other Float32Set) Float32Set {
	o := other.(*threadSafeFloat32Set)

	set.RLock()
	o.RLock()

	unsafeDifference := set.s.Difference(&o.s).(*threadUnsafeFloat32Set)
	ret := &threadSafeFloat32Set{s: *unsafeDifference}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeFloat32Set) SymmetricDifference(other Float32Set) Float32Set {
	o := other.(*threadSafeFloat32Set)

	set.RLock()
	o.RLock()

	unsafeDifference := set.s.SymmetricDifference(&o.s).(*threadUnsafeFloat32Set)
	ret := &threadSafeFloat32Set{s: *unsafeDifference}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeFloat32Set) Clear() {
	set.Lock()
	set.s = newThreadUnsafeFloat32Set()
	set.Unlock()
}

func (set *threadSafeFloat32Set) Remove(i float32) {
	set.Lock()
	delete(set.s, i)
	set.Unlock()
}

func (set *threadSafeFloat32Set) Cardinality() int {
	set.RLock()
	defer set.RUnlock()
	return len(set.s)
}

func (set *threadSafeFloat32Set) Each(cb func(float32) bool) {
	set.RLock()
	for elem := range set.s {
		if cb(elem) {
			break
		}
	}
	set.RUnlock()
}

func (set *threadSafeFloat32Set) Iter() <-chan float32 {
	ch := make(chan float32)
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

func (set *threadSafeFloat32Set) Iterator() *Float32Iterator {
	iterator, ch, stopCh := newFloat32Iterator()

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

func (set *threadSafeFloat32Set) Equal(other Float32Set) bool {
	o := other.(*threadSafeFloat32Set)

	set.RLock()
	o.RLock()

	ret := set.s.Equal(&o.s)
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeFloat32Set) Clone() Float32Set {
	set.RLock()

	unsafeClone := set.s.Clone().(*threadUnsafeFloat32Set)
	ret := &threadSafeFloat32Set{s: *unsafeClone}
	set.RUnlock()
	return ret
}

func (set *threadSafeFloat32Set) String() string {
	set.RLock()
	ret := set.s.String()
	set.RUnlock()
	return ret
}

/*
// Not yet supported
func (set *threadSafeFloat32Set) PowerSet() Float32Set {
    set.RLock()
    unsafePowerSet := set.s.PowerSet().(*threadUnsafeFloat32Set)
    set.RUnlock()

    ret := &threadSafeFloat32Set{s: newThreadUnsafeFloat32Set()}
    for subset := range unsafePowerSet.Iter() {
        unsafeSubset := subset.(*threadUnsafeFloat32Set)
        ret.Add(&threadSafeFloat32Set{s: *unsafeSubset})
    }
    return ret
}
*/

func (set *threadSafeFloat32Set) Pop() float32 {
	set.Lock()
	defer set.Unlock()
	return set.s.Pop()
}

/*
// Not yet supported
func (set *threadSafeFloat32Set) CartesianProduct(other Float32Set) Float32Set {
    o := other.(*threadSafeFloat32Set)

    set.RLock()
    o.RLock()

    unsafeCartProduct := set.s.CartesianProduct(&o.s).(*threadUnsafeFloat32Set)
    ret := &threadSafeFloat32Set{s: *unsafeCartProduct}
    set.RUnlock()
    o.RUnlock()
    return ret
}
*/

func (set *threadSafeFloat32Set) ToSlice() []float32 {
	keys := make([]float32, 0, set.Cardinality())
	set.RLock()
	for elem := range set.s {
		keys = append(keys, elem)
	}
	set.RUnlock()
	return keys
}

func (set *threadSafeFloat32Set) MarshalJSON() ([]byte, error) {
	set.RLock()
	b, err := set.s.MarshalJSON()
	set.RUnlock()

	return b, err
}

func (set *threadSafeFloat32Set) UnmarshalJSON(p []byte) error {
	set.RLock()
	err := set.s.UnmarshalJSON(p)
	set.RUnlock()

	return err
}
