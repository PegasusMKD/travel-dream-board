package middleware_test

import (
	"fmt"
	"github.com/PegasusMKD/travel-dream-board/internal/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestErrorHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := middleware.ErrorHandler()

	testCases := []struct {
		Name   string
		Errors []error
	}{
		{Name: "No tests"},
		{Name: "Single test", Errors: []error{fmt.Errorf("Error")}},
		{Name: "Only last test", Errors: []error{fmt.Errorf("First error"), fmt.Errorf("Last test")}},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			for _, err := range tc.Errors {
				w.Code = 500
				c.Errors = append(c.Errors, &gin.Error{Err: err})
			}

			handler(c)

			if tc.Errors == nil {
				assert.Equal(t, http.StatusOK, w.Code)
			} else {
				assert.Contains(t, w.Body.String(), tc.Errors[len(tc.Errors)-1].Error())
			}
		})
	}
}
