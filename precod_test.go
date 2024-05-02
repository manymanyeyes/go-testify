package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createURL(city string, count int) string {
	return fmt.Sprintf("/cafe?count=%d&city=%s", count, city)
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	city := "moscow"
	count := 10
	totalCount := len(cafeList[city])

	url := createURL(city, count)
	req := httptest.NewRequest("GET", url, nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, http.StatusOK, responseRecorder.Code, "Expected status code OK")

	body := responseRecorder.Body.String()
	list := strings.Split(body, ",")

	assert.Len(t, list, totalCount, "Expected number of cafes to match total count")
}

func TestMainHandlerWhenCityIsNotSupported(t *testing.T) {
	city := "saint-petersburg"
	count := 2

	url := createURL(city, count)
	req := httptest.NewRequest("GET", url, nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, http.StatusBadRequest, responseRecorder.Code, "Expected status code BadRequest")

	body := responseRecorder.Body.String()
	assert.Equal(t, body, "wrong city value", "Expected error message for unsupported city")
}

func TestMainHandlerWhenTheStatusIsCorrect(t *testing.T) {
	city := "moscow"
	count := 2

	url := createURL(city, count)
	req := httptest.NewRequest("GET", url, nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, http.StatusOK, responseRecorder.Code, "Expected status code OK")

	body := responseRecorder.Body
	assert.NotEmpty(t, body, "Expected non-empty response body")
}
