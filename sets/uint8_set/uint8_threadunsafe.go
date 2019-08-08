package mapsetuint8

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type threadUnsafeUint8Set map[uint8]struct{}

// An OrderedPair represents a 2-tuple of values.
type OrderedPair struct {
	First  uint8
	Second uint8
}

func newThreadUnsafeUint8Set() threadUnsafeUint8Set {
	return make(threadUnsafeUint8Set)
}

// Equal says whether two 2-tuples contain the same values in the same order.
func (pair *OrderedPair) Equal(other OrderedPair) bool {
	if pair.First == other.First &&
		pair.Second == other.Second {
		return true
	}

	return false
}

func (set *threadUnsafeUint8Set) Add(i uint8) bool {
	_, found := (*set)[i]
	if found {
		return false //False if it existed already
	}

	(*set)[i] = struct{}{}
	return true
}

func (set *threadUnsafeUint8Set) Contains(i ...uint8) bool {
	for _, val := range i {
		if _, ok := (*set)[val]; !ok {
			return false
		}
	}
	return true
}

func (set *threadUnsafeUint8Set) IsSubset(other Uint8Set) bool {
	_ = other.(*threadUnsafeUint8Set)
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

func (set *threadUnsafeUint8Set) IsProperSubset(other Uint8Set) bool {
	return set.IsSubset(other) && !set.Equal(other)
}

func (set *threadUnsafeUint8Set) IsSuperset(other Uint8Set) bool {
	return other.IsSubset(set)
}

func (set *threadUnsafeUint8Set) IsProperSuperset(other Uint8Set) bool {
	return set.IsSuperset(other) && !set.Equal(other)
}

func (set *threadUnsafeUint8Set) Union(other Uint8Set) Uint8Set {
	o := other.(*threadUnsafeUint8Set)

	unionedSet := newThreadUnsafeUint8Set()

	for elem := range *set {
		unionedSet.Add(elem)
	}
	for elem := range *o {
		unionedSet.Add(elem)
	}
	return &unionedSet
}

func (set *threadUnsafeUint8Set) Intersect(other Uint8Set) Uint8Set {
	o := other.(*threadUnsafeUint8Set)

	intersection := newThreadUnsafeUint8Set()
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

func (set *threadUnsafeUint8Set) Difference(other Uint8Set) Uint8Set {
	_ = other.(*threadUnsafeUint8Set)

	difference := newThreadUnsafeUint8Set()
	for elem := range *set {
		if !other.Contains(elem) {
			difference.Add(elem)
		}
	}
	return &difference
}

func (set *threadUnsafeUint8Set) SymmetricDifference(other Uint8Set) Uint8Set {
	_ = other.(*threadUnsafeUint8Set)

	aDiff := set.Difference(other)
	bDiff := other.Difference(set)
	return aDiff.Union(bDiff)
}

func (set *threadUnsafeUint8Set) Clear() {
	*set = newThreadUnsafeUint8Set()
}

func (set *threadUnsafeUint8Set) Remove(i uint8) {
	delete(*set, i)
}

func (set *threadUnsafeUint8Set) Cardinality() int {
	return len(*set)
}

func (set *threadUnsafeUint8Set) Each(cb func(uint8) bool) {
	for elem := range *set {
		if cb(elem) {
			break
		}
	}
}

func (set *threadUnsafeUint8Set) Iter() <-chan uint8 {
	ch := make(chan uint8)
	go func() {
		for elem := range *set {
			ch <- elem
		}
		close(ch)
	}()

	return ch
}

func (set *threadUnsafeUint8Set) Iterator() *Uint8Iterator {
	iterator, ch, stopCh := newUint8Iterator()

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

func (set *threadUnsafeUint8Set) Equal(other Uint8Set) bool {
	_ = other.(*threadUnsafeUint8Set)

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

func (set *threadUnsafeUint8Set) Clone() Uint8Set {
	clonedSet := newThreadUnsafeUint8Set()
	for elem := range *set {
		clonedSet.Add(elem)
	}
	return &clonedSet
}

func (set *threadUnsafeUint8Set) String() string {
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

func (set *threadUnsafeUint8Set) Pop() uint8 {
	for item := range *set {
		delete(*set, item)
		return item
	}
	return 0
}

/*
// Not yet supported
func (set *threadUnsafeUint8Set) PowerSet() Uint8Set {
	powSet := NewThreadUnsafeUint8Set()
	nullset := newThreadUnsafeUint8Set()
	powSet.Add(&nullset)

	for es := range *set {
		u := newThreadUnsafeUint8Set()
		j := powSet.Iter()
		for er := range j {
			p := newThreadUnsafeUint8Set()
			if reflect.TypeOf(er).Name() == "" {
				k := er.(*threadUnsafeUint8Set)
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
func (set *threadUnsafeUint8Set) CartesianProduct(other Uint8Set) Uint8Set {
	o := other.(*threadUnsafeUint8Set)
	cartProduct := NewThreadUnsafeUint8Set()

	for i := range *set {
		for j := range *o {
			elem := OrderedPair{First: i, Second: j}
			cartProduct.Add(elem)
		}
	}

	return cartProduct
}
*/

func (set *threadUnsafeUint8Set) ToSlice() []uint8 {
	keys := make([]uint8, 0, set.Cardinality())
	for elem := range *set {
		keys = append(keys, elem)
	}

	return keys
}

// MarshalJSON creates a JSON array from the set, it marshals all elements
func (set *threadUnsafeUint8Set) MarshalJSON() ([]byte, error) {
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
func (set *threadUnsafeUint8Set) UnmarshalJSON(b []byte) error {
	var i []uint8

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
