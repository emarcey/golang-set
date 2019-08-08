package mapsetint8

import (
	"sync"
)

type threadSafeInt8Set struct {
	s threadUnsafeInt8Set
	sync.RWMutex
}

func newThreadSafeInt8Set() threadSafeInt8Set {
	return threadSafeInt8Set{s: newThreadUnsafeInt8Set()}
}

func (set *threadSafeInt8Set) Add(i int8) bool {
	set.Lock()
	ret := set.s.Add(i)
	set.Unlock()
	return ret
}

func (set *threadSafeInt8Set) Contains(i ...int8) bool {
	set.RLock()
	ret := set.s.Contains(i...)
	set.RUnlock()
	return ret
}

func (set *threadSafeInt8Set) IsSubset(other Int8Set) bool {
	o := other.(*threadSafeInt8Set)

	set.RLock()
	o.RLock()

	ret := set.s.IsSubset(&o.s)
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeInt8Set) IsProperSubset(other Int8Set) bool {
	o := other.(*threadSafeInt8Set)

	set.RLock()
	defer set.RUnlock()
	o.RLock()
	defer o.RUnlock()

	return set.s.IsProperSubset(&o.s)
}

func (set *threadSafeInt8Set) IsSuperset(other Int8Set) bool {
	return other.IsSubset(set)
}

func (set *threadSafeInt8Set) IsProperSuperset(other Int8Set) bool {
	return other.IsProperSubset(set)
}

func (set *threadSafeInt8Set) Union(other Int8Set) Int8Set {
	o := other.(*threadSafeInt8Set)

	set.RLock()
	o.RLock()

	unsafeUnion := set.s.Union(&o.s).(*threadUnsafeInt8Set)
	ret := &threadSafeInt8Set{s: *unsafeUnion}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeInt8Set) Intersect(other Int8Set) Int8Set {
	o := other.(*threadSafeInt8Set)

	set.RLock()
	o.RLock()

	unsafeIntersection := set.s.Intersect(&o.s).(*threadUnsafeInt8Set)
	ret := &threadSafeInt8Set{s: *unsafeIntersection}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeInt8Set) Difference(other Int8Set) Int8Set {
	o := other.(*threadSafeInt8Set)

	set.RLock()
	o.RLock()

	unsafeDifference := set.s.Difference(&o.s).(*threadUnsafeInt8Set)
	ret := &threadSafeInt8Set{s: *unsafeDifference}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeInt8Set) SymmetricDifference(other Int8Set) Int8Set {
	o := other.(*threadSafeInt8Set)

	set.RLock()
	o.RLock()

	unsafeDifference := set.s.SymmetricDifference(&o.s).(*threadUnsafeInt8Set)
	ret := &threadSafeInt8Set{s: *unsafeDifference}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeInt8Set) Clear() {
	set.Lock()
	set.s = newThreadUnsafeInt8Set()
	set.Unlock()
}

func (set *threadSafeInt8Set) Remove(i int8) {
	set.Lock()
	delete(set.s, i)
	set.Unlock()
}

func (set *threadSafeInt8Set) Cardinality() int {
	set.RLock()
	defer set.RUnlock()
	return len(set.s)
}

func (set *threadSafeInt8Set) Each(cb func(int8) bool) {
	set.RLock()
	for elem := range set.s {
		if cb(elem) {
			break
		}
	}
	set.RUnlock()
}

func (set *threadSafeInt8Set) Iter() <-chan int8 {
	ch := make(chan int8)
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

func (set *threadSafeInt8Set) Iterator() *Int8Iterator {
	iterator, ch, stopCh := newInt8Iterator()

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

func (set *threadSafeInt8Set) Equal(other Int8Set) bool {
	o := other.(*threadSafeInt8Set)

	set.RLock()
	o.RLock()

	ret := set.s.Equal(&o.s)
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeInt8Set) Clone() Int8Set {
	set.RLock()

	unsafeClone := set.s.Clone().(*threadUnsafeInt8Set)
	ret := &threadSafeInt8Set{s: *unsafeClone}
	set.RUnlock()
	return ret
}

func (set *threadSafeInt8Set) String() string {
	set.RLock()
	ret := set.s.String()
	set.RUnlock()
	return ret
}

/*
// Not yet supported
func (set *threadSafeInt8Set) PowerSet() Int8Set {
    set.RLock()
    unsafePowerSet := set.s.PowerSet().(*threadUnsafeInt8Set)
    set.RUnlock()

    ret := &threadSafeInt8Set{s: newThreadUnsafeInt8Set()}
    for subset := range unsafePowerSet.Iter() {
        unsafeSubset := subset.(*threadUnsafeInt8Set)
        ret.Add(&threadSafeInt8Set{s: *unsafeSubset})
    }
    return ret
}
*/

func (set *threadSafeInt8Set) Pop() int8 {
	set.Lock()
	defer set.Unlock()
	return set.s.Pop()
}

/*
// Not yet supported
func (set *threadSafeInt8Set) CartesianProduct(other Int8Set) Int8Set {
    o := other.(*threadSafeInt8Set)

    set.RLock()
    o.RLock()

    unsafeCartProduct := set.s.CartesianProduct(&o.s).(*threadUnsafeInt8Set)
    ret := &threadSafeInt8Set{s: *unsafeCartProduct}
    set.RUnlock()
    o.RUnlock()
    return ret
}
*/

func (set *threadSafeInt8Set) ToSlice() []int8 {
	keys := make([]int8, 0, set.Cardinality())
	set.RLock()
	for elem := range set.s {
		keys = append(keys, elem)
	}
	set.RUnlock()
	return keys
}

func (set *threadSafeInt8Set) MarshalJSON() ([]byte, error) {
	set.RLock()
	b, err := set.s.MarshalJSON()
	set.RUnlock()

	return b, err
}

func (set *threadSafeInt8Set) UnmarshalJSON(p []byte) error {
	set.RLock()
	err := set.s.UnmarshalJSON(p)
	set.RUnlock()

	return err
}
