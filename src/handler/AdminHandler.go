package handler

import (
	"errors"
	"fmt"
	"net/http"
	"p2-mini-project/src/dto"
	"p2-mini-project/src/entity"
	"p2-mini-project/src/httputil"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminService struct {
	db *gorm.DB
}

func NewAdminService(db *gorm.DB) *AdminService {
	return &AdminService{db: db}
}

// Admin godoc
// @Summary Create car
// @Description Create new car
// @Tags 	 Admin
// @Accept   json
// @Produce  json
// @Param car body dto.Car true "Create new car"
// @Success 201 {object} object{message=string,car=entity.Car}
// @Failure 400 {object} httputil.HTTPError
// @Failure 401 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /admin/cars [post]
func (as *AdminService) CreateNewCar(c *gin.Context) {
	car := new(entity.Car)

	if err := c.ShouldBindJSON(&car); err != nil {
		c.Error(httputil.NewError(http.StatusBadRequest, "CreateNewCar: invalid body request", err))
		return
	}

	car.Status = "available"
	if res := as.db.Create(&car); res.Error != nil {
		c.Error(httputil.NewError(http.StatusInternalServerError, "CreateNewCar: failed to create new car", res.Error))
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "success create new car",
		"car":     car,
	})
}

// Admin godoc
// @Summary Update car
// @Description Update car by id
// @Tags 	 Admin
// @Accept   json
// @Produce  json
// @Param    car    query     int  true  "car update by car_id"
// @Param car body dto.Car true "Update car"
// @Success 200 {object} object{message=string,car=entity.Car}
// @Failure 400 {object} httputil.HTTPError
// @Failure 401 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /admin/cars/{car_id} [put]
func (as *AdminService) UpdateCar(c *gin.Context) {
	car_id := c.Param("car_id")

	car := new(entity.Car)

	if err := c.ShouldBindJSON(&car); err != nil {
		c.Error(httputil.NewError(http.StatusBadRequest, "CreateNewCar: invalid body request", err))
		return
	}

	car.ID, _ = strconv.Atoi(car_id)

	res := as.db.Model(&car).Updates(entity.Car{CategoryID: car.CategoryID, Name: car.Name, RentalCostPerDay: car.RentalCostPerDay, Capacity: car.Capacity})
	if res.Error != nil {
		msg := fmt.Sprintf("UpdateCar: failed to update car with ID [%d]", car.ID)
		c.Error(httputil.NewError(http.StatusInternalServerError, msg, res.Error))
		return
	}
	if res.RowsAffected == 0 {
		c.Error(httputil.NewError(http.StatusNotFound, "UpdateCar: car id not found", errors.New("car id not found")))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success update car with ID: " + car_id,
		"car":     car,
	})
}

// Admin godoc
// @Summary Delete car
// @Description Delete car by id
// @Tags 	 Admin
// @Accept   json
// @Produce  json
// @Param    car    query     int  true  "car delete by car_id"
// @Success 200 {object} object{message=string}
// @Failure 401 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /admin/cars/{car_id} [delete]
func (as *AdminService) DeleteCar(c *gin.Context) {
	car_id := c.Param("car_id")

	res := as.db.Delete(&entity.Car{}, car_id)
	if res.Error != nil {
		msg := fmt.Sprintf("DeleteCar: failed to delete car with ID [%s]", car_id)
		c.Error(httputil.NewError(http.StatusInternalServerError, msg, res.Error))
		return
	}
	if res.RowsAffected == 0 {
		c.Error(httputil.NewError(http.StatusNotFound, "DeleteCar: car id not found", errors.New("car id not found")))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success delete car with ID: " + car_id,
	})
}

// Admin godoc
// @Summary Get all users
// @Description Get all users
// @Tags 	 Admin
// @Produce  json
// @Success 200 {object} object{message=string,users=[]entity.User}
// @Failure 401 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /admin/users [get]
func (as *AdminService) GetAllUsers(c *gin.Context) {
	users := new([]entity.User)

	res := as.db.Omit("password").Find(&users)
	if res.Error != nil {
		c.Error(httputil.NewError(http.StatusInternalServerError, "GetAllUsers: failed to get all users", res.Error))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success get all users",
		"users":   users,
	})
}

// Admin godoc
// @Summary Get rental history
// @Description Get rental history
// @Tags 	 Admin
// @Produce  json
// @Success 200 {object} object{message=string,rental_history=[]dto.RentalHistory}
// @Failure 401 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /admin/rental-history [get]
func (as *AdminService) GetRentalHistory(c *gin.Context) {
	history := new([]dto.RentalHistory)

	rows, err := as.db.Raw("select r.rental_id, r.rental_date, r.return_date, u.user_id, u.fullname, u.address, c.car_id, c.name, p.total_price from rentals r join users u on r.user_id = u.user_id join cars c on r.car_id = c.car_id join payments p on r.rental_id = p.rental_id").Rows()
	if err != nil {
		c.Error(httputil.NewError(http.StatusInternalServerError, "GetRentalHistory: failed to query", err))
		return
	}
	for rows.Next() {
		var h dto.RentalHistory
		err := rows.Scan(&h.RentalID, &h.RentalDate, &h.ReturnDate, &h.UserID, &h.User.Fullname, &h.User.Address, &h.CarID, &h.Car.Name, &h.TotalPrice)
		if err != nil {
			c.Error(httputil.NewError(http.StatusInternalServerError, "GetRentalHistory: failed to scan query", err))
			return
		}
		*history = append(*history, h)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "success get rental history",
		"rental_history": history,
	})
}
