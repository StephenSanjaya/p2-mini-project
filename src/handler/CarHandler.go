package handler

import (
	"errors"
	"fmt"
	"net/http"
	"p2-mini-project/src/dto"
	"p2-mini-project/src/entity"
	"p2-mini-project/src/helpers"
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
	cars := new([]entity.Car)

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

	cars := new([]entity.Car)

	res := cs.db.Where("category_id = ?", id).Find(&cars)
	if res.RowsAffected == 0 {
		c.Error(httputil.NewError(http.StatusNotFound, "GetAllCarsByCategory: cateogry id not found", errors.New("cateogry id not found")))
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

	rental := new(dto.Rental)

	if err := c.ShouldBindJSON(&rental); err != nil {
		c.Error(httputil.NewError(http.StatusBadRequest, "RentalCar: invalid body request", err))
		return
	}
	price, err := helpers.GetPrice(cs.db, rental)
	if err != nil {
		c.Error(err)
		return
	}

	rental.UserID = int(c.GetFloat64("user_id"))
	rental.Price = price

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

	rental := new(dto.Rental)
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
	// payment.TotalPrice = CalculateTotalPrice(payment, rental)
	payment.TotalPrice = helpers.CalculateTotalPrice(payment, rental)
	payment.RentalID = rental.ID
	payment.PaymentStatus = "settlement"

	// check rental cost
	// currDeposit, err := GetUserDeposit(cs, int(c.GetFloat64("user_id")))
	currDeposit, err := helpers.GetUserDeposit(cs.db, int(c.GetFloat64("user_id")))
	if err != nil {
		c.Error(err)
		return
	}
	if currDeposit < payment.TotalPrice {
		paymentError := fmt.Sprintf("your deposit is %.2f while total payment is %.2f", currDeposit, payment.TotalPrice)
		c.Error(httputil.NewError(http.StatusBadRequest, "PayRentalCar: your deposit is not enough", errors.New(paymentError)))
		return
	}
	updatedDeposit := currDeposit - payment.TotalPrice

	txErr := cs.db.Transaction(func(tx *gorm.DB) error {

		// update deposit
		if res := tx.Model(&entity.User{}).Where("user_id = ?", int(c.GetFloat64("user_id"))).Update("deposit", updatedDeposit); res.Error != nil {
			return httputil.NewError(http.StatusInternalServerError, "PayRentalCar: failed to update deposit", res.Error)
		}

		// create payment
		if res := tx.Create(&payment); res.Error != nil {
			return httputil.NewError(http.StatusInternalServerError, "PayRentalCar: failed to pay rental car", res.Error)
		}

		// update status
		if res := tx.Model(&entity.Car{}).Where("car_id = ?", rental.CarID).Update("status", "rented"); res.Error != nil {
			return httputil.NewError(http.StatusInternalServerError, "PayRentalCar: failed to update status", res.Error)
		}

		return nil
	})
	if txErr != nil {
		c.Error(txErr)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "success pay rental car",
		"payment": payment,
	})
}

func (cs *CarService) ReturnRentalCar(c *gin.Context) {
	rental_id := c.Param("rental_id")

	rental := new(entity.Rental)

	res := cs.db.Where("rental_id = ?", rental_id).First(&rental)
	if res.Error == gorm.ErrRecordNotFound {
		c.Error(httputil.NewError(http.StatusNotFound, "ReturnRentalCar: id not found", res.Error))
		return
	}
	if res.Error != nil {
		c.Error(httputil.NewError(http.StatusInternalServerError, "ReturnRentalCar: failed to get rental", res.Error))
		return
	}

	// update status
	if res := cs.db.Model(&entity.Car{}).Where("car_id = ?", rental.CarID).Update("status", "available"); res.Error != nil {
		c.Error(httputil.NewError(http.StatusInternalServerError, "ReturnRentalCar: failed to update status", res.Error))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success return rental car",
	})
}
