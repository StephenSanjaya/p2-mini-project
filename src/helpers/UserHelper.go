package helpers

import (
	"net/http"
	"p2-mini-project/src/entity"
	"p2-mini-project/src/httputil"

	"gorm.io/gorm"
)

func GetCurrentUser(db *gorm.DB, user_id int) (*entity.User, *httputil.HTTPError) {
	user := new(entity.User)

	if res := db.Where("user_id = ?", user_id).First(&user); res.Error != nil {
		return nil, httputil.NewError(http.StatusInternalServerError, "GetCurrentUser: failed to get current user", res.Error)
	}

	return user, nil
}
