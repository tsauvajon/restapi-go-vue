package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var app App

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS products
(
id SERIAL,
name TEXT NOT NULL,
price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
CONSTRAINT products_pkey PRIMARY KEY (id)
)`

func ensureTableExists() {
	if _, err := app.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	app.DB.Exec("DELETE FROM products")
	app.DB.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

func addProducts(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		app.DB.Exec("INSERT INTO products(name, price) VALUES($1, $2)", "Product "+strconv.Itoa(i), (i+1.0)*10)
	}
}

func TestMain(m *testing.M) {
	app = App{}
	os.Setenv("TEST_DB_USERNAME", "testing")
	os.Setenv("TEST_DB_PASSWORD", "testing")
	os.Setenv("TEST_DB_NAME", "restapi-go-vue")
	os.Setenv("TEST_DB_HOST", "localhost")

	app.Initialize(
		os.Getenv("TEST_DB_USERNAME"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_NAME"),
		os.Getenv("TEST_DB_HOST"),
		"5432",
		"disable")

	ensureTableExists()
	code := m.Run()

	clearTable()

	os.Exit(code)
}

func TestEmptyTable(t *testing.T) {
	clearTable()
	req, _ := http.NewRequest("GET", "/api/products", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestNonExistantProduct(t *testing.T) {
	clearTable()
	req, _ := http.NewRequest("GET", "/api/products/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Product not found" {
		t.Errorf("Expected Product not found, got %s", m["error"])
	}
}

func TestCreateProduct(t *testing.T) {
	clearTable()

	productName := "test product"
	productPrice := 45.67

	payload := []byte(`{"name": "` + productName + `", "price": ` + fmt.Sprintf("%f", productPrice) + `}`)
	req, _ := http.NewRequest("POST", "/api/products", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var product map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &product)

	if product["name"] != productName {
		t.Errorf("Expected product name to be %s, got %v", productName, product["name"])
	}

	if product["price"] != productPrice {
		t.Errorf("Expected product price to be %f, got %v", productPrice, product["price"])
	}

	if product["id"] != 1.0 {
		t.Errorf("Expected product ID to be '1', got %v", product["id"])
	}
}

func TestGetProduct(t *testing.T) {
	clearTable()
	addProducts(1)

	req, _ := http.NewRequest("GET", "/api/products/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateProduct(t *testing.T) {
	clearTable()
	addProducts(1)

	productName := "test updated product"
	productPrice := 78.78

	req, _ := http.NewRequest("GET", "/api/products/1", nil)
	response := executeRequest(req)
	var originalProduct map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalProduct)

	payload := []byte(`{"name": "` + productName + `", "price": ` + fmt.Sprintf("%f", productPrice) + `}`)
	req, _ = http.NewRequest("PUT", "/api/products/1", bytes.NewBuffer(payload))
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var updatedProduct map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &updatedProduct)

	if updatedProduct["id"] != originalProduct["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalProduct["id"], updatedProduct["id"])
	}

	if updatedProduct["name"] == originalProduct["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalProduct["name"], productName, updatedProduct["name"])
	}

	if updatedProduct["price"] == originalProduct["price"] {
		t.Errorf("Expected the price to change from '%v' to '%v'. Got '%v'", originalProduct["price"], productPrice, updatedProduct["price"])
	}
}

func TestDeleteProduct(t *testing.T) {
	clearTable()
	addProducts(1)

	req, _ := http.NewRequest("GET", "/api/products/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/api/products/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/api/products/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}
