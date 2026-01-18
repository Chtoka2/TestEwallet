package login
//internal/http_router/handlers/auth/login
import (
	"context"
	"e-wallet/internal/lib/jwt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"bytes"
)

type Stub struct{}

func (s *Stub) EnterAuth(ctx context.Context, email string, password string) (uuid.UUID, error){
	return uuid.New(), nil
}

func TestHandler(t *testing.T){
	// Настройка
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "test-secret"
	}
	jwtSvc := jwt.NewJWTService(secretKey)

	router := chi.NewRouter()
	router.Post("/auth/login", New(slog.Default(), &Stub{}, jwtSvc))

	// Создаём запрос
	bodies := []string{
		`{"email": "test@example.com", "password": "12345678"}`,
		`{"email": "test@example.com"}`,
	} 
	for _, i := range bodies {
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(i))


		req.Header.Set("Content-Type", "application/json")

		// Записываем ответ
		w := httptest.NewRecorder()

		// Вызываем хендлер
		router.ServeHTTP(w, req)

		// Проверяем, что ответ пришёл
		if w.Code == 0 {
			t.Fatal("хендлер не вернул статус")
		}

		// Успешный вход должен вернуть 200
		if w.Code != http.StatusOK {
			t.Errorf("ожидался статус 200, получен %d", w.Code)
		}

		// Проверяем, что тело не пустое
		if w.Body.Len() == 0 {
			t.Error("тело ответа пустое")
		}
	}
}