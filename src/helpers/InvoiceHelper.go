package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"p2-mini-project/src/dto"
	"p2-mini-project/src/entity"
)

func CreateInvoicePayment(rap *dto.RentalAndPayment, user *entity.User, car *entity.Car) (*entity.Invoice, error) {
	apiKey := os.Getenv("XENDIT_API_KEY")
	apiUrl := "https://api.xendit.co/v2/invoices"

	bodyRequest := map[string]interface{}{
		"external_id":      "1",
		"amount":           rap.Payment.TotalPrice,
		"description":      "Dummy Invoice Mini Project",
		"invoice_duration": 86400,
		"customer": map[string]interface{}{
			"name":    user.Fullname,
			"address": user.Address,
			"email":   user.Email,
		},
		"currency": "IDR",
		"items": []interface{}{
			map[string]interface{}{
				"name":     car.Name,
				"quantity": 1,
				"price":    car.RentalCostPerDay,
			},
		},
	}

	reqBody, err := json.Marshal(bodyRequest)
	if err != nil {
		return nil, err
	}

	fmt.Println("reqbody: ", string(reqBody))

	client := &http.Client{}
	request, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	fmt.Println("request: ", request)
	request.SetBasicAuth(apiKey, "")
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	fmt.Println("response: ", response)

	defer response.Body.Close()

	var resInvoice entity.Invoice
	if err := json.NewDecoder(response.Body).Decode(&resInvoice); err != nil {
		return nil, err
	}

	fmt.Println("resinvoice: ", resInvoice)
	return &resInvoice, nil

}
