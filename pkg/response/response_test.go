package response_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"event-coming/internal/domain"
	"event-coming/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestResponse_Structure(t *testing.T) {
	resp := response.Response{
		Success: true,
		Data:    map[string]string{"key": "value"},
	}

	assert.True(t, resp.Success)
	assert.NotNil(t, resp.Data)
	assert.Nil(t, resp.Error)
}

func TestErrorInfo_Structure(t *testing.T) {
	errorInfo := response.ErrorInfo{
		Code:    "validation_error",
		Message: "Invalid input",
		Details: []string{"field is required"},
	}

	assert.Equal(t, "validation_error", errorInfo.Code)
	assert.Equal(t, "Invalid input", errorInfo.Message)
	assert.NotNil(t, errorInfo.Details)
}

func TestPaginationMeta_Structure(t *testing.T) {
	meta := response.PaginationMeta{
		Page:       2,
		PerPage:    10,
		Total:      55,
		TotalPages: 6,
	}

	assert.Equal(t, 2, meta.Page)
	assert.Equal(t, 10, meta.PerPage)
	assert.Equal(t, int64(55), meta.Total)
	assert.Equal(t, 6, meta.TotalPages)
}

func TestSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := map[string]string{"message": "hello"}
	response.Success(c, data)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestCreated(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := map[string]string{"id": "123"}
	response.Created(c, data)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestNoContent(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/test", nil)

	response.NoContent(c)

	assert.Equal(t, http.StatusNoContent, c.Writer.Status())
	assert.Empty(t, w.Body.String())
}

func TestError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		code       string
		message    string
	}{
		{
			name:       "Bad Request",
			statusCode: http.StatusBadRequest,
			code:       "bad_request",
			message:    "Invalid input",
		},
		{
			name:       "Unauthorized",
			statusCode: http.StatusUnauthorized,
			code:       "unauthorized",
			message:    "Not authenticated",
		},
		{
			name:       "Not Found",
			statusCode: http.StatusNotFound,
			code:       "not_found",
			message:    "Resource not found",
		},
		{
			name:       "Internal Server Error",
			statusCode: http.StatusInternalServerError,
			code:       "internal_error",
			message:    "Something went wrong",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			response.Error(c, tt.statusCode, tt.code, tt.message)

			assert.Equal(t, tt.statusCode, w.Code)

			var resp response.Response
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.False(t, resp.Success)
			assert.Equal(t, tt.code, resp.Error.Code)
			assert.Equal(t, tt.message, resp.Error.Message)
		})
	}
}

func TestValidationError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	details := map[string]string{
		"email": "Invalid email format",
		"name":  "Name is required",
	}

	response.ValidationError(c, details)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "validation_error", resp.Error.Code)
	assert.Equal(t, "Validation failed", resp.Error.Message)
	assert.NotNil(t, resp.Error.Details)
}

func TestHandleDomainError(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "Not Found",
			err:            domain.ErrNotFound,
			expectedStatus: http.StatusNotFound,
			expectedCode:   "not_found",
		},
		{
			name:           "Unauthorized",
			err:            domain.ErrUnauthorized,
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   "unauthorized",
		},
		{
			name:           "Forbidden",
			err:            domain.ErrForbidden,
			expectedStatus: http.StatusForbidden,
			expectedCode:   "forbidden",
		},
		{
			name:           "Conflict",
			err:            domain.ErrConflict,
			expectedStatus: http.StatusConflict,
			expectedCode:   "conflict",
		},
		{
			name:           "Invalid Input",
			err:            domain.ErrInvalidInput,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "invalid_input",
		},
		{
			name:           "Invalid Credentials",
			err:            domain.ErrInvalidCredentials,
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   "invalid_credentials",
		},
		{
			name:           "Token Expired",
			err:            domain.ErrTokenExpired,
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   "token_expired",
		},
		{
			name:           "Invalid Token",
			err:            domain.ErrInvalidToken,
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   "invalid_token",
		},
		{
			name:           "Unknown Error",
			err:            domain.ErrInternalServer,
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   "internal_error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			response.HandleDomainError(c, tt.err)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp response.Response
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.False(t, resp.Success)
			assert.Equal(t, tt.expectedCode, resp.Error.Code)
		})
	}
}

func TestPaginated(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := []string{"item1", "item2", "item3"}
	page := 1
	perPage := 10
	total := int64(25)

	response.Paginated(c, data, page, perPage, total)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.PaginatedResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.Meta)
	assert.Equal(t, page, resp.Meta.Page)
	assert.Equal(t, perPage, resp.Meta.PerPage)
	assert.Equal(t, total, resp.Meta.Total)
	assert.Equal(t, 3, resp.Meta.TotalPages) // 25 / 10 = 2.5, rounded up to 3
}

func TestPaginated_ExactPages(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := []string{"item1", "item2"}
	page := 2
	perPage := 5
	total := int64(10) // Exactly 2 pages

	response.Paginated(c, data, page, perPage, total)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.PaginatedResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 2, resp.Meta.TotalPages)
}

func TestPaginated_SinglePage(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := []string{"item1", "item2", "item3"}
	page := 1
	perPage := 10
	total := int64(3) // Less than perPage

	response.Paginated(c, data, page, perPage, total)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.PaginatedResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 1, resp.Meta.TotalPages)
}

func TestPaginated_EmptyResult(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := []string{}
	page := 1
	perPage := 10
	total := int64(0)

	response.Paginated(c, data, page, perPage, total)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.PaginatedResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 0, resp.Meta.TotalPages)
}
