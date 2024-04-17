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

	rentalAndPay := new(dto.RentalAndPayment)

	if err := c.ShouldBindJSON(&rentalAndPay); err != nil {
		c.Error(httputil.NewError(http.StatusBadRequest, "RentalCar: invalid body request", err))
		return
	}
	price, err := helpers.GetPrice(cs.db, &rentalAndPay.Rental)
	if err != nil {
		c.Error(err)
		return
	}
	err = helpers.CheckCarStatus(cs.db, rentalAndPay.Rental.CarID)
	if err != nil {
		c.Error(err)
		return
	}

	rentalAndPay.Rental.UserID = int(c.GetFloat64("user_id"))
	rentalAndPay.Rental.Price = price

	txErr := cs.db.Transaction(func(tx *gorm.DB) error {

		// create rental
		if res := tx.Create(&rentalAndPay.Rental); res.Error != nil {
			return httputil.NewError(http.StatusInternalServerError, "RentalCar: failed to rental car", res.Error)
		}

		rentalAndPay.Payment.PaymentDate = time.Now().Format("2006-01-02")
		rentalAndPay.Payment.RentalID = rentalAndPay.Rental.ID
		rentalAndPay.Payment.PaymentStatus = "pending"
		rentalAndPay.Payment.TotalPrice = helpers.CalculateTotalPrice(rentalAndPay.Payment.CouponID, &rentalAndPay.Rental)

		// create payment
		if res := tx.Create(&rentalAndPay.Payment); res.Error != nil {
			return httputil.NewError(http.StatusInternalServerError, "RentalCar: failed to create payment", res.Error)
		}

		return nil
	})
	if txErr != nil {
		c.Error(txErr)
		return
	}

	user, err := helpers.GetCurrentUser(cs.db, rentalAndPay.Rental.UserID)
	if err != nil {
		c.Error(err)
		return
	}

	car, err := helpers.GetCarByID(cs.db, rentalAndPay.Rental.CarID)
	if err != nil {
		c.Error(err)
		return
	}

	invoiceRes, errInvoice := helpers.CreateInvoicePayment(rentalAndPay, user, car)
	if errInvoice != nil {
		c.Error(httputil.NewError(http.StatusInternalServerError, "RentalCar: failed to create invoice", errInvoice))
		return
	}

	helpers.SendSuccessPayment(user.Email, invoiceRes.InvoiceUrl)

	c.JSON(http.StatusCreated, gin.H{
		"message":    "success rental a car",
		"rental_car": rentalAndPay,
		"invoice":    invoiceRes,
	})
}

func (cs *CarService) PayRentalCar(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")

	payment_id := c.Param("payment_id")

	payment := new(dto.Payment)
	res := cs.db.Where("payment_id = ?", payment_id).First(&payment)
	if res.Error == gorm.ErrRecordNotFound {
		c.Error(httputil.NewError(http.StatusNotFound, "PayRentalCar: payment id not found", res.Error))
		return
	}
	if res.Error != nil {
		c.Error(httputil.NewError(http.StatusInternalServerError, "PayRentalCar: failed to get rental", res.Error))
		return
	}

	if payment.PaymentStatus != "pending" {
		c.Error(httputil.NewError(http.StatusBadRequest, "PayRentalCar: already paid", errors.New("payment status already settlement")))
		return
	}

	rental := new(dto.Rental)
	res = cs.db.Where("rental_id = ?", payment.RentalID).First(&rental)
	if res.Error == gorm.ErrRecordNotFound {
		c.Error(httputil.NewError(http.StatusNotFound, "PayRentalCar: rental id not found", res.Error))
		return
	}
	if res.Error != nil {
		c.Error(httputil.NewError(http.StatusInternalServerError, "PayRentalCar: failed to get rental", res.Error))
		return
	}

	isLoginUser := helpers.CheckAuthorizeUser(int(c.GetFloat64("user_id")), rental.UserID)
	if !isLoginUser {
		c.Error(httputil.NewError(http.StatusUnauthorized, "ReturnRentalCar: failed to pay rental car", errors.New("only authorize user can do this action")))
		return
	}

	currDeposit, err := helpers.GetUserDeposit(cs.db, int(c.GetFloat64("user_id")))
	if err != nil {
		c.Error(err)
		return
	}
	// check balance and rental cost
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

		// update car status
		if res := tx.Model(&entity.Car{}).Where("car_id = ?", rental.CarID).Update("status", "rented"); res.Error != nil {
			return httputil.NewError(http.StatusInternalServerError, "PayRentalCar: failed to update car status", res.Error)
		}

		// update payment status
		if res := tx.Model(&payment).Where("payment_id = ? ", payment_id).Update("payment_status", "settlement"); res.Error != nil {
			return httputil.NewError(http.StatusInternalServerError, "PayRentalCar: failed to update payment status", res.Error)
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

	isLoginUser := helpers.CheckAuthorizeUser(int(c.GetFloat64("user_id")), rental.UserID)
	if !isLoginUser {
		c.Error(httputil.NewError(http.StatusUnauthorized, "ReturnRentalCar: failed to return rental car", errors.New("only authorize user can do this action")))
		return
	}

	paymentStatus := ""
	if res = cs.db.Model(&entity.Payment{}).Select("payment_status").Where("rental_id = ?", rental.ID).First(&paymentStatus); res.Error != nil {
		c.Error(httputil.NewError(http.StatusInternalServerError, "ReturnRentalCar: failed to get payment status", res.Error))
		return
	}
	car, err := helpers.GetCarByID(cs.db, rental.CarID)
	if err != nil {
		c.Error(err)
		return
	}
	if car.Status == "available" {
		c.Error(httputil.NewError(http.StatusBadRequest, "ReturnRentalCar: car already return", errors.New("car already return")))
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
