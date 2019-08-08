package mapsettimetime

import (
	"time"
)

// TimeTimeSet is the primary interface provided by the mapset package.  It
// represents an unordered set of data and a large number of
// operations that can be applied to that set.
type TimeTimeSet interface {
	// Adds an element to the set. Returns whether
	// the item was added.
	Add(i time.Time) bool

	// Returns the number of elements in the set.
	Cardinality() int

	// Removes all elements from the set, leaving
	// the empty set.
	Clear()

	// Returns a clone of the set using the same
	// implementation, duplicating all keys.
	Clone() TimeTimeSet

	// Returns whether the given items
	// are all in the set.
	Contains(i ...time.Time) bool

	// Returns the difference between this set
	// and other. The returned set will contain
	// all elements of this set that are not also
	// elements of other.
	//
	// Note that the argument to Difference
	// must be of the same type as the receiver
	// of the method. Otherwise, Difference will
	// panic.
	Difference(other TimeTimeSet) TimeTimeSet

	// Determines if two sets are equal to each
	// other. If they have the same cardinality
	// and contain the same elements, they are
	// considered equal. The order in which
	// the elements were added is irrelevant.
	//
	// Note that the argument to Equal must be
	// of the same type as the receiver of the
	// method. Otherwise, Equal will panic.
	Equal(other TimeTimeSet) bool

	// Returns a new set containing only the elements
	// that exist only in both sets.
	//
	// Note that the argument to Intersect
	// must be of the same type as the receiver
	// of the method. Otherwise, Intersect will
	// panic.
	Intersect(other TimeTimeSet) TimeTimeSet

	// Determines if every element in this set is in
	// the other set but the two sets are not equal.
	//
	// Note that the argument to IsProperSubset
	// must be of the same type as the receiver
	// of the method. Otherwise, IsProperSubset
	// will panic.
	IsProperSubset(other TimeTimeSet) bool

	// Determines if every element in the other set
	// is in this set but the two sets are not
	// equal.
	//
	// Note that the argument to IsSuperset
	// must be of the same type as the receiver
	// of the method. Otherwise, IsSuperset will
	// panic.
	IsProperSuperset(other TimeTimeSet) bool

	// Determines if every element in this set is in
	// the other set.
	//
	// Note that the argument to IsSubset
	// must be of the same type as the receiver
	// of the method. Otherwise, IsSubset will
	// panic.
	IsSubset(other TimeTimeSet) bool

	// Determines if every element in the other set
	// is in this set.
	//
	// Note that the argument to IsSuperset
	// must be of the same type as the receiver
	// of the method. Otherwise, IsSuperset will
	// panic.
	IsSuperset(other TimeTimeSet) bool

	// Iterates over elements and executes the passed func against each element.
	// If passed func returns true, stop iteration at the time.
	Each(func(time.Time) bool)

	// Returns a channel of elements that you can
	// range over.
	Iter() <-chan time.Time

	// Returns an Iterator object that you can
	// use to range over the set.
	Iterator() *TimeTimeIterator

	// Remove a single element from the set.
	Remove(i time.Time)

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
	SymmetricDifference(other TimeTimeSet) TimeTimeSet

	// Returns a new set with all elements in both sets.
	//
	// Note that the argument to Union must be of the

	// same type as the receiver of the method.
	// Otherwise, IsSuperset will panic.
	Union(other TimeTimeSet) TimeTimeSet

	// Pop removes and returns an arbitrary item from the set.
	Pop() time.Time

	// Returns all subsets of a given set (Power TimeTimeSet).
	// Not yet supported
	// PowerSet() TimeTimeSet

	// Returns the Cartesian Product of two sets.
	// Not yet supported
	// CartesianProduct(other TimeTimeSet) TimeTimeSet

	// Returns the members of the set as a slice.
	ToSlice() []time.Time
}

// NewTimeTimeSet creates and returns a reference to an empty set.  Operations
// on the resulting set are thread-safe.
func NewTimeTimeSet(s ...time.Time) TimeTimeSet {
	set := newThreadSafeTimeTimeSet()
	for _, item := range s {
		set.Add(item)
	}
	return &set
}

// NewTimeTimeSetWith creates and returns a new set with the given elements.
// Operations on the resulting set are thread-safe.
func NewTimeTimeSetWith(elts ...time.Time) TimeTimeSet {
	return NewTimeTimeSetFromSlice(elts)
}

// NewTimeTimeSetFromSlice creates and returns a reference to a set from an
// existing slice.  Operations on the resulting set are thread-safe.
func NewTimeTimeSetFromSlice(s []time.Time) TimeTimeSet {
	a := NewTimeTimeSet(s...)
	return a
}

// NewThreadUnsafeTimeTimeSet creates and returns a reference to an empty set.
// Operations on the resulting set are not thread-safe.
func NewThreadUnsafeTimeTimeSet() TimeTimeSet {
	set := newThreadUnsafeTimeTimeSet()
	return &set
}

// NewThreadUnsafeTimeTimeSetFromSlice creates and returns a reference to a
// set from an existing slice.  Operations on the resulting set are
// not thread-safe.
func NewThreadUnsafeTimeTimeSetFromSlice(s []time.Time) TimeTimeSet {
	a := NewThreadUnsafeTimeTimeSet()
	for _, item := range s {
		a.Add(item)
	}
	return a
}
