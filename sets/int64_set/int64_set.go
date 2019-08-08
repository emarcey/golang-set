package mapsetint64

// Int64Set is the primary interface provided by the mapset package.  It
// represents an unordered set of data and a large number of
// operations that can be applied to that set.
type Int64Set interface {
	// Adds an element to the set. Returns whether
	// the item was added.
	Add(i int64) bool

	// Returns the number of elements in the set.
	Cardinality() int

	// Removes all elements from the set, leaving
	// the empty set.
	Clear()

	// Returns a clone of the set using the same
	// implementation, duplicating all keys.
	Clone() Int64Set

	// Returns whether the given items
	// are all in the set.
	Contains(i ...int64) bool

	// Returns the difference between this set
	// and other. The returned set will contain
	// all elements of this set that are not also
	// elements of other.
	//
	// Note that the argument to Difference
	// must be of the same type as the receiver
	// of the method. Otherwise, Difference will
	// panic.
	Difference(other Int64Set) Int64Set

	// Determines if two sets are equal to each
	// other. If they have the same cardinality
	// and contain the same elements, they are
	// considered equal. The order in which
	// the elements were added is irrelevant.
	//
	// Note that the argument to Equal must be
	// of the same type as the receiver of the
	// method. Otherwise, Equal will panic.
	Equal(other Int64Set) bool

	// Returns a new set containing only the elements
	// that exist only in both sets.
	//
	// Note that the argument to Intersect
	// must be of the same type as the receiver
	// of the method. Otherwise, Intersect will
	// panic.
	Intersect(other Int64Set) Int64Set

	// Determines if every element in this set is in
	// the other set but the two sets are not equal.
	//
	// Note that the argument to IsProperSubset
	// must be of the same type as the receiver
	// of the method. Otherwise, IsProperSubset
	// will panic.
	IsProperSubset(other Int64Set) bool

	// Determines if every element in the other set
	// is in this set but the two sets are not
	// equal.
	//
	// Note that the argument to IsSuperset
	// must be of the same type as the receiver
	// of the method. Otherwise, IsSuperset will
	// panic.
	IsProperSuperset(other Int64Set) bool

	// Determines if every element in this set is in
	// the other set.
	//
	// Note that the argument to IsSubset
	// must be of the same type as the receiver
	// of the method. Otherwise, IsSubset will
	// panic.
	IsSubset(other Int64Set) bool

	// Determines if every element in the other set
	// is in this set.
	//
	// Note that the argument to IsSuperset
	// must be of the same type as the receiver
	// of the method. Otherwise, IsSuperset will
	// panic.
	IsSuperset(other Int64Set) bool

	// Iterates over elements and executes the passed func against each element.
	// If passed func returns true, stop iteration at the time.
	Each(func(int64) bool)

	// Returns a channel of elements that you can
	// range over.
	Iter() <-chan int64

	// Returns an Iterator object that you can
	// use to range over the set.
	Iterator() *Int64Iterator

	// Remove a single element from the set.
	Remove(i int64)

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
	SymmetricDifference(other Int64Set) Int64Set

	// Returns a new set with all elements in both sets.
	//
	// Note that the argument to Union must be of the

	// same type as the receiver of the method.
	// Otherwise, IsSuperset will panic.
	Union(other Int64Set) Int64Set

	// Pop removes and returns an arbitrary item from the set.
	Pop() int64

	// Returns all subsets of a given set (Power Int64Set).
	// Not yet supported
	// PowerSet() Int64Set

	// Returns the Cartesian Product of two sets.
	// Not yet supported
	// CartesianProduct(other Int64Set) Int64Set

	// Returns the members of the set as a slice.
	ToSlice() []int64
}

// NewInt64Set creates and returns a reference to an empty set.  Operations
// on the resulting set are thread-safe.
func NewInt64Set(s ...int64) Int64Set {
	set := newThreadSafeInt64Set()
	for _, item := range s {
		set.Add(item)
	}
	return &set
}

// NewInt64SetWith creates and returns a new set with the given elements.
// Operations on the resulting set are thread-safe.
func NewInt64SetWith(elts ...int64) Int64Set {
	return NewInt64SetFromSlice(elts)
}

// NewInt64SetFromSlice creates and returns a reference to a set from an
// existing slice.  Operations on the resulting set are thread-safe.
func NewInt64SetFromSlice(s []int64) Int64Set {
	a := NewInt64Set(s...)
	return a
}

// NewThreadUnsafeInt64Set creates and returns a reference to an empty set.
// Operations on the resulting set are not thread-safe.
func NewThreadUnsafeInt64Set() Int64Set {
	set := newThreadUnsafeInt64Set()
	return &set
}

// NewThreadUnsafeInt64SetFromSlice creates and returns a reference to a
// set from an existing slice.  Operations on the resulting set are
// not thread-safe.
func NewThreadUnsafeInt64SetFromSlice(s []int64) Int64Set {
	a := NewThreadUnsafeInt64Set()
	for _, item := range s {
		a.Add(item)
	}
	return a
}
