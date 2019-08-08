package mapsetuint16

import (
	"sync"
)

type threadSafeUint16Set struct {
	s threadUnsafeUint16Set
	sync.RWMutex
}

func newThreadSafeUint16Set() threadSafeUint16Set {
	return threadSafeUint16Set{s: newThreadUnsafeUint16Set()}
}

func (set *threadSafeUint16Set) Add(i uint16) bool {
	set.Lock()
	ret := set.s.Add(i)
	set.Unlock()
	return ret
}

func (set *threadSafeUint16Set) Contains(i ...uint16) bool {
	set.RLock()
	ret := set.s.Contains(i...)
	set.RUnlock()
	return ret
}

func (set *threadSafeUint16Set) IsSubset(other Uint16Set) bool {
	o := other.(*threadSafeUint16Set)

	set.RLock()
	o.RLock()

	ret := set.s.IsSubset(&o.s)
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeUint16Set) IsProperSubset(other Uint16Set) bool {
	o := other.(*threadSafeUint16Set)

	set.RLock()
	defer set.RUnlock()
	o.RLock()
	defer o.RUnlock()

	return set.s.IsProperSubset(&o.s)
}

func (set *threadSafeUint16Set) IsSuperset(other Uint16Set) bool {
	return other.IsSubset(set)
}

func (set *threadSafeUint16Set) IsProperSuperset(other Uint16Set) bool {
	return other.IsProperSubset(set)
}

func (set *threadSafeUint16Set) Union(other Uint16Set) Uint16Set {
	o := other.(*threadSafeUint16Set)

	set.RLock()
	o.RLock()

	unsafeUnion := set.s.Union(&o.s).(*threadUnsafeUint16Set)
	ret := &threadSafeUint16Set{s: *unsafeUnion}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeUint16Set) Intersect(other Uint16Set) Uint16Set {
	o := other.(*threadSafeUint16Set)

	set.RLock()
	o.RLock()

	unsafeIntersection := set.s.Intersect(&o.s).(*threadUnsafeUint16Set)
	ret := &threadSafeUint16Set{s: *unsafeIntersection}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeUint16Set) Difference(other Uint16Set) Uint16Set {
	o := other.(*threadSafeUint16Set)

	set.RLock()
	o.RLock()

	unsafeDifference := set.s.Difference(&o.s).(*threadUnsafeUint16Set)
	ret := &threadSafeUint16Set{s: *unsafeDifference}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeUint16Set) SymmetricDifference(other Uint16Set) Uint16Set {
	o := other.(*threadSafeUint16Set)

	set.RLock()
	o.RLock()

	unsafeDifference := set.s.SymmetricDifference(&o.s).(*threadUnsafeUint16Set)
	ret := &threadSafeUint16Set{s: *unsafeDifference}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeUint16Set) Clear() {
	set.Lock()
	set.s = newThreadUnsafeUint16Set()
	set.Unlock()
}

func (set *threadSafeUint16Set) Remove(i uint16) {
	set.Lock()
	delete(set.s, i)
	set.Unlock()
}

func (set *threadSafeUint16Set) Cardinality() int {
	set.RLock()
	defer set.RUnlock()
	return len(set.s)
}

func (set *threadSafeUint16Set) Each(cb func(uint16) bool) {
	set.RLock()
	for elem := range set.s {
		if cb(elem) {
			break
		}
	}
	set.RUnlock()
}

func (set *threadSafeUint16Set) Iter() <-chan uint16 {
	ch := make(chan uint16)
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

func (set *threadSafeUint16Set) Iterator() *Uint16Iterator {
	iterator, ch, stopCh := newUint16Iterator()

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

func (set *threadSafeUint16Set) Equal(other Uint16Set) bool {
	o := other.(*threadSafeUint16Set)

	set.RLock()
	o.RLock()

	ret := set.s.Equal(&o.s)
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeUint16Set) Clone() Uint16Set {
	set.RLock()

	unsafeClone := set.s.Clone().(*threadUnsafeUint16Set)
	ret := &threadSafeUint16Set{s: *unsafeClone}
	set.RUnlock()
	return ret
}

func (set *threadSafeUint16Set) String() string {
	set.RLock()
	ret := set.s.String()
	set.RUnlock()
	return ret
}

/*
// Not yet supported
func (set *threadSafeUint16Set) PowerSet() Uint16Set {
    set.RLock()
    unsafePowerSet := set.s.PowerSet().(*threadUnsafeUint16Set)
    set.RUnlock()

    ret := &threadSafeUint16Set{s: newThreadUnsafeUint16Set()}
    for subset := range unsafePowerSet.Iter() {
        unsafeSubset := subset.(*threadUnsafeUint16Set)
        ret.Add(&threadSafeUint16Set{s: *unsafeSubset})
    }
    return ret
}
*/

func (set *threadSafeUint16Set) Pop() uint16 {
	set.Lock()
	defer set.Unlock()
	return set.s.Pop()
}

/*
// Not yet supported
func (set *threadSafeUint16Set) CartesianProduct(other Uint16Set) Uint16Set {
    o := other.(*threadSafeUint16Set)

    set.RLock()
    o.RLock()

    unsafeCartProduct := set.s.CartesianProduct(&o.s).(*threadUnsafeUint16Set)
    ret := &threadSafeUint16Set{s: *unsafeCartProduct}
    set.RUnlock()
    o.RUnlock()
    return ret
}
*/

func (set *threadSafeUint16Set) ToSlice() []uint16 {
	keys := make([]uint16, 0, set.Cardinality())
	set.RLock()
	for elem := range set.s {
		keys = append(keys, elem)
	}
	set.RUnlock()
	return keys
}

func (set *threadSafeUint16Set) MarshalJSON() ([]byte, error) {
	set.RLock()
	b, err := set.s.MarshalJSON()
	set.RUnlock()

	return b, err
}

func (set *threadSafeUint16Set) UnmarshalJSON(p []byte) error {
	set.RLock()
	err := set.s.UnmarshalJSON(p)
	set.RUnlock()

	return err
}
