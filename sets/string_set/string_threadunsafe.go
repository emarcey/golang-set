package mapsetstring

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type threadUnsafeStringSet map[string]struct{}

// An OrderedPair represents a 2-tuple of values.
type OrderedPair struct {
	First  string
	Second string
}

func newThreadUnsafeStringSet() threadUnsafeStringSet {
	return make(threadUnsafeStringSet)
}

// Equal says whether two 2-tuples contain the same values in the same order.
func (pair *OrderedPair) Equal(other OrderedPair) bool {
	if pair.First == other.First &&
		pair.Second == other.Second {
		return true
	}

	return false
}

func (set *threadUnsafeStringSet) Add(i string) bool {
	_, found := (*set)[i]
	if found {
		return false //False if it existed already
	}

	(*set)[i] = struct{}{}
	return true
}

func (set *threadUnsafeStringSet) Contains(i ...string) bool {
	for _, val := range i {
		if _, ok := (*set)[val]; !ok {
			return false
		}
	}
	return true
}

func (set *threadUnsafeStringSet) IsSubset(other StringSet) bool {
	_ = other.(*threadUnsafeStringSet)
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

func (set *threadUnsafeStringSet) IsProperSubset(other StringSet) bool {
	return set.IsSubset(other) && !set.Equal(other)
}

func (set *threadUnsafeStringSet) IsSuperset(other StringSet) bool {
	return other.IsSubset(set)
}

func (set *threadUnsafeStringSet) IsProperSuperset(other StringSet) bool {
	return set.IsSuperset(other) && !set.Equal(other)
}

func (set *threadUnsafeStringSet) Union(other StringSet) StringSet {
	o := other.(*threadUnsafeStringSet)

	unionedSet := newThreadUnsafeStringSet()

	for elem := range *set {
		unionedSet.Add(elem)
	}
	for elem := range *o {
		unionedSet.Add(elem)
	}
	return &unionedSet
}

func (set *threadUnsafeStringSet) Intersect(other StringSet) StringSet {
	o := other.(*threadUnsafeStringSet)

	intersection := newThreadUnsafeStringSet()
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

func (set *threadUnsafeStringSet) Difference(other StringSet) StringSet {
	_ = other.(*threadUnsafeStringSet)

	difference := newThreadUnsafeStringSet()
	for elem := range *set {
		if !other.Contains(elem) {
			difference.Add(elem)
		}
	}
	return &difference
}

func (set *threadUnsafeStringSet) SymmetricDifference(other StringSet) StringSet {
	_ = other.(*threadUnsafeStringSet)

	aDiff := set.Difference(other)
	bDiff := other.Difference(set)
	return aDiff.Union(bDiff)
}

func (set *threadUnsafeStringSet) Clear() {
	*set = newThreadUnsafeStringSet()
}

func (set *threadUnsafeStringSet) Remove(i string) {
	delete(*set, i)
}

func (set *threadUnsafeStringSet) Cardinality() int {
	return len(*set)
}

func (set *threadUnsafeStringSet) Each(cb func(string) bool) {
	for elem := range *set {
		if cb(elem) {
			break
		}
	}
}

func (set *threadUnsafeStringSet) Iter() <-chan string {
	ch := make(chan string)
	go func() {
		for elem := range *set {
			ch <- elem
		}
		close(ch)
	}()

	return ch
}

func (set *threadUnsafeStringSet) Iterator() *StringIterator {
	iterator, ch, stopCh := newStringIterator()

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

func (set *threadUnsafeStringSet) Equal(other StringSet) bool {
	_ = other.(*threadUnsafeStringSet)

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

func (set *threadUnsafeStringSet) Clone() StringSet {
	clonedSet := newThreadUnsafeStringSet()
	for elem := range *set {
		clonedSet.Add(elem)
	}
	return &clonedSet
}

func (set *threadUnsafeStringSet) String() string {
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

func (set *threadUnsafeStringSet) Pop() string {
	for item := range *set {
		delete(*set, item)
		return item
	}
	return ""
}

/*
// Not yet supported
func (set *threadUnsafeStringSet) PowerSet() StringSet {
	powSet := NewThreadUnsafeStringSet()
	nullset := newThreadUnsafeStringSet()
	powSet.Add(&nullset)

	for es := range *set {
		u := newThreadUnsafeStringSet()
		j := powSet.Iter()
		for er := range j {
			p := newThreadUnsafeStringSet()
			if reflect.TypeOf(er).Name() == "" {
				k := er.(*threadUnsafeStringSet)
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
func (set *threadUnsafeStringSet) CartesianProduct(other StringSet) StringSet {
	o := other.(*threadUnsafeStringSet)
	cartProduct := NewThreadUnsafeStringSet()

	for i := range *set {
		for j := range *o {
			elem := OrderedPair{First: i, Second: j}
			cartProduct.Add(elem)
		}
	}

	return cartProduct
}
*/

func (set *threadUnsafeStringSet) ToSlice() []string {
	keys := make([]string, 0, set.Cardinality())
	for elem := range *set {
		keys = append(keys, elem)
	}

	return keys
}

// MarshalJSON creates a JSON array from the set, it marshals all elements
func (set *threadUnsafeStringSet) MarshalJSON() ([]byte, error) {
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
func (set *threadUnsafeStringSet) UnmarshalJSON(b []byte) error {
	var i []string

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
