package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	city := "moscow"
	count := 10
	totalCount := len(cafeList[city])

	url := createURL(city, count)
	req := httptest.NewRequest("GET", url, nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, http.StatusOK, responseRecorder.Code)

	body := responseRecorder.Body.String()
	list := strings.Split(body, ",")

	assert.Len(t, list, totalCount)
}

func TestMainHandlerWhenCityIsNotSupported(t *testing.T) {
	city := "saint-petersburg"
	count := 2

	url := createURL(city, count)
	req := httptest.NewRequest("GET", url, nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, http.StatusBadRequest, responseRecorder.Code)

	body := responseRecorder.Body.String()
	assert.Equal(t, body, "wrong city value")
}

func TestMainHandlerWhenTheStatusIsCorrect(t *testing.T) {
	city := "moscow"
	count := 2

	url := createURL(city, count)
	req := httptest.NewRequest("GET", url, nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, http.StatusOK, responseRecorder.Code)

	body := responseRecorder.Body
	assert.NotEmpty(t, body)
}
