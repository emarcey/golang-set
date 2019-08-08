package mapsetuint32

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type threadUnsafeUint32Set map[uint32]struct{}

// An OrderedPair represents a 2-tuple of values.
type OrderedPair struct {
	First  uint32
	Second uint32
}

func newThreadUnsafeUint32Set() threadUnsafeUint32Set {
	return make(threadUnsafeUint32Set)
}

// Equal says whether two 2-tuples contain the same values in the same order.
func (pair *OrderedPair) Equal(other OrderedPair) bool {
	if pair.First == other.First &&
		pair.Second == other.Second {
		return true
	}

	return false
}

func (set *threadUnsafeUint32Set) Add(i uint32) bool {
	_, found := (*set)[i]
	if found {
		return false //False if it existed already
	}

	(*set)[i] = struct{}{}
	return true
}

func (set *threadUnsafeUint32Set) Contains(i ...uint32) bool {
	for _, val := range i {
		if _, ok := (*set)[val]; !ok {
			return false
		}
	}
	return true
}

func (set *threadUnsafeUint32Set) IsSubset(other Uint32Set) bool {
	_ = other.(*threadUnsafeUint32Set)
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

func (set *threadUnsafeUint32Set) IsProperSubset(other Uint32Set) bool {
	return set.IsSubset(other) && !set.Equal(other)
}

func (set *threadUnsafeUint32Set) IsSuperset(other Uint32Set) bool {
	return other.IsSubset(set)
}

func (set *threadUnsafeUint32Set) IsProperSuperset(other Uint32Set) bool {
	return set.IsSuperset(other) && !set.Equal(other)
}

func (set *threadUnsafeUint32Set) Union(other Uint32Set) Uint32Set {
	o := other.(*threadUnsafeUint32Set)

	unionedSet := newThreadUnsafeUint32Set()

	for elem := range *set {
		unionedSet.Add(elem)
	}
	for elem := range *o {
		unionedSet.Add(elem)
	}
	return &unionedSet
}

func (set *threadUnsafeUint32Set) Intersect(other Uint32Set) Uint32Set {
	o := other.(*threadUnsafeUint32Set)

	intersection := newThreadUnsafeUint32Set()
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

func (set *threadUnsafeUint32Set) Difference(other Uint32Set) Uint32Set {
	_ = other.(*threadUnsafeUint32Set)

	difference := newThreadUnsafeUint32Set()
	for elem := range *set {
		if !other.Contains(elem) {
			difference.Add(elem)
		}
	}
	return &difference
}

func (set *threadUnsafeUint32Set) SymmetricDifference(other Uint32Set) Uint32Set {
	_ = other.(*threadUnsafeUint32Set)

	aDiff := set.Difference(other)
	bDiff := other.Difference(set)
	return aDiff.Union(bDiff)
}

func (set *threadUnsafeUint32Set) Clear() {
	*set = newThreadUnsafeUint32Set()
}

func (set *threadUnsafeUint32Set) Remove(i uint32) {
	delete(*set, i)
}

func (set *threadUnsafeUint32Set) Cardinality() int {
	return len(*set)
}

func (set *threadUnsafeUint32Set) Each(cb func(uint32) bool) {
	for elem := range *set {
		if cb(elem) {
			break
		}
	}
}

func (set *threadUnsafeUint32Set) Iter() <-chan uint32 {
	ch := make(chan uint32)
	go func() {
		for elem := range *set {
			ch <- elem
		}
		close(ch)
	}()

	return ch
}

func (set *threadUnsafeUint32Set) Iterator() *Uint32Iterator {
	iterator, ch, stopCh := newUint32Iterator()

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

func (set *threadUnsafeUint32Set) Equal(other Uint32Set) bool {
	_ = other.(*threadUnsafeUint32Set)

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

func (set *threadUnsafeUint32Set) Clone() Uint32Set {
	clonedSet := newThreadUnsafeUint32Set()
	for elem := range *set {
		clonedSet.Add(elem)
	}
	return &clonedSet
}

func (set *threadUnsafeUint32Set) String() string {
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

func (set *threadUnsafeUint32Set) Pop() uint32 {
	for item := range *set {
		delete(*set, item)
		return item
	}
	return 0
}

/*
// Not yet supported
func (set *threadUnsafeUint32Set) PowerSet() Uint32Set {
	powSet := NewThreadUnsafeUint32Set()
	nullset := newThreadUnsafeUint32Set()
	powSet.Add(&nullset)

	for es := range *set {
		u := newThreadUnsafeUint32Set()
		j := powSet.Iter()
		for er := range j {
			p := newThreadUnsafeUint32Set()
			if reflect.TypeOf(er).Name() == "" {
				k := er.(*threadUnsafeUint32Set)
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
func (set *threadUnsafeUint32Set) CartesianProduct(other Uint32Set) Uint32Set {
	o := other.(*threadUnsafeUint32Set)
	cartProduct := NewThreadUnsafeUint32Set()

	for i := range *set {
		for j := range *o {
			elem := OrderedPair{First: i, Second: j}
			cartProduct.Add(elem)
		}
	}

	return cartProduct
}
*/

func (set *threadUnsafeUint32Set) ToSlice() []uint32 {
	keys := make([]uint32, 0, set.Cardinality())
	for elem := range *set {
		keys = append(keys, elem)
	}

	return keys
}

// MarshalJSON creates a JSON array from the set, it marshals all elements
func (set *threadUnsafeUint32Set) MarshalJSON() ([]byte, error) {
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
func (set *threadUnsafeUint32Set) UnmarshalJSON(b []byte) error {
	var i []uint32

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
