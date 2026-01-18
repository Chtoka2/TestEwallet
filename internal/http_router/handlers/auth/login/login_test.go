package login

//cd internal/http_router/handlers/auth/login  

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"log/slog"

	"e-wallet/internal/lib/jwt"
)

// Моки
type MockStorage struct {
	enterAuthFn func(ctx context.Context, email, password string) (uuid.UUID, error)
}

func (m *MockStorage) EnterAuth(ctx context.Context, email, password string) (uuid.UUID, error) {
	if m.enterAuthFn != nil {
		return m.enterAuthFn(ctx, email, password)
	}
	return uuid.Nil, errors.New("not implemented")
}

type MockJWT struct {
	generateFn func(userID uuid.UUID, ttl time.Duration) (string, error)
}

func (m *MockJWT) Generate(userID uuid.UUID, ttl time.Duration) (string, error) {
	if m.generateFn != nil {
		return m.generateFn(userID, ttl)
	}
	return "", errors.New("not implemented")
}

// Тест
func TestNew(t *testing.T) {
	// фиксированный секрет
	tests := []struct {
		name         string
		body         string
		storageMock  func() *MockStorage
		jwtMock func() *MockJWT
		wantStatus   int
		wantCookie   bool
	}{
		{
			name: "успешный вход",
			body: `{"email": "test@example.com", "password": "12345678"}`,
			storageMock: func() *MockStorage {
				return &MockStorage{
					enterAuthFn: func(ctx context.Context, email, password string) (uuid.UUID, error) {
						return uuid.MustParse("12345678-1234-5678-1234-567812345678"), nil
					},
				}
			},
			jwtMock: func() *MockJWT {
				return &MockJWT{
					generateFn: func(userID uuid.UUID, ttl time.Duration) (string, error) {
						return "fake-jwt-token", nil
					},
				}
			},
			wantStatus: 200,
			wantCookie: true,
		},
		{
			name:       "невалидный JSON",
			body:       `{"email": "test@example.com", "password": 12345}`,
			wantStatus: 500,
		},
		{
			name: "ошибка аутентификации",
			body: `{"email": "test@example.com", "password": "wrong"}`,
			storageMock: func() *MockStorage {
				return &MockStorage{
					enterAuthFn: func(ctx context.Context, email, password string) (uuid.UUID, error) {
						return uuid.Nil, errors.New("invalid")
					},
				}
			},
			wantStatus: 400,
		},
		{
			name: "Ошибка генерации JWT",
			body: `{"email": "test@example.com", "password": "12345678"}`,
			storageMock: func() *MockStorage {
				return &MockStorage{
					enterAuthFn: func(ctx context.Context, email, password string) (uuid.UUID, error) {
						return uuid.MustParse("12345678-1234-5678-1234-567812345678"), nil
					},
				}
			},
			jwtMock: func() *MockJWT {
				return &MockJWT{
					generateFn: func(userID uuid.UUID, ttl time.Duration) (string, error) {
						return "", errors.New("fjkd")
					},
				}
			},
			wantStatus: 500,
			wantCookie: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var storage EnterInterface = &MockStorage{}
			if tt.storageMock != nil {
				storage = tt.storageMock()
			}
			var jwtSvc jwt.JWTGeneratorInterface
			if tt.jwtMock != nil{
				jwtSvc = tt.jwtMock()
			}
			router := chi.NewRouter()
			router.Post("/auth/login", New(slog.Default(), storage, jwtSvc))

			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("ожидался статус %d, получен %d", tt.wantStatus, w.Code)
			}

			cookies := w.Result().Cookies()
			if tt.wantCookie && len(cookies) == 0 {
				t.Error("ожидалась кука, но её нет")
			}
			if !tt.wantCookie && len(cookies) > 0 {
				t.Errorf("кука не ожидалась, но получена: %s", cookies[0].Name)
			}
		})
	}
}