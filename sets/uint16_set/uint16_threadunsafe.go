package mapsetuint16

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type threadUnsafeUint16Set map[uint16]struct{}

// An OrderedPair represents a 2-tuple of values.
type OrderedPair struct {
	First  uint16
	Second uint16
}

func newThreadUnsafeUint16Set() threadUnsafeUint16Set {
	return make(threadUnsafeUint16Set)
}

// Equal says whether two 2-tuples contain the same values in the same order.
func (pair *OrderedPair) Equal(other OrderedPair) bool {
	if pair.First == other.First &&
		pair.Second == other.Second {
		return true
	}

	return false
}

func (set *threadUnsafeUint16Set) Add(i uint16) bool {
	_, found := (*set)[i]
	if found {
		return false //False if it existed already
	}

	(*set)[i] = struct{}{}
	return true
}

func (set *threadUnsafeUint16Set) Contains(i ...uint16) bool {
	for _, val := range i {
		if _, ok := (*set)[val]; !ok {
			return false
		}
	}
	return true
}

func (set *threadUnsafeUint16Set) IsSubset(other Uint16Set) bool {
	_ = other.(*threadUnsafeUint16Set)
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

func (set *threadUnsafeUint16Set) IsProperSubset(other Uint16Set) bool {
	return set.IsSubset(other) && !set.Equal(other)
}

func (set *threadUnsafeUint16Set) IsSuperset(other Uint16Set) bool {
	return other.IsSubset(set)
}

func (set *threadUnsafeUint16Set) IsProperSuperset(other Uint16Set) bool {
	return set.IsSuperset(other) && !set.Equal(other)
}

func (set *threadUnsafeUint16Set) Union(other Uint16Set) Uint16Set {
	o := other.(*threadUnsafeUint16Set)

	unionedSet := newThreadUnsafeUint16Set()

	for elem := range *set {
		unionedSet.Add(elem)
	}
	for elem := range *o {
		unionedSet.Add(elem)
	}
	return &unionedSet
}

func (set *threadUnsafeUint16Set) Intersect(other Uint16Set) Uint16Set {
	o := other.(*threadUnsafeUint16Set)

	intersection := newThreadUnsafeUint16Set()
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

func (set *threadUnsafeUint16Set) Difference(other Uint16Set) Uint16Set {
	_ = other.(*threadUnsafeUint16Set)

	difference := newThreadUnsafeUint16Set()
	for elem := range *set {
		if !other.Contains(elem) {
			difference.Add(elem)
		}
	}
	return &difference
}

func (set *threadUnsafeUint16Set) SymmetricDifference(other Uint16Set) Uint16Set {
	_ = other.(*threadUnsafeUint16Set)

	aDiff := set.Difference(other)
	bDiff := other.Difference(set)
	return aDiff.Union(bDiff)
}

func (set *threadUnsafeUint16Set) Clear() {
	*set = newThreadUnsafeUint16Set()
}

func (set *threadUnsafeUint16Set) Remove(i uint16) {
	delete(*set, i)
}

func (set *threadUnsafeUint16Set) Cardinality() int {
	return len(*set)
}

func (set *threadUnsafeUint16Set) Each(cb func(uint16) bool) {
	for elem := range *set {
		if cb(elem) {
			break
		}
	}
}

func (set *threadUnsafeUint16Set) Iter() <-chan uint16 {
	ch := make(chan uint16)
	go func() {
		for elem := range *set {
			ch <- elem
		}
		close(ch)
	}()

	return ch
}

func (set *threadUnsafeUint16Set) Iterator() *Uint16Iterator {
	iterator, ch, stopCh := newUint16Iterator()

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

func (set *threadUnsafeUint16Set) Equal(other Uint16Set) bool {
	_ = other.(*threadUnsafeUint16Set)

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

func (set *threadUnsafeUint16Set) Clone() Uint16Set {
	clonedSet := newThreadUnsafeUint16Set()
	for elem := range *set {
		clonedSet.Add(elem)
	}
	return &clonedSet
}

func (set *threadUnsafeUint16Set) String() string {
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

func (set *threadUnsafeUint16Set) Pop() uint16 {
	for item := range *set {
		delete(*set, item)
		return item
	}
	return 0
}

/*
// Not yet supported
func (set *threadUnsafeUint16Set) PowerSet() Uint16Set {
	powSet := NewThreadUnsafeUint16Set()
	nullset := newThreadUnsafeUint16Set()
	powSet.Add(&nullset)

	for es := range *set {
		u := newThreadUnsafeUint16Set()
		j := powSet.Iter()
		for er := range j {
			p := newThreadUnsafeUint16Set()
			if reflect.TypeOf(er).Name() == "" {
				k := er.(*threadUnsafeUint16Set)
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
func (set *threadUnsafeUint16Set) CartesianProduct(other Uint16Set) Uint16Set {
	o := other.(*threadUnsafeUint16Set)
	cartProduct := NewThreadUnsafeUint16Set()

	for i := range *set {
		for j := range *o {
			elem := OrderedPair{First: i, Second: j}
			cartProduct.Add(elem)
		}
	}

	return cartProduct
}
*/

func (set *threadUnsafeUint16Set) ToSlice() []uint16 {
	keys := make([]uint16, 0, set.Cardinality())
	for elem := range *set {
		keys = append(keys, elem)
	}

	return keys
}

// MarshalJSON creates a JSON array from the set, it marshals all elements
func (set *threadUnsafeUint16Set) MarshalJSON() ([]byte, error) {
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
func (set *threadUnsafeUint16Set) UnmarshalJSON(b []byte) error {
	var i []uint16

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
