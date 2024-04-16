package handler

import (
	"errors"
	"fmt"
	"net/http"
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
