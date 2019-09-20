package mapsetstring

// StringSet is the primary interface provided by the mapset package.  It
// represents an unordered set of data and a large number of
// operations that can be applied to that set.
type StringSet interface {
	// Adds an element to the set. Returns whether
	// the item was added.
	Add(i string) bool

	// Returns the number of elements in the set.
	Cardinality() int

	// Removes all elements from the set, leaving
	// the empty set.
	Clear()

	// Returns a clone of the set using the same
	// implementation, duplicating all keys.
	Clone() StringSet

	// Returns whether the given items
	// are all in the set.
	Contains(i ...string) bool

	// Returns the difference between this set
	// and other. The returned set will contain
	// all elements of this set that are not also
	// elements of other.
	//
	// Note that the argument to Difference
	// must be of the same type as the receiver
	// of the method. Otherwise, Difference will
	// panic.
	Difference(other StringSet) StringSet

	// Determines if two sets are equal to each
	// other. If they have the same cardinality
	// and contain the same elements, they are
	// considered equal. The order in which
	// the elements were added is irrelevant.
	//
	// Note that the argument to Equal must be
	// of the same type as the receiver of the
	// method. Otherwise, Equal will panic.
	Equal(other StringSet) bool

	// Returns a new set containing only the elements
	// that exist only in both sets.
	//
	// Note that the argument to Intersect
	// must be of the same type as the receiver
	// of the method. Otherwise, Intersect will
	// panic.
	Intersect(other StringSet) StringSet

	// Determines if every element in this set is in
	// the other set but the two sets are not equal.
	//
	// Note that the argument to IsProperSubset
	// must be of the same type as the receiver
	// of the method. Otherwise, IsProperSubset
	// will panic.
	IsProperSubset(other StringSet) bool

	// Determines if every element in the other set
	// is in this set but the two sets are not
	// equal.
	//
	// Note that the argument to IsSuperset
	// must be of the same type as the receiver
	// of the method. Otherwise, IsSuperset will
	// panic.
	IsProperSuperset(other StringSet) bool

	// Determines if every element in this set is in
	// the other set.
	//
	// Note that the argument to IsSubset
	// must be of the same type as the receiver
	// of the method. Otherwise, IsSubset will
	// panic.
	IsSubset(other StringSet) bool

	// Determines if every element in the other set
	// is in this set.
	//
	// Note that the argument to IsSuperset
	// must be of the same type as the receiver
	// of the method. Otherwise, IsSuperset will
	// panic.
	IsSuperset(other StringSet) bool

	// Iterates over elements and executes the passed func against each element.
	// If passed func returns true, stop iteration at the time.
	Each(func(string) bool)

	// Returns a channel of elements that you can
	// range over.
	Iter() <-chan string

	// Returns an Iterator object that you can
	// use to range over the set.
	Iterator() *StringIterator

	// Remove a single element from the set.
	Remove(i string)

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
	SymmetricDifference(other StringSet) StringSet

	// Returns a new set with all elements in both sets.
	//
	// Note that the argument to Union must be of the

	// same type as the receiver of the method.
	// Otherwise, IsSuperset will panic.
	Union(other StringSet) StringSet

	// Pop removes and returns an arbitrary item from the set.
	Pop() string

	// Returns all subsets of a given set (Power StringSet).
	// Not yet supported
	// PowerSet() StringSet

	// Returns the Cartesian Product of two sets.
	// Not yet supported
	// CartesianProduct(other StringSet) StringSet

	// Returns the members of the set as a slice.
	ToSlice() []string
}

// NewStringSet creates and returns a reference to an empty set.  Operations
// on the resulting set are thread-safe.
func NewStringSet(s ...string) StringSet {
	set := newThreadSafeStringSet()
	for _, item := range s {
		set.Add(item)
	}
	return &set
}

// NewStringSetWith creates and returns a new set with the given elements.
// Operations on the resulting set are thread-safe.
func NewStringSetWith(elts ...string) StringSet {
	return NewStringSetFromSlice(elts)
}

// NewStringSetFromSlice creates and returns a reference to a set from an
// existing slice.  Operations on the resulting set are thread-safe.
func NewStringSetFromSlice(s []string) StringSet {
	a := NewStringSet(s...)
	return a
}

// NewThreadUnsafeStringSet creates and returns a reference to an empty set.
// Operations on the resulting set are not thread-safe.
func NewThreadUnsafeStringSet() StringSet {
	set := newThreadUnsafeStringSet()
	return &set
}

// NewThreadUnsafeStringSetFromSlice creates and returns a reference to a
// set from an existing slice.  Operations on the resulting set are
// not thread-safe.
func NewThreadUnsafeStringSetFromSlice(s []string) StringSet {
	a := NewThreadUnsafeStringSet()
	for _, item := range s {
		a.Add(item)
	}
	return a
}
