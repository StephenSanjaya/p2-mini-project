package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"p2-mini-project/src/dto"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func DbMock(t *testing.T) (*sql.DB, *gorm.DB, sqlmock.Sqlmock) {
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	gormdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqldb,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		t.Fatal(err)
	}
	return sqldb, gormdb, mock
}

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func MockJsonPost(c *gin.Context, content interface{}) {
	c.Request.Method = "POST"
	c.Request.Header.Set("Content-Type", "application/json")

	jsonbytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
}

func TestCreateNewCar_shouldSuccess(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	adminService := NewAdminService(db)

	addRow := sqlmock.NewRows([]string{"category_id", "name", "rental_cost_per_day", "capacity"}).AddRow(1, "test", 50000, 123)
	expectedSQL := "INSERT INTO \"cars\" (.+) VALUES (.+)"
	mock.ExpectBegin()
	mock.ExpectQuery(expectedSQL).WillReturnRows(addRow)
	mock.ExpectCommit()

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		Header: make(http.Header),
	}

	MockJsonPost(ctx, dto.Car{CategoryID: 1, Name: "toyota vios", RentalCostPerDay: 30000, Capacity: 4})
	adminService.CreateNewCar(ctx)

	assert.Equal(t, http.StatusCreated, ctx.Writer.Status())
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestUpdateCar_shouldSuccess(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	adminService := NewAdminService(db)

	updUserSQL := "UPDATE \"cars\" SET .+"
	mock.ExpectBegin()
	mock.ExpectExec(updUserSQL).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		Header: make(http.Header),
	}

	MockJsonPost(ctx, dto.Car{CategoryID: 1, Name: "toyota vios", RentalCostPerDay: 30000, Capacity: 4})
	ctx.AddParam("car_id", "1")
	adminService.UpdateCar(ctx)

	assert.Equal(t, http.StatusOK, ctx.Writer.Status())
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestDeleteCar_shouldSuccess(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	adminService := NewAdminService(db)

	delSQL := "DELETE FROM \"cars\" WHERE \"cars\".\"car_id\" = .+"
	mock.ExpectBegin()
	mock.ExpectExec(delSQL).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		Header: make(http.Header),
	}

	ctx.AddParam("car_id", "1")
	adminService.DeleteCar(ctx)

	assert.Equal(t, http.StatusOK, ctx.Writer.Status())
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetAllUsers_shouldSuccess(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	adminService := NewAdminService(db)

	users := sqlmock.NewRows([]string{"user_id", "full_name", "address", "email", "password", "role", "deposit"}).
		AddRow(1, "user", "jl. user", "user@email.com", "user123", "user", 0.0).
		AddRow(2, "user2", "jl. user2", "user2@email.com", "user123", "user", 0.0)

	expectedSQL := "SELECT (.+) FROM \"users\""
	mock.ExpectQuery(expectedSQL).WillReturnRows(users)

	log.Default().Println("test 123")
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		Header: make(http.Header),
	}

	adminService.GetAllUsers(ctx)

	assert.Equal(t, http.StatusOK, ctx.Writer.Status())
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetRentalHistory_shouldSuccess(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	adminService := NewAdminService(db)

	all_history := sqlmock.NewRows([]string{"rental_id", "rental_date", "return_date", "user_id", "fullname", "address", "car_id", "name", "total_price"}).
		AddRow(1, "2024-04-18", "2024-04-20", 1, "user", "jl user123", 1, "toyota", 30000).
		AddRow(2, "2024-04-18", "2024-04-20", 2, "user", "jl user124", 2, "camry", 50000)

	expectedSQL := "select r.rental_id, r.rental_date, r.return_date, u.user_id, u.fullname, u.address, c.car_id, c.name, p.total_price from rentals r join users u on r.user_id = u.user_id join cars c on r.car_id = c.car_id join payments p on r.rental_id = p.rental_id"
	mock.ExpectQuery(expectedSQL).WillReturnRows(all_history)

	log.Default().Println("test 123")
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		Header: make(http.Header),
	}

	adminService.GetRentalHistory(ctx)

	assert.Equal(t, http.StatusOK, ctx.Writer.Status())
	assert.Nil(t, mock.ExpectationsWereMet())
}
