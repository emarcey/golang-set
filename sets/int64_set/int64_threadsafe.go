package mapsetint64

import (
	"sync"
)

type threadSafeInt64Set struct {
	s threadUnsafeInt64Set
	sync.RWMutex
}

func newThreadSafeInt64Set() threadSafeInt64Set {
	return threadSafeInt64Set{s: newThreadUnsafeInt64Set()}
}

func (set *threadSafeInt64Set) Add(i int64) bool {
	set.Lock()
	ret := set.s.Add(i)
	set.Unlock()
	return ret
}

func (set *threadSafeInt64Set) Contains(i ...int64) bool {
	set.RLock()
	ret := set.s.Contains(i...)
	set.RUnlock()
	return ret
}

func (set *threadSafeInt64Set) IsSubset(other Int64Set) bool {
	o := other.(*threadSafeInt64Set)

	set.RLock()
	o.RLock()

	ret := set.s.IsSubset(&o.s)
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeInt64Set) IsProperSubset(other Int64Set) bool {
	o := other.(*threadSafeInt64Set)

	set.RLock()
	defer set.RUnlock()
	o.RLock()
	defer o.RUnlock()

	return set.s.IsProperSubset(&o.s)
}

func (set *threadSafeInt64Set) IsSuperset(other Int64Set) bool {
	return other.IsSubset(set)
}

func (set *threadSafeInt64Set) IsProperSuperset(other Int64Set) bool {
	return other.IsProperSubset(set)
}

func (set *threadSafeInt64Set) Union(other Int64Set) Int64Set {
	o := other.(*threadSafeInt64Set)

	set.RLock()
	o.RLock()

	unsafeUnion := set.s.Union(&o.s).(*threadUnsafeInt64Set)
	ret := &threadSafeInt64Set{s: *unsafeUnion}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeInt64Set) Intersect(other Int64Set) Int64Set {
	o := other.(*threadSafeInt64Set)

	set.RLock()
	o.RLock()

	unsafeIntersection := set.s.Intersect(&o.s).(*threadUnsafeInt64Set)
	ret := &threadSafeInt64Set{s: *unsafeIntersection}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeInt64Set) Difference(other Int64Set) Int64Set {
	o := other.(*threadSafeInt64Set)

	set.RLock()
	o.RLock()

	unsafeDifference := set.s.Difference(&o.s).(*threadUnsafeInt64Set)
	ret := &threadSafeInt64Set{s: *unsafeDifference}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeInt64Set) SymmetricDifference(other Int64Set) Int64Set {
	o := other.(*threadSafeInt64Set)

	set.RLock()
	o.RLock()

	unsafeDifference := set.s.SymmetricDifference(&o.s).(*threadUnsafeInt64Set)
	ret := &threadSafeInt64Set{s: *unsafeDifference}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeInt64Set) Clear() {
	set.Lock()
	set.s = newThreadUnsafeInt64Set()
	set.Unlock()
}

func (set *threadSafeInt64Set) Remove(i int64) {
	set.Lock()
	delete(set.s, i)
	set.Unlock()
}

func (set *threadSafeInt64Set) Cardinality() int {
	set.RLock()
	defer set.RUnlock()
	return len(set.s)
}

func (set *threadSafeInt64Set) Each(cb func(int64) bool) {
	set.RLock()
	for elem := range set.s {
		if cb(elem) {
			break
		}
	}
	set.RUnlock()
}

func (set *threadSafeInt64Set) Iter() <-chan int64 {
	ch := make(chan int64)
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

func (set *threadSafeInt64Set) Iterator() *Int64Iterator {
	iterator, ch, stopCh := newInt64Iterator()

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

func (set *threadSafeInt64Set) Equal(other Int64Set) bool {
	o := other.(*threadSafeInt64Set)

	set.RLock()
	o.RLock()

	ret := set.s.Equal(&o.s)
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeInt64Set) Clone() Int64Set {
	set.RLock()

	unsafeClone := set.s.Clone().(*threadUnsafeInt64Set)
	ret := &threadSafeInt64Set{s: *unsafeClone}
	set.RUnlock()
	return ret
}

func (set *threadSafeInt64Set) String() string {
	set.RLock()
	ret := set.s.String()
	set.RUnlock()
	return ret
}

/*
// Not yet supported
func (set *threadSafeInt64Set) PowerSet() Int64Set {
    set.RLock()
    unsafePowerSet := set.s.PowerSet().(*threadUnsafeInt64Set)
    set.RUnlock()

    ret := &threadSafeInt64Set{s: newThreadUnsafeInt64Set()}
    for subset := range unsafePowerSet.Iter() {
        unsafeSubset := subset.(*threadUnsafeInt64Set)
        ret.Add(&threadSafeInt64Set{s: *unsafeSubset})
    }
    return ret
}
*/

func (set *threadSafeInt64Set) Pop() int64 {
	set.Lock()
	defer set.Unlock()
	return set.s.Pop()
}

/*
// Not yet supported
func (set *threadSafeInt64Set) CartesianProduct(other Int64Set) Int64Set {
    o := other.(*threadSafeInt64Set)

    set.RLock()
    o.RLock()

    unsafeCartProduct := set.s.CartesianProduct(&o.s).(*threadUnsafeInt64Set)
    ret := &threadSafeInt64Set{s: *unsafeCartProduct}
    set.RUnlock()
    o.RUnlock()
    return ret
}
*/

func (set *threadSafeInt64Set) ToSlice() []int64 {
	keys := make([]int64, 0, set.Cardinality())
	set.RLock()
	for elem := range set.s {
		keys = append(keys, elem)
	}
	set.RUnlock()
	return keys
}

func (set *threadSafeInt64Set) MarshalJSON() ([]byte, error) {
	set.RLock()
	b, err := set.s.MarshalJSON()
	set.RUnlock()

	return b, err
}

func (set *threadSafeInt64Set) UnmarshalJSON(p []byte) error {
	set.RLock()
	err := set.s.UnmarshalJSON(p)
	set.RUnlock()

	return err
}
