package mapsetint16

import (
	"sync"
)

type threadSafeInt16Set struct {
	s threadUnsafeInt16Set
	sync.RWMutex
}

func newThreadSafeInt16Set() threadSafeInt16Set {
	return threadSafeInt16Set{s: newThreadUnsafeInt16Set()}
}

func (set *threadSafeInt16Set) Add(i int16) bool {
	set.Lock()
	ret := set.s.Add(i)
	set.Unlock()
	return ret
}

func (set *threadSafeInt16Set) Contains(i ...int16) bool {
	set.RLock()
	ret := set.s.Contains(i...)
	set.RUnlock()
	return ret
}

func (set *threadSafeInt16Set) IsSubset(other Int16Set) bool {
	o := other.(*threadSafeInt16Set)

	set.RLock()
	o.RLock()

	ret := set.s.IsSubset(&o.s)
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeInt16Set) IsProperSubset(other Int16Set) bool {
	o := other.(*threadSafeInt16Set)

	set.RLock()
	defer set.RUnlock()
	o.RLock()
	defer o.RUnlock()

	return set.s.IsProperSubset(&o.s)
}

func (set *threadSafeInt16Set) IsSuperset(other Int16Set) bool {
	return other.IsSubset(set)
}

func (set *threadSafeInt16Set) IsProperSuperset(other Int16Set) bool {
	return other.IsProperSubset(set)
}

func (set *threadSafeInt16Set) Union(other Int16Set) Int16Set {
	o := other.(*threadSafeInt16Set)

	set.RLock()
	o.RLock()

	unsafeUnion := set.s.Union(&o.s).(*threadUnsafeInt16Set)
	ret := &threadSafeInt16Set{s: *unsafeUnion}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeInt16Set) Intersect(other Int16Set) Int16Set {
	o := other.(*threadSafeInt16Set)

	set.RLock()
	o.RLock()

	unsafeIntersection := set.s.Intersect(&o.s).(*threadUnsafeInt16Set)
	ret := &threadSafeInt16Set{s: *unsafeIntersection}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeInt16Set) Difference(other Int16Set) Int16Set {
	o := other.(*threadSafeInt16Set)

	set.RLock()
	o.RLock()

	unsafeDifference := set.s.Difference(&o.s).(*threadUnsafeInt16Set)
	ret := &threadSafeInt16Set{s: *unsafeDifference}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeInt16Set) SymmetricDifference(other Int16Set) Int16Set {
	o := other.(*threadSafeInt16Set)

	set.RLock()
	o.RLock()

	unsafeDifference := set.s.SymmetricDifference(&o.s).(*threadUnsafeInt16Set)
	ret := &threadSafeInt16Set{s: *unsafeDifference}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeInt16Set) Clear() {
	set.Lock()
	set.s = newThreadUnsafeInt16Set()
	set.Unlock()
}

func (set *threadSafeInt16Set) Remove(i int16) {
	set.Lock()
	delete(set.s, i)
	set.Unlock()
}

func (set *threadSafeInt16Set) Cardinality() int {
	set.RLock()
	defer set.RUnlock()
	return len(set.s)
}

func (set *threadSafeInt16Set) Each(cb func(int16) bool) {
	set.RLock()
	for elem := range set.s {
		if cb(elem) {
			break
		}
	}
	set.RUnlock()
}

func (set *threadSafeInt16Set) Iter() <-chan int16 {
	ch := make(chan int16)
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

func (set *threadSafeInt16Set) Iterator() *Int16Iterator {
	iterator, ch, stopCh := newInt16Iterator()

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

func (set *threadSafeInt16Set) Equal(other Int16Set) bool {
	o := other.(*threadSafeInt16Set)

	set.RLock()
	o.RLock()

	ret := set.s.Equal(&o.s)
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeInt16Set) Clone() Int16Set {
	set.RLock()

	unsafeClone := set.s.Clone().(*threadUnsafeInt16Set)
	ret := &threadSafeInt16Set{s: *unsafeClone}
	set.RUnlock()
	return ret
}

func (set *threadSafeInt16Set) String() string {
	set.RLock()
	ret := set.s.String()
	set.RUnlock()
	return ret
}

/*
// Not yet supported
func (set *threadSafeInt16Set) PowerSet() Int16Set {
    set.RLock()
    unsafePowerSet := set.s.PowerSet().(*threadUnsafeInt16Set)
    set.RUnlock()

    ret := &threadSafeInt16Set{s: newThreadUnsafeInt16Set()}
    for subset := range unsafePowerSet.Iter() {
        unsafeSubset := subset.(*threadUnsafeInt16Set)
        ret.Add(&threadSafeInt16Set{s: *unsafeSubset})
    }
    return ret
}
*/

func (set *threadSafeInt16Set) Pop() int16 {
	set.Lock()
	defer set.Unlock()
	return set.s.Pop()
}

/*
// Not yet supported
func (set *threadSafeInt16Set) CartesianProduct(other Int16Set) Int16Set {
    o := other.(*threadSafeInt16Set)

    set.RLock()
    o.RLock()

    unsafeCartProduct := set.s.CartesianProduct(&o.s).(*threadUnsafeInt16Set)
    ret := &threadSafeInt16Set{s: *unsafeCartProduct}
    set.RUnlock()
    o.RUnlock()
    return ret
}
*/

func (set *threadSafeInt16Set) ToSlice() []int16 {
	keys := make([]int16, 0, set.Cardinality())
	set.RLock()
	for elem := range set.s {
		keys = append(keys, elem)
	}
	set.RUnlock()
	return keys
}

func (set *threadSafeInt16Set) MarshalJSON() ([]byte, error) {
	set.RLock()
	b, err := set.s.MarshalJSON()
	set.RUnlock()

	return b, err
}

func (set *threadSafeInt16Set) UnmarshalJSON(p []byte) error {
	set.RLock()
	err := set.s.UnmarshalJSON(p)
	set.RUnlock()

	return err
}
