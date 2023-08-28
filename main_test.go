package main

import (
	"encoding/json"
	"example.com/m/controllers"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestGetLongUrl(t *testing.T) {
	// Initialize the Gin router

	router := gin.Default()

	// Define your POST route
	router.POST("/url", AuthMiddleware(), controllers.GetLongUrl)

	// Create a JSON payload
	formData := url.Values{}
	formData.Set("long_url", "https://www.w3schools.com/css/css_form.asp")

	// Create a new HTTP request
	req := httptest.NewRequest("POST", "/url", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTMyMTYwODcsInVzZXJfaWQiOjh9.Ee5LyoI1gSGOTQkROpNoK7yGyCdIqhw7NHVGgY_WTpc")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	res := w.Body
	body, _ := ioutil.ReadAll(res)
	bodyString := string(body)
	var result struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal([]byte(bodyString), &result); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
	fmt.Println("Extracted error message:", result.Error)
	if result.Error == "Token is either expired or not active yet" {
		fmt.Println("Exception : ", result.Error)
		return
	}
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, but got %d", http.StatusOK, w.Code)
	}
	expectedBody := `{"shortUrl":"http://localhost:3000/hyvObMHG"}`
	if bodyString == "" {
		t.Errorf("Expected body '%s', but got '%s'", expectedBody, w.Body.String())
	}
}

func TestSignup(t *testing.T) {
	router := gin.Default()
	router.POST("/signup", controllers.Signup)
	jsonData := `{
	"username": "HamzssaMalic",
	"password": "Malick12334"
	}`
	req := httptest.NewRequest("POST", "/signup", strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	res := w.Body
	body, _ := ioutil.ReadAll(res)
	bodyString := string(body)
	fmt.Println(bodyString)
	var result struct {
		Error string `json:"error"`
	}

	if err := json.Unmarshal([]byte(bodyString), &result); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
	if result.Error == "In valid Username" {
		t.Errorf(result.Error)
	}
	if result.Error == "Length of password must be greater than 8" {
		t.Errorf(result.Error)
	}
	if result.Error == "User already exists" {
		fmt.Println(result.Error)
		return
	}
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, but got %d", http.StatusOK, w.Code)
	}

}

func TestLogin(t *testing.T) {
	router := gin.Default()
	router.POST("/login", controllers.Login)
	jsonData := `{
	"username": "malik_afaf3419",
	"password": "1234"
	}`
	req := httptest.NewRequest("POST", "/login", strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	res := w.Body
	body, _ := ioutil.ReadAll(res)
	bodyString := string(body)
	fmt.Println(bodyString)
	var result struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal([]byte(bodyString), &result); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
	if w.Code == 401 || result.Error == "Invalid credentials" {
		t.Errorf(result.Error)
		return
	}
	if result.Error == "Could not generate token" {
		t.Errorf(result.Error)
		return
	}
	if w.Code == http.StatusOK {
		fmt.Println(body)
	}
}

func TestRedirectingURl(t *testing.T) {
	router := gin.Default()
	router.GET("/:shortcode", controllers.RedirectingToOrignalUrl) // Replace with your actual handler

	// Create a new HTTP GET request
	req := httptest.NewRequest("GET", "/UlqcVUR8", nil)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Your assertions here, for example:
	if w.Code != http.StatusMovedPermanently {
		t.Fatalf("Expected to get status %d but instead got %d", http.StatusMovedPermanently, w.Code)
	}
	// Check if the Location header is set to the long URL
	location, ok := w.HeaderMap["Location"]
	if !ok || len(location) != 1 || location[0] != "https://htmlpdfapi.com/blog/free_html5_invoice_templates" {
		t.Fatalf("Expected the 'Location' header to be set to 'your_expected_long_url_here'")
	}
}
