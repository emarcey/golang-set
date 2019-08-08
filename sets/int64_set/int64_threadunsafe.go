package mapsetint64

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type threadUnsafeInt64Set map[int64]struct{}

// An OrderedPair represents a 2-tuple of values.
type OrderedPair struct {
	First  int64
	Second int64
}

func newThreadUnsafeInt64Set() threadUnsafeInt64Set {
	return make(threadUnsafeInt64Set)
}

// Equal says whether two 2-tuples contain the same values in the same order.
func (pair *OrderedPair) Equal(other OrderedPair) bool {
	if pair.First == other.First &&
		pair.Second == other.Second {
		return true
	}

	return false
}

func (set *threadUnsafeInt64Set) Add(i int64) bool {
	_, found := (*set)[i]
	if found {
		return false //False if it existed already
	}

	(*set)[i] = struct{}{}
	return true
}

func (set *threadUnsafeInt64Set) Contains(i ...int64) bool {
	for _, val := range i {
		if _, ok := (*set)[val]; !ok {
			return false
		}
	}
	return true
}

func (set *threadUnsafeInt64Set) IsSubset(other Int64Set) bool {
	_ = other.(*threadUnsafeInt64Set)
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

func (set *threadUnsafeInt64Set) IsProperSubset(other Int64Set) bool {
	return set.IsSubset(other) && !set.Equal(other)
}

func (set *threadUnsafeInt64Set) IsSuperset(other Int64Set) bool {
	return other.IsSubset(set)
}

func (set *threadUnsafeInt64Set) IsProperSuperset(other Int64Set) bool {
	return set.IsSuperset(other) && !set.Equal(other)
}

func (set *threadUnsafeInt64Set) Union(other Int64Set) Int64Set {
	o := other.(*threadUnsafeInt64Set)

	unionedSet := newThreadUnsafeInt64Set()

	for elem := range *set {
		unionedSet.Add(elem)
	}
	for elem := range *o {
		unionedSet.Add(elem)
	}
	return &unionedSet
}

func (set *threadUnsafeInt64Set) Intersect(other Int64Set) Int64Set {
	o := other.(*threadUnsafeInt64Set)

	intersection := newThreadUnsafeInt64Set()
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

func (set *threadUnsafeInt64Set) Difference(other Int64Set) Int64Set {
	_ = other.(*threadUnsafeInt64Set)

	difference := newThreadUnsafeInt64Set()
	for elem := range *set {
		if !other.Contains(elem) {
			difference.Add(elem)
		}
	}
	return &difference
}

func (set *threadUnsafeInt64Set) SymmetricDifference(other Int64Set) Int64Set {
	_ = other.(*threadUnsafeInt64Set)

	aDiff := set.Difference(other)
	bDiff := other.Difference(set)
	return aDiff.Union(bDiff)
}

func (set *threadUnsafeInt64Set) Clear() {
	*set = newThreadUnsafeInt64Set()
}

func (set *threadUnsafeInt64Set) Remove(i int64) {
	delete(*set, i)
}

func (set *threadUnsafeInt64Set) Cardinality() int {
	return len(*set)
}

func (set *threadUnsafeInt64Set) Each(cb func(int64) bool) {
	for elem := range *set {
		if cb(elem) {
			break
		}
	}
}

func (set *threadUnsafeInt64Set) Iter() <-chan int64 {
	ch := make(chan int64)
	go func() {
		for elem := range *set {
			ch <- elem
		}
		close(ch)
	}()

	return ch
}

func (set *threadUnsafeInt64Set) Iterator() *Int64Iterator {
	iterator, ch, stopCh := newInt64Iterator()

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

func (set *threadUnsafeInt64Set) Equal(other Int64Set) bool {
	_ = other.(*threadUnsafeInt64Set)

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

func (set *threadUnsafeInt64Set) Clone() Int64Set {
	clonedSet := newThreadUnsafeInt64Set()
	for elem := range *set {
		clonedSet.Add(elem)
	}
	return &clonedSet
}

func (set *threadUnsafeInt64Set) String() string {
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

func (set *threadUnsafeInt64Set) Pop() int64 {
	for item := range *set {
		delete(*set, item)
		return item
	}
	return 0
}

/*
// Not yet supported
func (set *threadUnsafeInt64Set) PowerSet() Int64Set {
	powSet := NewThreadUnsafeInt64Set()
	nullset := newThreadUnsafeInt64Set()
	powSet.Add(&nullset)

	for es := range *set {
		u := newThreadUnsafeInt64Set()
		j := powSet.Iter()
		for er := range j {
			p := newThreadUnsafeInt64Set()
			if reflect.TypeOf(er).Name() == "" {
				k := er.(*threadUnsafeInt64Set)
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
func (set *threadUnsafeInt64Set) CartesianProduct(other Int64Set) Int64Set {
	o := other.(*threadUnsafeInt64Set)
	cartProduct := NewThreadUnsafeInt64Set()

	for i := range *set {
		for j := range *o {
			elem := OrderedPair{First: i, Second: j}
			cartProduct.Add(elem)
		}
	}

	return cartProduct
}
*/

func (set *threadUnsafeInt64Set) ToSlice() []int64 {
	keys := make([]int64, 0, set.Cardinality())
	for elem := range *set {
		keys = append(keys, elem)
	}

	return keys
}

// MarshalJSON creates a JSON array from the set, it marshals all elements
func (set *threadUnsafeInt64Set) MarshalJSON() ([]byte, error) {
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
func (set *threadUnsafeInt64Set) UnmarshalJSON(b []byte) error {
	var i []int64

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
