package handler

import (
	"errors"
	"fmt"
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
	price, err := GetPrice(cs, c, rental)
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
	payment.TotalPrice = CalculateTotalPrice(payment, rental)
	payment.RentalID = rental.ID

	txErr := cs.db.Transaction(func(tx *gorm.DB) error {
		// check rental cost & update deposit
		currDeposit, err := GetUserDeposit(tx, int(c.GetFloat64("user_id")))
		if err != nil {
			return err
		}
		if currDeposit < payment.TotalPrice {
			paymentError := fmt.Sprintf("your deposit is %.2f while total payment is %.2f", currDeposit, payment.TotalPrice)
			return httputil.NewError(http.StatusBadRequest, "PayRentalCar: your deposit is not enough", errors.New(paymentError))
		}
		updatedDeposit := currDeposit - payment.TotalPrice

		if res := tx.Model(&entity.User{}).Where("user_id = ?", int(c.GetFloat64("user_id"))).Update("deposit", updatedDeposit); res.Error != nil {
			return httputil.NewError(http.StatusInternalServerError, "PayRentalCar: failed to update deposit", res.Error)
		}

		if res := tx.Create(&payment); res.Error != nil {
			return httputil.NewError(http.StatusInternalServerError, "PayRentalCar: failed to pay rental car", res.Error)
		}

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

func GetUserDeposit(tx *gorm.DB, user_id int) (float64, *httputil.HTTPError) {
	deposit := 0.0
	if res := tx.Table("users").Select("deposit").Where("user_id = ?", user_id).Scan(&deposit); res.Error != nil {
		return -1, httputil.NewError(http.StatusInternalServerError, "GetUserDeposit: fail to get deposit user", res.Error)
	}
	return deposit, nil
}

func GetPrice(cs *CarService, c *gin.Context, r *dto.Rental) (float64, *httputil.HTTPError) {
	price := 0.0
	res := cs.db.Model(&entity.Car{}).Select("rental_cost_per_day").Where("car_id = ?", r.CarID).First(&price)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return -1, httputil.NewError(http.StatusNotFound, "GetPrice: car id not found", res.Error)
	}
	if res.Error != nil {
		return -1, httputil.NewError(http.StatusInternalServerError, "GetPrice: fail to check stock", res.Error)
	}
	return price, nil
}

func CalculateTotalPrice(p *dto.Payment, r *dto.Rental) float64 {

	returnDate, _ := time.Parse(time.RFC3339, r.ReturnDate)
	rentalDate, _ := time.Parse(time.RFC3339, r.RentalDate)

	diffDay := int(returnDate.Sub(rentalDate).Hours() / 24)

	total_price := (r.Price * float64(diffDay))

	switch p.CouponID {
	case 1:
		total_price = total_price - (total_price * 0.1)
	case 2:
		total_price = total_price - (total_price * 0.2)
	case 3:
		total_price = total_price - (total_price * 0.3)
	default:
		fmt.Println("can't detect coupon")
	}

	return total_price
}
