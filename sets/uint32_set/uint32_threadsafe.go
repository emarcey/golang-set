package mapsetuint32

import (
	"sync"
)

type threadSafeUint32Set struct {
	s threadUnsafeUint32Set
	sync.RWMutex
}

func newThreadSafeUint32Set() threadSafeUint32Set {
	return threadSafeUint32Set{s: newThreadUnsafeUint32Set()}
}

func (set *threadSafeUint32Set) Add(i uint32) bool {
	set.Lock()
	ret := set.s.Add(i)
	set.Unlock()
	return ret
}

func (set *threadSafeUint32Set) Contains(i ...uint32) bool {
	set.RLock()
	ret := set.s.Contains(i...)
	set.RUnlock()
	return ret
}

func (set *threadSafeUint32Set) IsSubset(other Uint32Set) bool {
	o := other.(*threadSafeUint32Set)

	set.RLock()
	o.RLock()

	ret := set.s.IsSubset(&o.s)
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeUint32Set) IsProperSubset(other Uint32Set) bool {
	o := other.(*threadSafeUint32Set)

	set.RLock()
	defer set.RUnlock()
	o.RLock()
	defer o.RUnlock()

	return set.s.IsProperSubset(&o.s)
}

func (set *threadSafeUint32Set) IsSuperset(other Uint32Set) bool {
	return other.IsSubset(set)
}

func (set *threadSafeUint32Set) IsProperSuperset(other Uint32Set) bool {
	return other.IsProperSubset(set)
}

func (set *threadSafeUint32Set) Union(other Uint32Set) Uint32Set {
	o := other.(*threadSafeUint32Set)

	set.RLock()
	o.RLock()

	unsafeUnion := set.s.Union(&o.s).(*threadUnsafeUint32Set)
	ret := &threadSafeUint32Set{s: *unsafeUnion}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeUint32Set) Intersect(other Uint32Set) Uint32Set {
	o := other.(*threadSafeUint32Set)

	set.RLock()
	o.RLock()

	unsafeIntersection := set.s.Intersect(&o.s).(*threadUnsafeUint32Set)
	ret := &threadSafeUint32Set{s: *unsafeIntersection}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeUint32Set) Difference(other Uint32Set) Uint32Set {
	o := other.(*threadSafeUint32Set)

	set.RLock()
	o.RLock()

	unsafeDifference := set.s.Difference(&o.s).(*threadUnsafeUint32Set)
	ret := &threadSafeUint32Set{s: *unsafeDifference}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeUint32Set) SymmetricDifference(other Uint32Set) Uint32Set {
	o := other.(*threadSafeUint32Set)

	set.RLock()
	o.RLock()

	unsafeDifference := set.s.SymmetricDifference(&o.s).(*threadUnsafeUint32Set)
	ret := &threadSafeUint32Set{s: *unsafeDifference}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeUint32Set) Clear() {
	set.Lock()
	set.s = newThreadUnsafeUint32Set()
	set.Unlock()
}

func (set *threadSafeUint32Set) Remove(i uint32) {
	set.Lock()
	delete(set.s, i)
	set.Unlock()
}

func (set *threadSafeUint32Set) Cardinality() int {
	set.RLock()
	defer set.RUnlock()
	return len(set.s)
}

func (set *threadSafeUint32Set) Each(cb func(uint32) bool) {
	set.RLock()
	for elem := range set.s {
		if cb(elem) {
			break
		}
	}
	set.RUnlock()
}

func (set *threadSafeUint32Set) Iter() <-chan uint32 {
	ch := make(chan uint32)
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

func (set *threadSafeUint32Set) Iterator() *Uint32Iterator {
	iterator, ch, stopCh := newUint32Iterator()

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

func (set *threadSafeUint32Set) Equal(other Uint32Set) bool {
	o := other.(*threadSafeUint32Set)

	set.RLock()
	o.RLock()

	ret := set.s.Equal(&o.s)
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeUint32Set) Clone() Uint32Set {
	set.RLock()

	unsafeClone := set.s.Clone().(*threadUnsafeUint32Set)
	ret := &threadSafeUint32Set{s: *unsafeClone}
	set.RUnlock()
	return ret
}

func (set *threadSafeUint32Set) String() string {
	set.RLock()
	ret := set.s.String()
	set.RUnlock()
	return ret
}

/*
// Not yet supported
func (set *threadSafeUint32Set) PowerSet() Uint32Set {
    set.RLock()
    unsafePowerSet := set.s.PowerSet().(*threadUnsafeUint32Set)
    set.RUnlock()

    ret := &threadSafeUint32Set{s: newThreadUnsafeUint32Set()}
    for subset := range unsafePowerSet.Iter() {
        unsafeSubset := subset.(*threadUnsafeUint32Set)
        ret.Add(&threadSafeUint32Set{s: *unsafeSubset})
    }
    return ret
}
*/

func (set *threadSafeUint32Set) Pop() uint32 {
	set.Lock()
	defer set.Unlock()
	return set.s.Pop()
}

/*
// Not yet supported
func (set *threadSafeUint32Set) CartesianProduct(other Uint32Set) Uint32Set {
    o := other.(*threadSafeUint32Set)

    set.RLock()
    o.RLock()

    unsafeCartProduct := set.s.CartesianProduct(&o.s).(*threadUnsafeUint32Set)
    ret := &threadSafeUint32Set{s: *unsafeCartProduct}
    set.RUnlock()
    o.RUnlock()
    return ret
}
*/

func (set *threadSafeUint32Set) ToSlice() []uint32 {
	keys := make([]uint32, 0, set.Cardinality())
	set.RLock()
	for elem := range set.s {
		keys = append(keys, elem)
	}
	set.RUnlock()
	return keys
}

func (set *threadSafeUint32Set) MarshalJSON() ([]byte, error) {
	set.RLock()
	b, err := set.s.MarshalJSON()
	set.RUnlock()

	return b, err
}

func (set *threadSafeUint32Set) UnmarshalJSON(p []byte) error {
	set.RLock()
	err := set.s.UnmarshalJSON(p)
	set.RUnlock()

	return err
}
