package helpers

import (
	"net/http"
	"p2-mini-project/src/entity"
	"p2-mini-project/src/httputil"

	"gorm.io/gorm"
)

func GetUserByID(db *gorm.DB, user_id int) (*entity.User, *httputil.HTTPError) {
	user := new(entity.User)

	if res := db.Where("user_id = ?", user_id).First(&user); res.Error != nil {
		return nil, httputil.NewError(http.StatusInternalServerError, "GetCurrentUser: failed to get current user", res.Error)
	}

	return user, nil
}

func GetUserDeposit(cs *gorm.DB, user_id int) (float64, *httputil.HTTPError) {
	deposit := 0.0
	if res := cs.Table("users").Select("deposit").Where("user_id = ?", user_id).Scan(&deposit); res.Error != nil {
		return -1, httputil.NewError(http.StatusInternalServerError, "GetUserDeposit: fail to get deposit user", res.Error)
	}
	return deposit, nil
}

func GetUserEmail(cs *gorm.DB, user_id int) (string, *httputil.HTTPError) {
	email := ""
	if res := cs.Table("users").Select("email").Where("user_id = ?", user_id).Scan(&email); res.Error != nil {
		return "", httputil.NewError(http.StatusInternalServerError, "GetUserEmail: fail to get email user", res.Error)
	}
	return email, nil
}
