package handler

import (
	"net/http"
	"p2-mini-project/src/dto"
	"p2-mini-project/src/entity"
	"p2-mini-project/src/httputil"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CarService struct {
	db *gorm.DB
}

func NewCarService(db *gorm.DB) *CarService {
	return &CarService{db: db}
}

func (cs *CarService) GetAllCars(c *gin.Context) {
	cars := new(entity.Car)

	if res := cs.db.Find(&cars); res.Error != nil {
		c.Error(httputil.NewError(http.StatusInternalServerError, "GetAllCars: fail to get all cars", res.Error))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success get all cars",
		"cars":    cars,
	})
}

func (cs *CarService) GetAllCarsByCategory(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")

	id := c.Param("category_id")

	cars := new(entity.Car)

	res := cs.db.Where("category_id = ?", id).Find(&cars)
	if res.Error == gorm.ErrRecordNotFound {
		c.Error(httputil.NewError(http.StatusNotFound, "GetAllCarsByCategory: cateogry id not found", res.Error))
		return
	}
	if res.Error != nil {
		c.Error(httputil.NewError(http.StatusInternalServerError, "GetAllCarsByCategory: fail to get all cars by category", res.Error))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success get all cars by category",
		"cars":    cars,
	})
}

func (cs *CarService) RentalCar(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")

	rental := new(entity.Rental)

	if err := c.ShouldBindJSON(&rental); err != nil {
		c.Error(httputil.NewError(http.StatusBadRequest, "RentalCar: invalid body request", err))
		return
	}
	rental.UserID = c.GetInt("user_id")

	if res := cs.db.Create(&rental); res.Error != nil {
		c.Error(httputil.NewError(http.StatusInternalServerError, "RentalCar: failed to rental car", res.Error))
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "success rental a car",
		"rental_car": rental,
	})
}

func (cs *CarService) PayRentalCar(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")

	payment := new(dto.Payment)
	rental_id := c.Param("rental_id")

	if err := c.ShouldBindJSON(&payment); err != nil {
		c.Error(httputil.NewError(http.StatusBadRequest, "PayRentalCar: invalid body request", err))
		return
	}

	rental := new(entity.Rental)
	res := cs.db.Where("rental_id = ?", rental_id).First(&rental)
	if res.Error == gorm.ErrRecordNotFound {
		c.Error(httputil.NewError(http.StatusNotFound, "PayRentalCar: rental id not found", res.Error))
		return
	}
	if res.Error != nil {
		c.Error(httputil.NewError(http.StatusInternalServerError, "PayRentalCar: failed to get rental", res.Error))
		return
	}
	payment.PaymentDate = time.Now().Format("2006-01-02")
	payment.TotalPrice = CalculateTotalPrice(payment, rental)

	if res := cs.db.Create(&payment); res.Error != nil {
		c.Error(httputil.NewError(http.StatusInternalServerError, "PayRentalCar: failed to rental car", res.Error))
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "success pay rental a car",
		"rental_car": rental,
	})
}

func CalculateTotalPrice(p *dto.Payment, r *entity.Rental) float64 {
	total_price := r.Price * float64(r.Quantity)

	switch p.CouponID {
	case 1:
		total_price = total_price - (total_price * 0.1)
	case 2:
		total_price = total_price - (total_price * 0.2)
	case 3:
		total_price = total_price - (total_price * 0.3)
	}

	return total_price
}
