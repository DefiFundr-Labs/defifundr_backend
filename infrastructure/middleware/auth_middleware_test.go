package middleware

// import (
// 	"errors"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// 	"time"

// 	tokenMaker "github.com/demola234/defifundr/pkg/token_maker"

// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// // MockTokenMaker is a mock for the token.Maker interface
// type MockTokenMaker struct {
// 	mock.Mock
// }

// // Mock the VerifyToken method to satisfy the token.Maker interface
// func (m *MockTokenMaker) VerifyToken(token string) (*tokenMaker.Payload, error) {
// 	args := m.Called(token)
// 	if payload, ok := args.Get(0).(*tokenMaker.Payload); ok {
// 		return payload, args.Error(1)
// 	}
// 	return nil, args.Error(1)
// }

// // Mock the CreateToken method to satisfy the token.Maker interface
// func (m *MockTokenMaker) CreateToken(email string, userID uuid.UUID, duration time.Duration) (string, *tokenMaker.Payload, error) {
// 	args := m.Called(email, userID, duration)
// 	if payload, ok := args.Get(1).(*tokenMaker.Payload); ok {
// 		return args.String(0), payload, args.Error(2)
// 	}
// 	return args.String(0), nil, args.Error(2)
// }

// func TestAuthMiddleware(t *testing.T) {
// 	gin.SetMode(gin.TestMode)
// 	mockTokenMaker := new(MockTokenMaker)
// 	authMiddleware := AuthMiddleware(mockTokenMaker)

// 	router := gin.New()
// 	router.Use(authMiddleware)
// 	router.GET("/test", func(ctx *gin.Context) {
// 		if payload, exists := ctx.Get(authorizationPayloadKey); exists {
// 			ctx.JSON(http.StatusOK, gin.H{"message": "success", "user_id": payload.(*tokenMaker.Payload).UserID}) // Replace 'ID' with the correct field name
// 		} else {
// 			ctx.JSON(http.StatusOK, gin.H{"message": "payload not found"})
// 		}
// 	})

// 	t.Run("missing authorization header", func(t *testing.T) {
// 		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
// 		w := httptest.NewRecorder()
// 		router.ServeHTTP(w, req)

// 		assert.Equal(t, http.StatusUnauthorized, w.Code)
// 		assert.Contains(t, w.Body.String(), "authorization header not found")
// 	})

// 	t.Run("invalid authorization header format", func(t *testing.T) {
// 		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
// 		req.Header.Set(authorizationHeader, "bearer")
// 		w := httptest.NewRecorder()
// 		router.ServeHTTP(w, req)

// 		assert.Equal(t, http.StatusUnauthorized, w.Code)
// 		assert.Contains(t, w.Body.String(), "invalid authorization header format")
// 	})

// 	t.Run("unsupported authorization type", func(t *testing.T) {
// 		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
// 		req.Header.Set(authorizationHeader, "basic token")
// 		w := httptest.NewRecorder()
// 		router.ServeHTTP(w, req)

// 		assert.Equal(t, http.StatusUnauthorized, w.Code)
// 		assert.Contains(t, w.Body.String(), "unsupported authorization type basic")
// 	})

// 	t.Run("invalid token", func(t *testing.T) {
// 		mockTokenMaker.On("VerifyToken", "invalid_token").Return(nil, errors.New("invalid token")).Once()

// 		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
// 		req.Header.Set(authorizationHeader, "bearer invalid_token")
// 		w := httptest.NewRecorder()
// 		router.ServeHTTP(w, req)

// 		assert.Equal(t, http.StatusUnauthorized, w.Code)
// 		assert.Contains(t, w.Body.String(), "invalid token")
// 		mockTokenMaker.AssertExpectations(t)
// 	})

// 	t.Run("valid token", func(t *testing.T) {
// 		payload := &tokenMaker.Payload{Email: "12345", UserID: uuid.New()} // Replace 'ID' with the actual field name
// 		mockTokenMaker.On("VerifyToken", "valid_token").Return(payload, nil).Once()

// 		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
// 		req.Header.Set(authorizationHeader, "bearer valid_token")
// 		w := httptest.NewRecorder()
// 		router.ServeHTTP(w, req)

// 		assert.Equal(t, http.StatusOK, w.Code)
// 		assert.Contains(t, w.Body.String(), `"user_id":"12345"`)
// 		mockTokenMaker.AssertExpectations(t)
// 	})
// }
