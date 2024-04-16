package handler

import (
	"net/http"
	"p2-mini-project/src/entity"
	"p2-mini-project/src/httputil"

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

	if res := as.db.Create(&car); res.Error != nil {
		c.Error(httputil.NewError(http.StatusInternalServerError, "CreateNewCar: failed to create new car", res.Error))
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "success create new car",
		"car":     car,
	})
}
