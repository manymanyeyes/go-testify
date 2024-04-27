package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var cafeList = map[string][]string{
	"moscow": []string{"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"},
}

func createURL(city string, count int) string {
	return fmt.Sprintf("/cafe?count=%d&city=%s", count, city)
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
	countStr := req.URL.Query().Get("count")
	if countStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("count missing"))
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong count value"))
		return
	}

	city := req.URL.Query().Get("city")

	cafe, ok := cafeList[city]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong city value"))
		return
	}

	if count > len(cafe) {
		count = len(cafe)
	}

	answer := strings.Join(cafe[:count], ",")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(answer))
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

	msg := fmt.Sprintf("expected status code: %d, got %d", http.StatusOK, responseRecorder.Code)
	require.Equalf(t, http.StatusOK, responseRecorder.Code, msg)

	body := responseRecorder.Body.String()
	list := strings.Split(body, ",")

	assert.Len(t, list, totalCount)
}

func TestMainHandlerWhenCityIsNotSupported(t *testing.T) {
	city := "Saint Petersburg"
	count := 2

	url := createURL(city, count)
	req := httptest.NewRequest("GET", url, nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	msg := fmt.Sprintf("expected status code: %d, got %d", http.StatusBadRequest, responseRecorder.Code)
	require.Equalf(t, http.StatusBadRequest, responseRecorder.Code, msg)

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

	msg := fmt.Sprintf("expected status code: %d, got %d", http.StatusOK, responseRecorder.Code)
	require.Equalf(t, http.StatusOK, responseRecorder.Code, msg)

	body := responseRecorder.Body.String()
	assert.NotEmpty(t, body)
}
