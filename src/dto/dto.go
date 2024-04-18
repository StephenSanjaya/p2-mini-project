package dto

type User struct {
	Fullname string  `json:"fullname" binding:"required"`
	Address  string  `json:"address" binding:"required"`
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password,omitempty" binding:"required" swaggerignore:"true"`
	Role     string  `json:"role" swaggerignore:"true"`
	Deposit  float64 `json:"deposit"`
}

type Login struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type Rental struct {
	ID         int     `json:"rental_id" gorm:"column:rental_id" swaggerignore:"true"`
	UserID     int     `json:"user_id" swaggerignore:"true"`
	CarID      int     `json:"car_id" binding:"required"`
	CouponID   int     `json:"coupon_id"`
	Price      float64 `json:"price" swaggerignore:"true"`
	RentalDate string  `json:"rental_date" binding:"required"`
	ReturnDate string  `json:"return_date" binding:"required"`
}

type Car struct {
	CategoryID       int     `json:"category_id"`
	Name             string  `json:"name"`
	RentalCostPerDay float64 `json:"rental_cost_per_day"`
	Capacity         float64 `json:"capacity"`
}

type Payment struct {
	PaymentMethodID int     `json:"payment_method_id" binding:"required"`
	RentalID        int     `json:"rental_id" swaggerignore:"true"`
	TotalPrice      float64 `json:"total_price" swaggerignore:"true"`
	PaymentDate     string  `json:"payment_date" swaggerignore:"true"`
	PaymentStatus   string  `json:"payment_status" swaggerignore:"true"`
}

type RentalAndPayment struct {
	Rental  Rental  `json:"rental"`
	Payment Payment `json:"payment"`
}

type RentalHistory struct {
	RentalID   int               `json:"rental_id"`
	RentalDate string            `json:"rental_date"`
	ReturnDate string            `json:"return_date"`
	UserID     int               `json:"user_id"`
	User       UserRentalHistory `json:"user"`
	CarID      int               `json:"car_id"`
	Car        CarRentalHistory  `json:"car"`
	TotalPrice float64           `json:"total_price"`
}

type UserRentalHistory struct {
	Fullname string `json:"fullname"`
	Address  string `json:"address"`
}

type CarRentalHistory struct {
	Name string `json:"name"`
}

type TopUp struct {
	Amount float64 `json:"amount"`
}
