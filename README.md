# rental-car

- tech stackaa
  Go - Gin - PostgreSQL - GORM - Swaggo - Heroku

- ERD
  ![erd rental car](https://github.com/StephenSanjaya/p2-mini-project/blob/dev/stephen/car_rental_erd.jpg)

- Web API dapat diakses pada https://tranquil-dawn-18450-e961ca3b239f.herokuapp.com/
- Swagger doc dapat diakses pada https://tranquil-dawn-18450-e961ca3b239f.herokuapp.com/swagger/index.html

- Web API memiliki endpoint sebagai berikut:

  - <b>POST</b> /api/v1/users/register
    - request body -> `{ fullname, address, email, password }`
  - <b>POST</b> /api/v1/users/login
    - request body -> `{ email, password }`
  - <b>POST</b> /api/v1/users/topup
    - request headers -> `{ authorization }`
    - request body -> `{ amount }`
  - <b>GET</b> /api/v1/cars
    - request headers -> `{ authorization }`
  - <b>GET</b> /api/v1/cars/:category_id
    - request headers -> `{ authorization }`
  - <b>POST</b> /api/v1/cars/rental
    - request headers -> `{ authorization }`
    - request body -> `{ car_id, rental_date, return_date, coupon_id }`
  - <b>POST</b> /api/v1/cars/pay/:payment_id
    - request headers -> `{ authorization }`
    - request body -> `{ payment_method_id }`
  - <b>POST</b> /api/v1/cars/return/:rental_id
    - request headers -> `{ authorization }`
  - <b>POST</b> /api/v1/admin/cars
    - request headers -> `{ authorization }`
    - request body -> `{ category_id, name, rental_cost_per_day, capacity }`
  - <b>PUT</b> /api/v1/admin/cars/:car_id
    - request headers -> `{ authorization }`
    - request body -> `{ category_id, name, rental_cost_per_day, capacity }`
  - <b>DELETE</b> /api/v1/admin/cars/:car_id
    - request headers -> `{ authorization }`
  - <b>GET</b> /api/v1/admin/users
    - request headers -> `{ authorization }`
  - <b>GET</b> /api/v1/admin/rental-history
    - request headers -> `{ authorization }`
