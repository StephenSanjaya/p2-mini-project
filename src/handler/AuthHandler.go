package handler

import (
	"net/http"
	"p2-mini-project/src/dto"
	"p2-mini-project/src/entity"
	"p2-mini-project/src/helpers"
	"p2-mini-project/src/httputil"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

// Auth godoc
// @Summary Create Users
// @Description Create new users
// @Tags 	 User
// @Accept   json
// @Produce  json
// @Param user body dto.User true "Create new user"
// @Success 201 {object} object{message=string,user=entity.User}
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /users/register [post]
func (as *AuthService) RegisterHandler(c *gin.Context) {
	user := new(dto.User)

	if err := c.ShouldBindJSON(&user); err != nil {
		c.Error(httputil.NewError(http.StatusBadRequest, "RegisterHandler: invalid body request", err))
		return
	}

	hashedPassword := helpers.HashPassword(user.Password)
	user.Password = hashedPassword
	user.Role = "user"

	if res := as.db.Create(&user); res.Error != nil {
		c.Error(httputil.NewError(http.StatusInternalServerError, "RegisterHandler: register failed", res.Error))
		return
	}

	helpers.SendSuccessRegister(user.Email)

	user.Password = ""

	c.JSON(http.StatusCreated, gin.H{
		"message": "registration successful",
		"user":    user,
	})
}

// Auth godoc
// @Summary User login
// @Description User do login
// @Tags 	 User
// @Accept   json
// @Produce  json
// @Param user body dto.Login true "login user"
// @Success 200 {object} object{message=string,token=string}
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /users/login [post]
func (as *AuthService) LoginHandler(c *gin.Context) {
	login := new(dto.Login)

	if err := c.ShouldBindJSON(&login); err != nil {
		c.Error(httputil.NewError(http.StatusBadRequest, "LoginHandler: invalid body request", err))
		return
	}

	user := new(entity.User)
	res := as.db.Where("email = ?", login.Email).First(&user)
	if res.Error == gorm.ErrRecordNotFound {
		c.Error(httputil.NewError(http.StatusNotFound, "LoginHandler: email not found", res.Error))
		return
	}
	if res.Error != nil {
		c.Error(httputil.NewError(http.StatusInternalServerError, "LoginHandler: login failed", res.Error))
		return
	}

	if err := helpers.CheckHashPassword(user.Password, login.Password); err != nil {
		c.Error(httputil.NewError(http.StatusBadRequest, "LoginHandler: password not match", err))
		return
	}

	tokenString, err := helpers.CreateJWT(user)
	if err != nil {
		c.Error(httputil.NewError(http.StatusInternalServerError, "failed create token", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"token":   tokenString,
	})
}
