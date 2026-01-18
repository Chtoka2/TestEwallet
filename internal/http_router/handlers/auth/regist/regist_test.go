package regist

//cd internal/http_router/handlers/auth/regist

import (
	"bytes"
	"context"
	"e-wallet/internal/storage"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type MockStorage struct {
	enterAuthFn func(ctx context.Context, email, password string) (uuid.UUID, error)
}

func (m *MockStorage) RegistAUTH(ctx context.Context, email string, password string) (uuid.UUID, error) {
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

func Test(t *testing.T){
	tests := []struct{
		name string
		body string
		storageMock  func() *MockStorage
		jwtMock func() *MockJWT
		wantStatus int
		wantCookie bool
	}{
		{
			name: "Успешная регистрация",
			body: `{"email":"test@mail.ru", "password":"12345678"}`,
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
			name: "Невалидный json",
			body: `{kjsdkjs:kjksjdf, "kldskd":klkd}`,
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
			wantStatus: 400,
			wantCookie: false,
		},
		{
			name: "Ошибка регистрации(email существует)",
			body: `{"email":"test@mail.ru", "password":"12345678"}`,
			storageMock: func() *MockStorage {
				return &MockStorage{
					enterAuthFn: func(ctx context.Context, email, password string) (uuid.UUID, error) {
						return uuid.Nil, storage.ErrEmailAlredyExists
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
			wantStatus: 409,
			wantCookie: false,
		},
		{
			name: "Ошибка данных",
			body: `{"email":"test@mail.ru", "password":"12345678"}`,
			storageMock: func() *MockStorage {
				return &MockStorage{
					enterAuthFn: func(ctx context.Context, email, password string) (uuid.UUID, error) {
						return uuid.Nil, storage.ErrInvalidInput
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
			wantStatus: 400,
			wantCookie: false,
		},
		{
			name: "Непонятная ошибка данных",
			body: `{"email":"test@mail.ru", "password":"12345678"}`,
			storageMock: func() *MockStorage {
				return &MockStorage{
					enterAuthFn: func(ctx context.Context, email, password string) (uuid.UUID, error) {
						return uuid.Nil, errors.New("jdjffkj")
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
			wantStatus: 500,
			wantCookie: false,
		},
		{
			name: "Ошибка генерации JWT токена",
			body: `{"email":"test@mail.ru", "password":"12345678"}`,
			storageMock: func() *MockStorage {
				return &MockStorage{
					enterAuthFn: func(ctx context.Context, email, password string) (uuid.UUID, error) {
						return uuid.New(), nil
					},
				}
			},
			jwtMock: func() *MockJWT {
				return &MockJWT{
					generateFn: func(userID uuid.UUID, ttl time.Duration) (string, error) {
						return "", errors.New("Some problem")
					},
				}
			},
			wantStatus: 500,
			wantCookie: false,
		},
	}

	//Test logic
	for _, i := range tests{
		t.Run(i.name, func(t *testing.T) {
			storage := i.storageMock()
			jwtSvc := i.jwtMock()
			r := chi.NewRouter()
			r.Post("/auth/regist", New(slog.Default(), storage, jwtSvc))
			req := httptest.NewRequest(http.MethodPost, "/auth/regist", bytes.NewBufferString(i.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != i.wantStatus {
				t.Errorf("ожидался статус %d, получен %d", i.wantStatus, w.Code)
			}

			cookies := w.Result().Cookies()
			if i.wantCookie && len(cookies) == 0 {
				t.Error("ожидалась кука, но её нет")
			}
			if !i.wantCookie && len(cookies) > 0 {
				t.Errorf("кука не ожидалась, но получена: %s", cookies[0].Name)
			}
		})
	}
}