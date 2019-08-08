package mapsetuint64

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type threadUnsafeUint64Set map[uint64]struct{}

// An OrderedPair represents a 2-tuple of values.
type OrderedPair struct {
	First  uint64
	Second uint64
}

func newThreadUnsafeUint64Set() threadUnsafeUint64Set {
	return make(threadUnsafeUint64Set)
}

// Equal says whether two 2-tuples contain the same values in the same order.
func (pair *OrderedPair) Equal(other OrderedPair) bool {
	if pair.First == other.First &&
		pair.Second == other.Second {
		return true
	}

	return false
}

func (set *threadUnsafeUint64Set) Add(i uint64) bool {
	_, found := (*set)[i]
	if found {
		return false //False if it existed already
	}

	(*set)[i] = struct{}{}
	return true
}

func (set *threadUnsafeUint64Set) Contains(i ...uint64) bool {
	for _, val := range i {
		if _, ok := (*set)[val]; !ok {
			return false
		}
	}
	return true
}

func (set *threadUnsafeUint64Set) IsSubset(other Uint64Set) bool {
	_ = other.(*threadUnsafeUint64Set)
	if set.Cardinality() > other.Cardinality() {
		return false
	}
	for elem := range *set {
		if !other.Contains(elem) {
			return false
		}
	}
	return true
}

func (set *threadUnsafeUint64Set) IsProperSubset(other Uint64Set) bool {
	return set.IsSubset(other) && !set.Equal(other)
}

func (set *threadUnsafeUint64Set) IsSuperset(other Uint64Set) bool {
	return other.IsSubset(set)
}

func (set *threadUnsafeUint64Set) IsProperSuperset(other Uint64Set) bool {
	return set.IsSuperset(other) && !set.Equal(other)
}

func (set *threadUnsafeUint64Set) Union(other Uint64Set) Uint64Set {
	o := other.(*threadUnsafeUint64Set)

	unionedSet := newThreadUnsafeUint64Set()

	for elem := range *set {
		unionedSet.Add(elem)
	}
	for elem := range *o {
		unionedSet.Add(elem)
	}
	return &unionedSet
}

func (set *threadUnsafeUint64Set) Intersect(other Uint64Set) Uint64Set {
	o := other.(*threadUnsafeUint64Set)

	intersection := newThreadUnsafeUint64Set()
	// loop over smaller set
	if set.Cardinality() < other.Cardinality() {
		for elem := range *set {
			if other.Contains(elem) {
				intersection.Add(elem)
			}
		}
	} else {
		for elem := range *o {
			if set.Contains(elem) {
				intersection.Add(elem)
			}
		}
	}
	return &intersection
}

func (set *threadUnsafeUint64Set) Difference(other Uint64Set) Uint64Set {
	_ = other.(*threadUnsafeUint64Set)

	difference := newThreadUnsafeUint64Set()
	for elem := range *set {
		if !other.Contains(elem) {
			difference.Add(elem)
		}
	}
	return &difference
}

func (set *threadUnsafeUint64Set) SymmetricDifference(other Uint64Set) Uint64Set {
	_ = other.(*threadUnsafeUint64Set)

	aDiff := set.Difference(other)
	bDiff := other.Difference(set)
	return aDiff.Union(bDiff)
}

func (set *threadUnsafeUint64Set) Clear() {
	*set = newThreadUnsafeUint64Set()
}

func (set *threadUnsafeUint64Set) Remove(i uint64) {
	delete(*set, i)
}

func (set *threadUnsafeUint64Set) Cardinality() int {
	return len(*set)
}

func (set *threadUnsafeUint64Set) Each(cb func(uint64) bool) {
	for elem := range *set {
		if cb(elem) {
			break
		}
	}
}

func (set *threadUnsafeUint64Set) Iter() <-chan uint64 {
	ch := make(chan uint64)
	go func() {
		for elem := range *set {
			ch <- elem
		}
		close(ch)
	}()

	return ch
}

func (set *threadUnsafeUint64Set) Iterator() *Uint64Iterator {
	iterator, ch, stopCh := newUint64Iterator()

	go func() {
	L:
		for elem := range *set {
			select {
			case <-stopCh:
				break L
			case ch <- elem:
			}
		}
		close(ch)
	}()

	return iterator
}

func (set *threadUnsafeUint64Set) Equal(other Uint64Set) bool {
	_ = other.(*threadUnsafeUint64Set)

	if set.Cardinality() != other.Cardinality() {
		return false
	}
	for elem := range *set {
		if !other.Contains(elem) {
			return false
		}
	}
	return true
}

func (set *threadUnsafeUint64Set) Clone() Uint64Set {
	clonedSet := newThreadUnsafeUint64Set()
	for elem := range *set {
		clonedSet.Add(elem)
	}
	return &clonedSet
}

func (set *threadUnsafeUint64Set) String() string {
	items := make([]string, 0, len(*set))

	for elem := range *set {
		items = append(items, fmt.Sprintf("%v", elem))
	}
	return fmt.Sprintf("Set{%s}", strings.Join(items, ", "))
}

// String outputs a 2-tuple in the form "(A, B)".
func (pair OrderedPair) String() string {
	return fmt.Sprintf("(%v, %v)", pair.First, pair.Second)
}

func (set *threadUnsafeUint64Set) Pop() uint64 {
	for item := range *set {
		delete(*set, item)
		return item
	}
	return 0
}

/*
// Not yet supported
func (set *threadUnsafeUint64Set) PowerSet() Uint64Set {
	powSet := NewThreadUnsafeUint64Set()
	nullset := newThreadUnsafeUint64Set()
	powSet.Add(&nullset)

	for es := range *set {
		u := newThreadUnsafeUint64Set()
		j := powSet.Iter()
		for er := range j {
			p := newThreadUnsafeUint64Set()
			if reflect.TypeOf(er).Name() == "" {
				k := er.(*threadUnsafeUint64Set)
				for ek := range *(k) {
					p.Add(ek)
				}
			} else {
				p.Add(er)
			}
			p.Add(es)
			u.Add(&p)
		}

		powSet = powSet.Union(&u)
	}

	return powSet
}
*/

/*
// Not yet supported
func (set *threadUnsafeUint64Set) CartesianProduct(other Uint64Set) Uint64Set {
	o := other.(*threadUnsafeUint64Set)
	cartProduct := NewThreadUnsafeUint64Set()

	for i := range *set {
		for j := range *o {
			elem := OrderedPair{First: i, Second: j}
			cartProduct.Add(elem)
		}
	}

	return cartProduct
}
*/

func (set *threadUnsafeUint64Set) ToSlice() []uint64 {
	keys := make([]uint64, 0, set.Cardinality())
	for elem := range *set {
		keys = append(keys, elem)
	}

	return keys
}

// MarshalJSON creates a JSON array from the set, it marshals all elements
func (set *threadUnsafeUint64Set) MarshalJSON() ([]byte, error) {
	items := make([]string, 0, set.Cardinality())

	for elem := range *set {
		b, err := json.Marshal(elem)
		if err != nil {
			return nil, err
		}

		items = append(items, string(b))
	}

	return []byte(fmt.Sprintf("[%s]", strings.Join(items, ","))), nil
}

// UnmarshalJSON recreates a set from a JSON array, it only decodes
// primitive types. Numbers are decoded as json.Number.
func (set *threadUnsafeUint64Set) UnmarshalJSON(b []byte) error {
	var i []uint64

	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()
	err := d.Decode(&i)
	if err != nil {
		return err
	}

	for _, v := range i {
		set.Add(v)
	}

	return nil
}
