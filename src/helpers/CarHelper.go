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

func GetUserDeposit(cs *gorm.DB, user_id int) (float64, *httputil.HTTPError) {
	deposit := 0.0
	if res := cs.Table("users").Select("deposit").Where("user_id = ?", user_id).Scan(&deposit); res.Error != nil {
		return -1, httputil.NewError(http.StatusInternalServerError, "GetUserDeposit: fail to get deposit user", res.Error)
	}
	return deposit, nil
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
	if curr_user_id != return_user_id {
		return false
	}
	return true
}
