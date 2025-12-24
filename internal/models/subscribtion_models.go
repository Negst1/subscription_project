package models

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	ServiceName string     `json:"service_name" db:"service_name" validate:"required"`
	Price       int        `json:"price" db:"price" validate:"required,min=1"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id" validate:"required"`
	StartDate   time.Time  `json:"start_date" db:"start_date" validate:"required"`
	EndDate     *time.Time `json:"end_date,omitempty" db:"end_date"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

type CreateSubscriptionRequest struct {
	ServiceName string  `json:"service_name" validate:"required"`
	Price       int     `json:"price" validate:"required,min=1"`
	UserID      string  `json:"user_id" validate:"required,uuid4"`
	StartDate   string  `json:"start_date" validate:"required,datetime=01-2006"`
	EndDate     *string `json:"end_date,omitempty" validate:"omitempty,datetime=01-2006"`
}

type UpdateSubscriptionRequest struct {
	ServiceName *string `json:"service_name,omitempty"`
	Price       *int    `json:"price,omitempty" validate:"omitempty,min=1"`
	EndDate     *string `json:"end_date,omitempty" validate:"omitempty,datetime=01-2006"`
}

type SubscriptionSummary struct {
	TotalCost int `json:"total_cost"`
}

type SummaryRequest struct {
	StartDate   string  `json:"start_date" validate:"required,datetime=01-2006"`
	EndDate     string  `json:"end_date" validate:"required,datetime=01-2006"`
	UserID      *string `json:"user_id,omitempty" validate:"omitempty,uuid4"`
	ServiceName *string `json:"service_name,omitempty"`
}
