package mapsetint8

// Int8Set is the primary interface provided by the mapset package.  It
// represents an unordered set of data and a large number of
// operations that can be applied to that set.
type Int8Set interface {
	// Adds an element to the set. Returns whether
	// the item was added.
	Add(i int8) bool

	// Returns the number of elements in the set.
	Cardinality() int

	// Removes all elements from the set, leaving
	// the empty set.
	Clear()

	// Returns a clone of the set using the same
	// implementation, duplicating all keys.
	Clone() Int8Set

	// Returns whether the given items
	// are all in the set.
	Contains(i ...int8) bool

	// Returns the difference between this set
	// and other. The returned set will contain
	// all elements of this set that are not also
	// elements of other.
	//
	// Note that the argument to Difference
	// must be of the same type as the receiver
	// of the method. Otherwise, Difference will
	// panic.
	Difference(other Int8Set) Int8Set

	// Determines if two sets are equal to each
	// other. If they have the same cardinality
	// and contain the same elements, they are
	// considered equal. The order in which
	// the elements were added is irrelevant.
	//
	// Note that the argument to Equal must be
	// of the same type as the receiver of the
	// method. Otherwise, Equal will panic.
	Equal(other Int8Set) bool

	// Returns a new set containing only the elements
	// that exist only in both sets.
	//
	// Note that the argument to Intersect
	// must be of the same type as the receiver
	// of the method. Otherwise, Intersect will
	// panic.
	Intersect(other Int8Set) Int8Set

	// Determines if every element in this set is in
	// the other set but the two sets are not equal.
	//
	// Note that the argument to IsProperSubset
	// must be of the same type as the receiver
	// of the method. Otherwise, IsProperSubset
	// will panic.
	IsProperSubset(other Int8Set) bool

	// Determines if every element in the other set
	// is in this set but the two sets are not
	// equal.
	//
	// Note that the argument to IsSuperset
	// must be of the same type as the receiver
	// of the method. Otherwise, IsSuperset will
	// panic.
	IsProperSuperset(other Int8Set) bool

	// Determines if every element in this set is in
	// the other set.
	//
	// Note that the argument to IsSubset
	// must be of the same type as the receiver
	// of the method. Otherwise, IsSubset will
	// panic.
	IsSubset(other Int8Set) bool

	// Determines if every element in the other set
	// is in this set.
	//
	// Note that the argument to IsSuperset
	// must be of the same type as the receiver
	// of the method. Otherwise, IsSuperset will
	// panic.
	IsSuperset(other Int8Set) bool

	// Iterates over elements and executes the passed func against each element.
	// If passed func returns true, stop iteration at the time.
	Each(func(int8) bool)

	// Returns a channel of elements that you can
	// range over.
	Iter() <-chan int8

	// Returns an Iterator object that you can
	// use to range over the set.
	Iterator() *Int8Iterator

	// Remove a single element from the set.
	Remove(i int8)

	// Provides a convenient string representation
	// of the current state of the set.
	String() string

	// Returns a new set with all elements which are
	// in either this set or the other set but not in both.
	//
	// Note that the argument to SymmetricDifference
	// must be of the same type as the receiver
	// of the method. Otherwise, SymmetricDifference
	// will panic.
	SymmetricDifference(other Int8Set) Int8Set

	// Returns a new set with all elements in both sets.
	//
	// Note that the argument to Union must be of the

	// same type as the receiver of the method.
	// Otherwise, IsSuperset will panic.
	Union(other Int8Set) Int8Set

	// Pop removes and returns an arbitrary item from the set.
	Pop() int8

	// Returns all subsets of a given set (Power Int8Set).
	// Not yet supported
	// PowerSet() Int8Set

	// Returns the Cartesian Product of two sets.
	// Not yet supported
	// CartesianProduct(other Int8Set) Int8Set

	// Returns the members of the set as a slice.
	ToSlice() []int8
}

// NewInt8Set creates and returns a reference to an empty set.  Operations
// on the resulting set are thread-safe.
func NewInt8Set(s ...int8) Int8Set {
	set := newThreadSafeInt8Set()
	for _, item := range s {
		set.Add(item)
	}
	return &set
}

// NewInt8SetWith creates and returns a new set with the given elements.
// Operations on the resulting set are thread-safe.
func NewInt8SetWith(elts ...int8) Int8Set {
	return NewInt8SetFromSlice(elts)
}

// NewInt8SetFromSlice creates and returns a reference to a set from an
// existing slice.  Operations on the resulting set are thread-safe.
func NewInt8SetFromSlice(s []int8) Int8Set {
	a := NewInt8Set(s...)
	return a
}

// NewThreadUnsafeInt8Set creates and returns a reference to an empty set.
// Operations on the resulting set are not thread-safe.
func NewThreadUnsafeInt8Set() Int8Set {
	set := newThreadUnsafeInt8Set()
	return &set
}

// NewThreadUnsafeInt8SetFromSlice creates and returns a reference to a
// set from an existing slice.  Operations on the resulting set are
// not thread-safe.
func NewThreadUnsafeInt8SetFromSlice(s []int8) Int8Set {
	a := NewThreadUnsafeInt8Set()
	for _, item := range s {
		a.Add(item)
	}
	return a
}
