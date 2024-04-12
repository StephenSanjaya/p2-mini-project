package handler

import (
	"net/http"
	"p2-mini-project/src/dto"
	"p2-mini-project/src/httputil"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

func (as *AuthService) RegisterHandler(c *gin.Context) {
	user := new(dto.Register)

	if err := c.Bind(&user); err != nil {
		c.Error(httputil.NewError(http.StatusBadRequest, "RegisterHandler: invalid body request", err))
	}

	hashedPassword := HashPassword(user.Password)
	user.Password = hashedPassword

	if res := as.db.Create(&user); res.Error != nil {
		c.Error(httputil.NewError(http.StatusInternalServerError, "RegisterHandler: register failed", res.Error))
	}

	user.Password = ""

	httputil.NewSuccess(c, http.StatusCreated, "registration successful", user)

}

func HashPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword)
}
