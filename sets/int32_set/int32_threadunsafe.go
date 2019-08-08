package mapsetint32

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type threadUnsafeInt32Set map[int32]struct{}

// An OrderedPair represents a 2-tuple of values.
type OrderedPair struct {
	First  int32
	Second int32
}

func newThreadUnsafeInt32Set() threadUnsafeInt32Set {
	return make(threadUnsafeInt32Set)
}

// Equal says whether two 2-tuples contain the same values in the same order.
func (pair *OrderedPair) Equal(other OrderedPair) bool {
	if pair.First == other.First &&
		pair.Second == other.Second {
		return true
	}

	return false
}

func (set *threadUnsafeInt32Set) Add(i int32) bool {
	_, found := (*set)[i]
	if found {
		return false //False if it existed already
	}

	(*set)[i] = struct{}{}
	return true
}

func (set *threadUnsafeInt32Set) Contains(i ...int32) bool {
	for _, val := range i {
		if _, ok := (*set)[val]; !ok {
			return false
		}
	}
	return true
}

func (set *threadUnsafeInt32Set) IsSubset(other Int32Set) bool {
	_ = other.(*threadUnsafeInt32Set)
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

func (set *threadUnsafeInt32Set) IsProperSubset(other Int32Set) bool {
	return set.IsSubset(other) && !set.Equal(other)
}

func (set *threadUnsafeInt32Set) IsSuperset(other Int32Set) bool {
	return other.IsSubset(set)
}

func (set *threadUnsafeInt32Set) IsProperSuperset(other Int32Set) bool {
	return set.IsSuperset(other) && !set.Equal(other)
}

func (set *threadUnsafeInt32Set) Union(other Int32Set) Int32Set {
	o := other.(*threadUnsafeInt32Set)

	unionedSet := newThreadUnsafeInt32Set()

	for elem := range *set {
		unionedSet.Add(elem)
	}
	for elem := range *o {
		unionedSet.Add(elem)
	}
	return &unionedSet
}

func (set *threadUnsafeInt32Set) Intersect(other Int32Set) Int32Set {
	o := other.(*threadUnsafeInt32Set)

	intersection := newThreadUnsafeInt32Set()
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

func (set *threadUnsafeInt32Set) Difference(other Int32Set) Int32Set {
	_ = other.(*threadUnsafeInt32Set)

	difference := newThreadUnsafeInt32Set()
	for elem := range *set {
		if !other.Contains(elem) {
			difference.Add(elem)
		}
	}
	return &difference
}

func (set *threadUnsafeInt32Set) SymmetricDifference(other Int32Set) Int32Set {
	_ = other.(*threadUnsafeInt32Set)

	aDiff := set.Difference(other)
	bDiff := other.Difference(set)
	return aDiff.Union(bDiff)
}

func (set *threadUnsafeInt32Set) Clear() {
	*set = newThreadUnsafeInt32Set()
}

func (set *threadUnsafeInt32Set) Remove(i int32) {
	delete(*set, i)
}

func (set *threadUnsafeInt32Set) Cardinality() int {
	return len(*set)
}

func (set *threadUnsafeInt32Set) Each(cb func(int32) bool) {
	for elem := range *set {
		if cb(elem) {
			break
		}
	}
}

func (set *threadUnsafeInt32Set) Iter() <-chan int32 {
	ch := make(chan int32)
	go func() {
		for elem := range *set {
			ch <- elem
		}
		close(ch)
	}()

	return ch
}

func (set *threadUnsafeInt32Set) Iterator() *Int32Iterator {
	iterator, ch, stopCh := newInt32Iterator()

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

func (set *threadUnsafeInt32Set) Equal(other Int32Set) bool {
	_ = other.(*threadUnsafeInt32Set)

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

func (set *threadUnsafeInt32Set) Clone() Int32Set {
	clonedSet := newThreadUnsafeInt32Set()
	for elem := range *set {
		clonedSet.Add(elem)
	}
	return &clonedSet
}

func (set *threadUnsafeInt32Set) String() string {
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

func (set *threadUnsafeInt32Set) Pop() int32 {
	for item := range *set {
		delete(*set, item)
		return item
	}
	return 0
}

/*
// Not yet supported
func (set *threadUnsafeInt32Set) PowerSet() Int32Set {
	powSet := NewThreadUnsafeInt32Set()
	nullset := newThreadUnsafeInt32Set()
	powSet.Add(&nullset)

	for es := range *set {
		u := newThreadUnsafeInt32Set()
		j := powSet.Iter()
		for er := range j {
			p := newThreadUnsafeInt32Set()
			if reflect.TypeOf(er).Name() == "" {
				k := er.(*threadUnsafeInt32Set)
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
func (set *threadUnsafeInt32Set) CartesianProduct(other Int32Set) Int32Set {
	o := other.(*threadUnsafeInt32Set)
	cartProduct := NewThreadUnsafeInt32Set()

	for i := range *set {
		for j := range *o {
			elem := OrderedPair{First: i, Second: j}
			cartProduct.Add(elem)
		}
	}

	return cartProduct
}
*/

func (set *threadUnsafeInt32Set) ToSlice() []int32 {
	keys := make([]int32, 0, set.Cardinality())
	for elem := range *set {
		keys = append(keys, elem)
	}

	return keys
}

// MarshalJSON creates a JSON array from the set, it marshals all elements
func (set *threadUnsafeInt32Set) MarshalJSON() ([]byte, error) {
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
func (set *threadUnsafeInt32Set) UnmarshalJSON(b []byte) error {
	var i []int32

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
