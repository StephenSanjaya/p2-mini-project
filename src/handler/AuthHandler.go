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
	user := new(dto.User)

	if err := c.Bind(&user); err != nil {
		c.Error(httputil.NewError(http.StatusBadRequest, "RegisterHandler: invalid body request", err))
		return
	}

	hashedPassword := HashPassword(user.Password)
	user.Password = hashedPassword

	if res := as.db.Create(&user); res.Error != nil {
		c.Error(httputil.NewError(http.StatusInternalServerError, "RegisterHandler: register failed", res.Error))
		return
	}

	user.Password = ""

	c.JSON(http.StatusCreated, gin.H{
		"message": "registration successful",
		"user":    user,
	})
}

func (as *AuthService) LoginHandler(c *gin.Context) {
	login := new(dto.Login)

	if err := c.ShouldBindJSON(&login); err != nil {
		c.Error(httputil.NewError(http.StatusBadRequest, "LoginHandler: invalid body request", err))
		return
	}

	user := new(dto.User)
	res := as.db.Where("email = ?", login.Email).First(&user)
	if res.Error == gorm.ErrRecordNotFound {
		c.Error(httputil.NewError(http.StatusNotFound, "LoginHandler: email not found", res.Error))
		return
	}
	if res.Error != nil {
		c.Error(httputil.NewError(http.StatusInternalServerError, "LoginHandler: login failed", res.Error))
		return
	}

	if err := CheckHashPassword(user.Password, login.Password); err != nil {
		c.Error(httputil.NewError(http.StatusBadRequest, "LoginHandler: password not match", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
	})
}

func HashPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword)
}

func CheckHashPassword(hashedPass string, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(password)); err != nil {
		return err
	}
	return nil
}
