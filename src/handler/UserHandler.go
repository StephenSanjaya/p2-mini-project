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

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// User godoc
// @Summary User top up
// @Description User top up
// @Tags 	 User
// @Accept   json
// @Produce  json
// @Param car body dto.TopUp true "top up"
// @Success 201 {object} object{message=string,invoice=entity.Invoice}
// @Failure 400 {object} httputil.HTTPError
// @Failure 401 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /users/topup [post]
func (us *UserService) TopUp(c *gin.Context) {

	topup := new(dto.TopUp)

	if err := c.ShouldBindJSON(&topup); err != nil {
		c.Error(httputil.NewError(http.StatusBadRequest, "TopUp: invalid body request", err))
		return
	}

	user_id := int(c.GetFloat64("user_id"))

	user, err := helpers.GetUserByID(us.db, user_id)
	if err != nil {
		c.Error(err)
	}

	updatedDeposit := user.Deposit + topup.Amount

	if res := us.db.Model(&entity.User{}).Where("user_id = ?", user_id).Update("deposit", updatedDeposit); res.Error != nil {
		c.Error(httputil.NewError(http.StatusInternalServerError, "TopUp: failed to top up", res.Error))
	}

	invoiceRes, errInvoice := helpers.CreateInvoiceTopUp(user, topup.Amount)
	if errInvoice != nil {
		c.Error(httputil.NewError(http.StatusInternalServerError, "TopUp: failed to create invoice", errInvoice))
		return
	}

	helpers.SendSuccessTopUp(user.Email, invoiceRes.InvoiceUrl)

	c.JSON(http.StatusOK, gin.H{
		"message": "success top up",
		"invoice": invoiceRes,
	})
}
