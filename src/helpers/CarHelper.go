package helpers

import (
	"errors"
	"fmt"
	"net/http"
	"p2-mini-project/src/dto"
	"p2-mini-project/src/entity"
	"p2-mini-project/src/httputil"
	"time"

	"gorm.io/gorm"
)

func GetCarByID(db *gorm.DB, car_id int) (*entity.Car, *httputil.HTTPError) {
	car := new(entity.Car)

	if res := db.Where("car_id = ?", car_id).First(&car); res.Error != nil {
		return nil, httputil.NewError(http.StatusInternalServerError, "GetCarByID: failed get car by car id", res.Error)
	}

	return car, nil
}

func GetRentalByID(db *gorm.DB, rental_id int) (*dto.Rental, *httputil.HTTPError) {
	rental := new(dto.Rental)

	res := db.Where("rental_id = ?", rental_id).First(&rental)
	if res.Error == gorm.ErrRecordNotFound {
		return nil, httputil.NewError(http.StatusNotFound, "GetRentalByID: rental id not found", res.Error)
	}
	if res.Error != nil {
		return nil, httputil.NewError(http.StatusInternalServerError, "GetRentalByID: failed to get rental", res.Error)
	}

	return rental, nil
}

func CalculateTotalPriceWithFormatStr(r *dto.Rental) float64 {
	returnDate, _ := time.Parse("2006-01-02", r.ReturnDate)
	rentalDate, _ := time.Parse("2006-01-02", r.RentalDate)

	return CalculateTotalPrice(r, rentalDate, returnDate)
}

func CalculateTotalPriceWithFormatRFC(r *dto.Rental) float64 {
	returnDate, _ := time.Parse(time.RFC3339, r.ReturnDate)
	rentalDate, _ := time.Parse(time.RFC3339, r.RentalDate)

	return CalculateTotalPrice(r, rentalDate, returnDate)
}

func CalculateTotalPrice(r *dto.Rental, rentalDate time.Time, returnDate time.Time) float64 {

	diffDay := int(returnDate.Sub(rentalDate).Hours() / 24)

	total_price := (r.Price * float64(diffDay))

	switch r.CouponID {
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

func GetPrice(cs *gorm.DB, r *dto.Rental) (float64, *httputil.HTTPError) {
	price := 0.0
	res := cs.Model(&entity.Car{}).Select("rental_cost_per_day").Where("car_id = ?", r.CarID).First(&price)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return -1, httputil.NewError(http.StatusNotFound, "GetPrice: car id not found", res.Error)
	}
	if res.Error != nil {
		return -1, httputil.NewError(http.StatusInternalServerError, "GetPrice: fail to check stock", res.Error)
	}
	return price, nil
}

func CheckCarStatus(cs *gorm.DB, car_id int) *httputil.HTTPError {
	status := ""
	if res := cs.Model(&entity.Car{}).Where("car_id = ?", car_id).Select("status").First(&status); res.Error != nil {
		return httputil.NewError(http.StatusInternalServerError, "CheckCarStatus: failed to check car status", res.Error)
	}
	if status == "rented" {
		return httputil.NewError(http.StatusInternalServerError, "CheckCarStatus: car is rented", errors.New("can't rental car, status is rented"))
	}

	return nil
}

func CheckAuthorizeUser(curr_user_id, return_user_id int) bool {
	return curr_user_id == return_user_id
}
