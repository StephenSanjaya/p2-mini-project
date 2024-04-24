package handler

import (
	"errors"
	"fmt"
	"net/http"
	"p2-mini-project/src/dto"
	"p2-mini-project/src/entity"
	"p2-mini-project/src/helpers"
	"p2-mini-project/src/httputil"
	"strconv"
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

// Car godoc
// @Summary Get all cars
// @Description Get all cars
// @Tags 	 Car
// @Produce  json
// @Success 200 {object} object{message=string,cars=[]entity.Car}
// @Failure 401 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /cars [get]
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

// Car godoc
// @Summary Get cars by category
// @Description Get cars by category id
// @Tags 	 Car
// @Accept   json
// @Produce  json
// @Param    category    query     int  true  "cars search by category_id"
// @Success 200 {object} object{message=string,cars=[]entity.Car}
// @Failure 401 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /cars/{category_id} [get]
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

// Car godoc
// @Summary Rent a car
// @Description Rent a car
// @Tags 	 Car
// @Accept   json
// @Produce  json
// @Param rental body dto.Rental true "user rent a car"
// @Success 201 {object} object{message=string,rental=entity.Rental,invoice=entity.Invoice}
// @Failure 400 {object} httputil.HTTPError
// @Failure 401 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /cars [post]
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
	err = helpers.CheckCarStatus(cs.db, rental.CarID)
	if err != nil {
		c.Error(err)
		return
	}

	rental.UserID = int(c.GetFloat64("user_id"))
	rental.Price = price

	// create rental
	if res := cs.db.Create(&rental); res.Error != nil {
		c.Error(httputil.NewError(http.StatusInternalServerError, "RentalCar: failed to rental car", res.Error))
		return
	}

	user, err := helpers.GetUserByID(cs.db, rental.UserID)
	if err != nil {
		c.Error(err)
		return
	}

	car, err := helpers.GetCarByID(cs.db, rental.CarID)
	if err != nil {
		c.Error(err)
		return
	}

	totalPrice := helpers.CalculateTotalPriceWithFormatStr(rental)
	invoiceRes, errInvoice := helpers.CreateInvoiceRental(&totalPrice, user, car)
	if errInvoice != nil {
		c.Error(httputil.NewError(http.StatusInternalServerError, "RentalCar: failed to create invoice", errInvoice))
		return
	}

	helpers.SendSuccessRental(user.Email, invoiceRes.InvoiceUrl)

	c.JSON(http.StatusCreated, gin.H{
		"message": "success rental a car",
		"rental":  rental,
		"invoice": invoiceRes,
	})
}

// Car godoc
// @Summary Pay rented car
// @Description Pay rented car
// @Tags 	 Car
// @Accept   json
// @Produce  json
// @Param    pay    query     int  true  "pay rental car by rental_id"
// @Param pay body dto.Payment true "user pay rented a car"
// @Success 200 {object} object{message=string,payment=entity.Payment}
// @Failure 401 {object} httputil.HTTPError
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /cars/pay/{rental_id} [post]
func (cs *CarService) PayRentalCar(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")

	rental_id, _ := strconv.Atoi(c.Param("rental_id"))

	payment := new(dto.Payment)

	if err := c.ShouldBindJSON(&payment); err != nil {
		c.Error(httputil.NewError(http.StatusBadRequest, "PayRentalCar: invalid body request", err))
		return
	}

	rental, err := helpers.GetRentalByID(cs.db, rental_id)
	if err != nil {
		c.Error(err)
		return
	}

	payment.PaymentStatus = "settlement"
	payment.PaymentDate = time.Now().Format("2006-01-02")
	payment.TotalPrice = helpers.CalculateTotalPriceWithFormatRFC(rental)
	payment.RentalID = rental.ID

	isLoginUser := helpers.CheckAuthorizeUser(int(c.GetFloat64("user_id")), rental.UserID)
	if !isLoginUser {
		c.Error(httputil.NewError(http.StatusUnauthorized, "PayRentalCar: failed to pay rental car", errors.New("only authorize user can do this action")))
		return
	}

	currDeposit, err := helpers.GetUserDeposit(cs.db, rental.UserID)
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
		if res := tx.Model(&entity.User{}).Where("user_id = ?", rental.UserID).Update("deposit", updatedDeposit); res.Error != nil {
			return httputil.NewError(http.StatusInternalServerError, "PayRentalCar: failed to update deposit", res.Error)
		}

		// update car status
		if res := tx.Model(&entity.Car{}).Where("car_id = ?", rental.CarID).Update("status", "rented"); res.Error != nil {
			return httputil.NewError(http.StatusInternalServerError, "PayRentalCar: failed to update car status", res.Error)
		}

		// create payment
		res := tx.Create(&payment)
		if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
			return httputil.NewError(http.StatusBadRequest, "PayRentalCar: already paid", res.Error)
		}
		if res.Error != nil {
			return httputil.NewError(http.StatusInternalServerError, "PayRentalCar: failed to create payment", res.Error)
		}

		return nil
	})
	if txErr != nil {
		c.Error(txErr)
		return
	}

	email, err := helpers.GetUserEmail(cs.db, rental.UserID)
	if err != nil {
		c.Error(err)
		return
	}

	helpers.SendSuccessPayment(email, payment.TotalPrice)

	c.JSON(http.StatusCreated, gin.H{
		"message": "success pay rental car",
		"payment": payment,
	})
}

// Car godoc
// @Summary Return rented car
// @Description Return rented car
// @Tags 	 Car
// @Accept   json
// @Produce  json
// @Param    rental    query     int  true  "return rental car by rental_id"
// @Success 200 {object} object{message=string}
// @Failure 401 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /cars/return/{rental_id} [post]
func (cs *CarService) ReturnRentalCar(c *gin.Context) {
	rental_id, _ := strconv.Atoi(c.Param("rental_id"))

	rental, err := helpers.GetRentalByID(cs.db, rental_id)
	if err != nil {
		c.Error(err)
		return
	}

	isLoginUser := helpers.CheckAuthorizeUser(int(c.GetFloat64("user_id")), rental.UserID)
	if !isLoginUser {
		c.Error(httputil.NewError(http.StatusUnauthorized, "ReturnRentalCar: failed to return rental car", errors.New("only authorize user can do this action")))
		return
	}

	car, err := helpers.GetCarByID(cs.db, rental.CarID)
	if err != nil {
		c.Error(err)
		return
	}
	if car.Status == "available" {
		c.Error(httputil.NewError(http.StatusBadRequest, "ReturnRentalCar: car already return", errors.New("car status is available")))
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
