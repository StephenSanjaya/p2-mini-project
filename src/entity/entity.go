package entity

import (
	"gorm.io/datatypes"
)

type User struct {
	ID       int      `json:"user_id" gorm:"primaryKey;column:user_id"`
	Fullname string   `json:"fullname" gorm:"type:string;size:255;not null;"`
	Address  string   `json:"address" gorm:"type:string;size:255;not null;"`
	Email    string   `json:"email" gorm:"type:string;size:255;not null;unique;"`
	Password string   `json:"password,omitempty" gorm:"type:string;size:255;not null;" swaggerignore:"true"`
	Role     string   `json:"role" gorm:"type:string;size:255;not null;"`
	Deposit  float64  `json:"deposit,omitempty" gorm:"not null;default:0.0"`
	Rentals  []Rental `json:"rentals,omitempty" swaggerignore:"true"`
}

type Car struct {
	ID               int      `json:"car_id" gorm:"primaryKey;column:car_id"`
	CategoryID       int      `json:"category_id" gorm:"not null"`
	Name             string   `json:"name" gorm:"type:string;size:255;not null;"`
	Status           string   `json:"status,omitempty" gorm:"not null"`
	RentalCostPerDay float64  `json:"rental_cost_per_day" gorm:"not null"`
	Capacity         float64  `json:"capacity" gorm:"not null"`
	Rentals          []Rental `json:"rentals,omitempty" swaggerignore:"true"`
}

type Category struct {
	ID   int    `json:"category_id" gorm:"primaryKey;column:category_id"`
	Type string `json:"type" gorm:"type:string;size:255;not null;"`
	Cars []Car  `json:"cars,omitempty"`
}

type Rental struct {
	ID         int            `json:"rental_id" gorm:"primaryKey;column:rental_id" swaggerignore:"true"`
	UserID     int            `json:"user_id" gorm:"not null"`
	CarID      int            `json:"car_id" gorm:"not null"`
	CouponID   int            `json:"coupon_id" gorm:"not null"`
	Price      float64        `json:"price" gorm:"not null"`
	RentalDate datatypes.Date `json:"rental_date" gorm:"not null"`
	ReturnDate datatypes.Date `json:"return_date" gorm:"not null"`
}

type Payment struct {
	ID              int            `json:"payment_id" gorm:"primaryKey;column:payment_id" swaggerignore:"true"`
	RentalID        int            `json:"rental_id" gorm:"unique;not null" `
	Rental          Rental         `json:"rental" swaggerignore:"true"`
	PaymentMethodID int            `json:"payment_method_id" gorm:"not null"`
	TotalPrice      float64        `json:"total_price" gorm:"not null"`
	PaymentStatus   string         `json:"payment_status" gorm:"not null;default:settlement"`
	PaymentDate     datatypes.Date `json:"payment_date" gorm:"not null"`
}

type Coupon struct {
	ID         int      `json:"coupon_id" gorm:"primaryKey;column:coupon_id"`
	CouponName string   `json:"coupon_name" gorm:"type:string;size:255;not null;"`
	Payments   []Rental `json:"payments,omitempty"`
}

type PaymentMethod struct {
	ID          int       `json:"payment_method_id" gorm:"primaryKey;column:payment_method_id"`
	PaymentName string    `json:"payment_name" gorm:"type:string;size:255;not null;"`
	Payments    []Payment `json:"payments,omitempty"`
}

type Invoice struct {
	ID         string `json:"id"`
	InvoiceUrl string `json:"invoice_url"`
}
