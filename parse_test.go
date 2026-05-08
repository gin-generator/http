package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type TestRequest struct {
	UserID  int    `path:"user_id"`
	Page    int    `query:"page" validate:"required,min=1"`
	Name    string `json:"name" validate:"required"`
	Phone   string `json:"phone" validate:"phone"`
	Email   string `json:"email" validate:"omitempty,email"`
	Age     int    `json:"age" validate:"omitempty,min=1,max=150"`
	Enabled bool   `query:"enabled"`
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestParse_Success(t *testing.T) {
	router := setupRouter()
	router.POST("/users/:user_id/orders", func(c *gin.Context) {
		var req TestRequest
		if err := Parse(c, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, req)
	})

	body := `{"name":"Alice","phone":"13800138000","email":"alice@example.com","age":25}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users/123/orders?page=2&enabled=true", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestParse_ValidationError_RequiredField(t *testing.T) {
	router := setupRouter()
	router.POST("/users/:user_id/orders", func(c *gin.Context) {
		var req TestRequest
		if err := Parse(c, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, req)
	})

	// Missing required "name" field
	body := `{"phone":"13800138000"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users/123/orders?page=1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestParse_ValidationError_MinValue(t *testing.T) {
	router := setupRouter()
	router.POST("/users/:user_id/orders", func(c *gin.Context) {
		var req TestRequest
		if err := Parse(c, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, req)
	})

	// page=0 violates min=1
	body := `{"name":"Alice","phone":"13800138000"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users/123/orders?page=0", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestParse_ValidationError_CustomPhone(t *testing.T) {
	router := setupRouter()
	router.POST("/users/:user_id/orders", func(c *gin.Context) {
		var req TestRequest
		if err := Parse(c, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, req)
	})

	// Invalid phone format
	body := `{"name":"Alice","phone":"12345"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users/123/orders?page=1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestParse_ValidationError_Email(t *testing.T) {
	router := setupRouter()
	router.POST("/users/:user_id/orders", func(c *gin.Context) {
		var req TestRequest
		if err := Parse(c, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, req)
	})

	// Invalid email format
	body := `{"name":"Alice","phone":"13800138000","email":"invalid-email"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users/123/orders?page=1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestParse_ValidationError_AgeRange(t *testing.T) {
	router := setupRouter()
	router.POST("/users/:user_id/orders", func(c *gin.Context) {
		var req TestRequest
		if err := Parse(c, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, req)
	})

	// age=200 violates max=150
	body := `{"name":"Alice","phone":"13800138000","age":200}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users/123/orders?page=1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestParse_MissingQueryParam(t *testing.T) {
	router := setupRouter()
	router.POST("/users/:user_id/orders", func(c *gin.Context) {
		var req TestRequest
		if err := Parse(c, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, req)
	})

	// Missing required query param "page"
	body := `{"name":"Alice","phone":"13800138000"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users/123/orders", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}
