package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

// function for finding the minimum
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	city := "moscow"
	numberOfCafes := len(cafeList[city])

	requests := []struct {
		count int // the transmitted value
		want  int // expected quantity
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, Min(numberOfCafes, 100)},
	}

	for _, value := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/cafe?count="+strconv.Itoa(value.count)+"&city="+city, nil)
		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code) // checking the request processing

		var RespSplit []string
		StrResponse := response.Body.String()

		if StrResponse == "" {
			RespSplit = []string{}
		} else {
			RespSplit = strings.Split(strings.TrimSpace(StrResponse), ",") // converting the response body into a slice
		}

		assert.Equal(t, value.want, len(RespSplit))
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	city := "moscow"

	requests := []struct {
		search    string
		wantCount int
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}

	for _, value := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/cafe?city="+city+"&search="+value.search, nil)
		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code) // checking the request processing

		var RespSplit []string
		StrResponse := response.Body.String()

		if StrResponse == "" {
			RespSplit = []string{}
		} else {
			RespSplit = strings.Split(strings.TrimSpace(StrResponse), ",") // converting the response body into a slice
		}

		assert.Equal(t, value.wantCount, len(RespSplit)) // checking the number

		for _, v := range RespSplit {
			assert.True(t, strings.Contains(strings.ToLower(v), strings.ToLower(value.search)))
		}

	}
}
