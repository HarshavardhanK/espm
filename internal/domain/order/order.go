package order

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusDraft     Status = "DRAFT"
	StatusSubmitted Status = "SUBMITTED"
	StatusCancelled Status = "CANCELLED"
)

// Order represents the order aggregate root
type Order struct {
	ID         uuid.UUID
	CustomerID uuid.UUID

	Status Status
	Items  []OrderItem

	TotalAmount float64

	CreatedAt time.Time
	UpdatedAt time.Time

	Version int
}

type OrderItem struct {
	// OrderItem represents an item in the order
	ProductID uuid.UUID
	Quantity  int
	UnitPrice float64
}

func NewOrder(customerID uuid.UUID) *Order {

	now := time.Now()

	return &Order{

		ID: uuid.New(),

		CustomerID: customerID,
		Status:     StatusDraft,

		Items: make([]OrderItem, 0),

		CreatedAt: now,
		UpdatedAt: now,

		Version: 1,
	}
}

func (o *Order) AddItem(productID uuid.UUID, quantity int, unitPrice float64) error {

	if o.Status != StatusDraft {
		return ErrOrderNotInDraftState
	}

	item := OrderItem{

		ProductID: productID,
		Quantity:  quantity,
		UnitPrice: unitPrice,
	}

	o.Items = append(o.Items, item)
	o.TotalAmount += float64(quantity) * unitPrice

	o.UpdatedAt = time.Now()
	o.Version++

	return nil
}

func (o *Order) RemoveItem(productID uuid.UUID) error {

	if o.Status != StatusDraft {
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

func (o *Order) Submit() error {

	if o.Status != StatusDraft {
		return ErrOrderNotInDraftState
	}

	if len(o.Items) == 0 {
		return ErrOrderHasNoItems
	}

	o.Status = StatusSubmitted

	o.UpdatedAt = time.Now()
	o.Version++

	return nil
}

func (o *Order) Cancel() error {

	if o.Status != StatusDraft && o.Status != StatusSubmitted {
		return ErrOrderCannotBeCancelled
	}

	o.Status = StatusCancelled

	o.UpdatedAt = time.Now()

	o.Version++

	return nil
}
