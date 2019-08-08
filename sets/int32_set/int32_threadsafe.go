package mapsetint32

import (
	"sync"
)

type threadSafeInt32Set struct {
	s threadUnsafeInt32Set
	sync.RWMutex
}

func newThreadSafeInt32Set() threadSafeInt32Set {
	return threadSafeInt32Set{s: newThreadUnsafeInt32Set()}
}

func (set *threadSafeInt32Set) Add(i int32) bool {
	set.Lock()
	ret := set.s.Add(i)
	set.Unlock()
	return ret
}

func (set *threadSafeInt32Set) Contains(i ...int32) bool {
	set.RLock()
	ret := set.s.Contains(i...)
	set.RUnlock()
	return ret
}

func (set *threadSafeInt32Set) IsSubset(other Int32Set) bool {
	o := other.(*threadSafeInt32Set)

	set.RLock()
	o.RLock()

	ret := set.s.IsSubset(&o.s)
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeInt32Set) IsProperSubset(other Int32Set) bool {
	o := other.(*threadSafeInt32Set)

	set.RLock()
	defer set.RUnlock()
	o.RLock()
	defer o.RUnlock()

	return set.s.IsProperSubset(&o.s)
}

func (set *threadSafeInt32Set) IsSuperset(other Int32Set) bool {
	return other.IsSubset(set)
}

func (set *threadSafeInt32Set) IsProperSuperset(other Int32Set) bool {
	return other.IsProperSubset(set)
}

func (set *threadSafeInt32Set) Union(other Int32Set) Int32Set {
	o := other.(*threadSafeInt32Set)

	set.RLock()
	o.RLock()

	unsafeUnion := set.s.Union(&o.s).(*threadUnsafeInt32Set)
	ret := &threadSafeInt32Set{s: *unsafeUnion}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeInt32Set) Intersect(other Int32Set) Int32Set {
	o := other.(*threadSafeInt32Set)

	set.RLock()
	o.RLock()

	unsafeIntersection := set.s.Intersect(&o.s).(*threadUnsafeInt32Set)
	ret := &threadSafeInt32Set{s: *unsafeIntersection}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeInt32Set) Difference(other Int32Set) Int32Set {
	o := other.(*threadSafeInt32Set)

	set.RLock()
	o.RLock()

	unsafeDifference := set.s.Difference(&o.s).(*threadUnsafeInt32Set)
	ret := &threadSafeInt32Set{s: *unsafeDifference}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeInt32Set) SymmetricDifference(other Int32Set) Int32Set {
	o := other.(*threadSafeInt32Set)

	set.RLock()
	o.RLock()

	unsafeDifference := set.s.SymmetricDifference(&o.s).(*threadUnsafeInt32Set)
	ret := &threadSafeInt32Set{s: *unsafeDifference}
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeInt32Set) Clear() {
	set.Lock()
	set.s = newThreadUnsafeInt32Set()
	set.Unlock()
}

func (set *threadSafeInt32Set) Remove(i int32) {
	set.Lock()
	delete(set.s, i)
	set.Unlock()
}

func (set *threadSafeInt32Set) Cardinality() int {
	set.RLock()
	defer set.RUnlock()
	return len(set.s)
}

func (set *threadSafeInt32Set) Each(cb func(int32) bool) {
	set.RLock()
	for elem := range set.s {
		if cb(elem) {
			break
		}
	}
	set.RUnlock()
}

func (set *threadSafeInt32Set) Iter() <-chan int32 {
	ch := make(chan int32)
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

func (set *threadSafeInt32Set) Iterator() *Int32Iterator {
	iterator, ch, stopCh := newInt32Iterator()

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

func (set *threadSafeInt32Set) Equal(other Int32Set) bool {
	o := other.(*threadSafeInt32Set)

	set.RLock()
	o.RLock()

	ret := set.s.Equal(&o.s)
	set.RUnlock()
	o.RUnlock()
	return ret
}

func (set *threadSafeInt32Set) Clone() Int32Set {
	set.RLock()

	unsafeClone := set.s.Clone().(*threadUnsafeInt32Set)
	ret := &threadSafeInt32Set{s: *unsafeClone}
	set.RUnlock()
	return ret
}

func (set *threadSafeInt32Set) String() string {
	set.RLock()
	ret := set.s.String()
	set.RUnlock()
	return ret
}

/*
// Not yet supported
func (set *threadSafeInt32Set) PowerSet() Int32Set {
    set.RLock()
    unsafePowerSet := set.s.PowerSet().(*threadUnsafeInt32Set)
    set.RUnlock()

    ret := &threadSafeInt32Set{s: newThreadUnsafeInt32Set()}
    for subset := range unsafePowerSet.Iter() {
        unsafeSubset := subset.(*threadUnsafeInt32Set)
        ret.Add(&threadSafeInt32Set{s: *unsafeSubset})
    }
    return ret
}
*/

func (set *threadSafeInt32Set) Pop() int32 {
	set.Lock()
	defer set.Unlock()
	return set.s.Pop()
}

/*
// Not yet supported
func (set *threadSafeInt32Set) CartesianProduct(other Int32Set) Int32Set {
    o := other.(*threadSafeInt32Set)

    set.RLock()
    o.RLock()

    unsafeCartProduct := set.s.CartesianProduct(&o.s).(*threadUnsafeInt32Set)
    ret := &threadSafeInt32Set{s: *unsafeCartProduct}
    set.RUnlock()
    o.RUnlock()
    return ret
}
*/

func (set *threadSafeInt32Set) ToSlice() []int32 {
	keys := make([]int32, 0, set.Cardinality())
	set.RLock()
	for elem := range set.s {
		keys = append(keys, elem)
	}
	set.RUnlock()
	return keys
}

func (set *threadSafeInt32Set) MarshalJSON() ([]byte, error) {
	set.RLock()
	b, err := set.s.MarshalJSON()
	set.RUnlock()

	return b, err
}

func (set *threadSafeInt32Set) UnmarshalJSON(p []byte) error {
	set.RLock()
	err := set.s.UnmarshalJSON(p)
	set.RUnlock()

	return err
}
