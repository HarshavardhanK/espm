package order

import "errors"

var (
	// ErrOrderNotInDraftState is returned when an operation requires the order to be in draft state
	ErrOrderNotInDraftState = errors.New("order is not in draft state")

	// ErrItemNotFound is returned when trying to remove a non-existent item
	ErrItemNotFound = errors.New("item not found in order")

	// ErrOrderHasNoItems is returned when trying to submit an order without items
	ErrOrderHasNoItems = errors.New("order has no items")

	// ErrOrderCannotBeCancelled is returned when trying to cancel an order in an invalid state
	ErrOrderCannotBeCancelled = errors.New("order cannot be cancelled in current state")
)
