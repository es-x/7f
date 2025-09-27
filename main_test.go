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

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	for k := range cafeList {
		requests := []struct {
			count int // передаваемое значение count
			want  int // ожидаемое количество кафе в ответе
		}{
			{count: 0, want: 0},
			{count: 1, want: 1},
			{count: 2, want: 2},
			{count: 100, want: min(len(cafeList[k]), 100)},
		}

		for _, v := range requests {
			response := httptest.NewRecorder()
			count := strconv.Itoa(v.count)

			req := httptest.NewRequest("GET", "/cafe?city="+k+"&count="+count, nil)
			handler.ServeHTTP(response, req)

			require.Equal(t, http.StatusOK, response.Code)
			str := strings.Split(response.Body.String(), ",")
			if str[0] != "" {
				assert.Len(t, str, v.want)
			} else {
				assert.Len(t, str, v.want+1)
			}
		}

	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		search    string // передаваемое значение search
		wantCount int    // ожидаемое количество кафе в ответе
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/cafe?city=moscow&search="+v.search, nil)

		handler.ServeHTTP(response, req)
		require.Equal(t, http.StatusOK, response.Code)

		str := strings.ToLower(response.Body.String())
		spl := strings.Split(str, ",")

		isTrue := strings.Contains(str, v.search)

		if isTrue {
			assert.Len(t, spl, v.wantCount)
		} else {
			assert.Len(t, spl, v.wantCount+1)
		}
	}
}
