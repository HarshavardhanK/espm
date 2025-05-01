package domain

import (
	"time"

	"github.com/google/uuid"
)

// OrderStatus represents the current state of an order
type OrderStatus string

const (
	OrderStatusDraft     OrderStatus = "DRAFT"
	OrderStatusSubmitted OrderStatus = "SUBMITTED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

// Order represents the main aggregate root
type Order struct {
	ID          uuid.UUID
	CustomerID  uuid.UUID
	Status      OrderStatus
	Items       []OrderItem
	TotalAmount float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Version     int
}

// OrderItem represents an item in an order
type OrderItem struct {
	ProductID uuid.UUID
	Quantity  int
	UnitPrice float64
}

// NewOrder creates a new order in DRAFT state
func NewOrder(customerID uuid.UUID) *Order {
	now := time.Now()
	return &Order{
		ID:         uuid.New(),
		CustomerID: customerID,
		Status:     OrderStatusDraft,
		Items:      make([]OrderItem, 0),
		CreatedAt:  now,
		UpdatedAt:  now,
		Version:    1,
	}
}

// AddItem adds an item to the order
func (o *Order) AddItem(productID uuid.UUID, quantity int, unitPrice float64) error {
	if o.Status != OrderStatusDraft {
		return ErrOrderNotInDraftState
	}

	o.Items = append(o.Items, OrderItem{
		ProductID: productID,
		Quantity:  quantity,
		UnitPrice: unitPrice,
	})

	o.TotalAmount += float64(quantity) * unitPrice
	o.UpdatedAt = time.Now()
	o.Version++

	return nil
}

// RemoveItem removes an item from the order
func (o *Order) RemoveItem(productID uuid.UUID) error {
	if o.Status != OrderStatusDraft {
		return ErrOrderNotInDraftState
	}

	for i, item := range o.Items {
		if item.ProductID == productID {
			o.TotalAmount -= float64(item.Quantity) * item.UnitPrice
			o.Items = append(o.Items[:i], o.Items[i+1:]...)
			o.UpdatedAt = time.Now()
			o.Version++
			return nil
		}
	}

	return ErrItemNotFound
}

// Submit changes the order status to SUBMITTED
func (o *Order) Submit() error {
	if o.Status != OrderStatusDraft {
		return ErrOrderNotInDraftState
	}

	if len(o.Items) == 0 {
		return ErrOrderHasNoItems
	}

	o.Status = OrderStatusSubmitted
	o.UpdatedAt = time.Now()
	o.Version++

	return nil
}

// Cancel changes the order status to CANCELLED
func (o *Order) Cancel() error {
	if o.Status != OrderStatusSubmitted {
		return ErrOrderCannotBeCancelled
	}

	o.Status = OrderStatusCancelled
	o.UpdatedAt = time.Now()
	o.Version++

	return nil
}

// Domain errors
var (
	ErrOrderNotInDraftState   = NewDomainError("order is not in draft state")
	ErrItemNotFound           = NewDomainError("item not found in order")
	ErrOrderHasNoItems        = NewDomainError("order has no items")
	ErrOrderCannotBeCancelled = NewDomainError("order cannot be cancelled")
)

// DomainError represents a domain-specific error
type DomainError struct {
	message string
}

func NewDomainError(message string) error {
	return &DomainError{message: message}
}

func (e *DomainError) Error() string {
	return e.message
}
