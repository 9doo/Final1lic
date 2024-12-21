package application_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/IvanSolomatin/calc_go/aplication"
	"github.com/IvanSolomatin/calc_go/pkg/calculation"
)


// структура запроса
type RequestBody struct {
	Expression string `json:"expression"`
}
// верный запрос
func TestCalcHandler_Success(t *testing.T) {

	handler := http.HandlerFunc(application.CalcHandler)
	server := httptest.NewServer(handler)
	defer server.Close()

	// Запрос с правильным выражением
	requestBody := RequestBody{
		Expression: "1+1",
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Error marshalling request body: %v", err)
	}

	// Создание POST-запроса
	req, err := http.NewRequest("POST", server.URL+"/api/v1/calculate", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// Отправка запроса и получение ответа
	resp, err := server.Client().Do(req)
	if err != nil {
		t.Fatalf("Error sending request: %v", err)
	}

	// проверка, что статус ответа 200
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %d", resp.StatusCode)
	}

	// Проверка ответа
	var response application.Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}
	expectedResult := "2.000000"
	if response.Result != expectedResult {
		t.Fatalf("Expected result %s, got %s", expectedResult, response.Result)
	}
}

// неверное выражение
func TestCalcHandler_InvalidExpression(t *testing.T) {

	handler := http.HandlerFunc(application.CalcHandler)
	server := httptest.NewServer(handler)
	defer server.Close()

	// запрос с некорректным выражением
	requestBody := RequestBody{
		Expression: "1+/",
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Error marshalling request body: %v", err)
	}

	// Создание POST-запроса
	req, err := http.NewRequest("POST", server.URL+"/api/v1/calculate", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Отправка запроса и получение ответа
	resp, err := server.Client().Do(req)
	if err != nil {
		t.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	// проверка что статус ответа 422
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("Expected status 422, got %d", resp.StatusCode)
	}

	// проверка ответа
	var response application.Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}
	if response.Error != calc.ErrInvalidExpression.Error() {
		t.Fatalf("Expected error %v, got %v", calc.ErrInvalidExpression, response.Error)
	}
}

// Тест ошибки деления на ноль
func TestCalcHandler_DivisionByZero(t *testing.T) {
	handler := http.HandlerFunc(application.CalcHandler)
	server := httptest.NewServer(handler)
	defer server.Close()

	// запрос с делением на ноль
	requestBody := RequestBody{
		Expression: "1/0",
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Error marshalling request body: %v", err)
	}

	// Создание POST-запроса
	req, err := http.NewRequest("POST", server.URL+"/api/v1/calculate", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// отправка запроса и получение ответа
	resp, err := server.Client().Do(req)
	if err != nil {
		t.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	// проверка, что статус ответа 422
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("Expected status 422, got %d", resp.StatusCode)
	}

	// Проверка ответа
	var response application.Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}
	if response.Error != calc.ErrDivisionByZero.Error() {
		t.Fatalf("Expected error %v, got %v", calc.ErrDivisionByZero, response.Error)
	}
}

// Тест на пустое выражения
func TestCalcHandler_EmptyExpression(t *testing.T) {
	handler := http.HandlerFunc(application.CalcHandler)
	server := httptest.NewServer(handler)
	defer server.Close()

	// запрос с пустым выражением
	requestBody := RequestBody{
		Expression: "",
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Error marshalling request body: %v", err)
	}

	// Создание POST-запроса
	req, err := http.NewRequest("POST", server.URL+"/api/v1/calculate", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// отправка запроса и получение ответа
	resp, err := server.Client().Do(req)
	if err != nil {
		t.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Проверка, что статус ответа 422
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("Expected status 422, got %d", resp.StatusCode)
	}

	// проверка ответа
	var response application.Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}
	if response.Error != calc.ErrEmptyExpression.Error() {
		t.Fatalf("Expected error %v, got %v", calc.ErrEmptyExpression, response.Error)
	}
}
