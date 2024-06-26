package helpers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"p2-mini-project/src/entity"
)

func CreateInvoiceRental(totalPrice *float64, user *entity.User, car *entity.Car) (*entity.Invoice, error) {
	apiKey := os.Getenv("XENDIT_API_KEY")
	apiUrl := "https://api.xendit.co/v2/invoices"

	bodyRequest := map[string]interface{}{
		"external_id":      "1",
		"amount":           totalPrice,
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

	client := &http.Client{}
	request, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	request.SetBasicAuth(apiKey, "")
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var resInvoice entity.Invoice
	if err := json.NewDecoder(response.Body).Decode(&resInvoice); err != nil {
		return nil, err
	}

	return &resInvoice, nil
}

func CreateInvoiceTopUp(user *entity.User, amount float64) (*entity.Invoice, error) {
	apiKey := os.Getenv("XENDIT_API_KEY")
	apiUrl := "https://api.xendit.co/v2/invoices"

	bodyRequest := map[string]interface{}{
		"external_id":      "1",
		"amount":           amount,
		"description":      "Dummy Invoice Mini Project",
		"invoice_duration": 86400,
		"customer": map[string]interface{}{
			"name":    user.Fullname,
			"address": user.Address,
			"email":   user.Email,
		},
		"currency": "IDR",
	}

	reqBody, err := json.Marshal(bodyRequest)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	request, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	request.SetBasicAuth(apiKey, "")
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var resInvoice entity.Invoice
	if err := json.NewDecoder(response.Body).Decode(&resInvoice); err != nil {
		return nil, err
	}

	return &resInvoice, nil
}
