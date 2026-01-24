package ErrHandler

import (
	"e-wallet/internal/storage"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

var(
	ErrFailedDecodeJSON = errors.New("Failed to decode JSON")
)

type Response struct{
	Status string
	Error string
}

func ErrHandler(
	w http.ResponseWriter,
	r *http.Request,
	log *slog.Logger,
	err error,
) {
	switch {
	// === 4xx: Ошибки клиента ===
	case errors.Is(err, storage.ErrWalletsNotFound):
		// Кошелёк(и) не найдены — клиент запросил несуществующий ресурс.
		// Аналогично "user not found", но для wallet → 404.
		log.Error("Wallets not found", slog.Any("error", err))
		Responser(w, r, err.Error(), http.StatusNotFound)
		return

	case errors.Is(err, storage.ErrWalletWithCurrencyNotFound):
		// У пользователя нет кошелька в указанной валюте — это ошибка запроса.
		// Ресурс (wallet+currency) отсутствует → 404.
		log.Error("Wallet with currency not found", slog.Any("error", err))
		Responser(w, r, err.Error(), http.StatusNotFound)
		return

	case errors.Is(err, storage.ErrInsufficientFunds):
		// Недостаточно средств — семантически корректный запрос, но условие не выполнено.
		// Используем 400 (Bad Request), т.к. операция невозможна при текущем состоянии.
		// Альтернатива: 422 (Unprocessable Entity), но 400 проще и широко принят.
		log.Error("Insufficient funds", slog.Any("error", err))
		Responser(w, r, err.Error(), http.StatusBadRequest)
		return

	case errors.Is(err, storage.ErrSameCurrency):
		// Попытка конвертации в ту же валюту — логическая ошибка клиента.
		// Запрос бессмыслен → 400.
		log.Error("Same currency conversion attempted", slog.Any("error", err))
		Responser(w, r, err.Error(), http.StatusBadRequest)
		return

	case errors.Is(err, storage.ErrInvalidAmount):
		// Сумма <= 0 или нецелая при работе с minor units — некорректный ввод.
		// Классическая ошибка валидации → 400.
		log.Error("Invalid amount", slog.Any("error", err))
		Responser(w, r, err.Error(), http.StatusBadRequest)
		return

	case errors.Is(err, storage.ErrEmailAlredyExists):
		// Email уже занят — нарушение уникальности при регистрации.
		// Стандарт для "resource already exists" → 409 Conflict.
		log.Error("Email already exists", slog.Any("error", err))
		Responser(w, r, err.Error(), http.StatusConflict)
		return

	case errors.Is(err, storage.ErrShortPass):
		// Пароль не соответствует требованиям — ошибка валидации.
		// → 400 Bad Request.
		log.Error("Password too short", slog.Any("error", err))
		Responser(w, r, err.Error(), http.StatusBadRequest)
		return

	case errors.Is(err, storage.ErrInvalidInput):
		// Общая ошибка валидации входных данных (email, имя и т.п.).
		// → 400.
		log.Error("Invalid input", slog.Any("error", err))
		Responser(w, r, err.Error(), http.StatusBadRequest)
		return

	case errors.Is(err, storage.ErrPasswordIncorrect):
		// Неверный пароль при входе — аутентификация не удалась.
		// Стандарт: 401 Unauthorized (не путать с 403!).
		log.Error("Incorrect password", slog.Any("error", err))
		Responser(w, r, err.Error(), http.StatusUnauthorized)
		return

	case errors.Is(err, storage.ErrCurencyNotInCurrencies):
		// Валюта не поддерживается системой — клиент отправил недопустимое значение.
		// → 400 Bad Request.
		log.Error("Currency not supported", slog.Any("error", err))
		Responser(w, r, err.Error(), http.StatusBadRequest)
		return

	case errors.Is(err, storage.ErrCurrencyWalletExist):
		// Попытка создать дубликат кошелька в той же валюте.
		// Ресурс уже существует → 409 Conflict.
		log.Error("Wallet in this currency already exists", slog.Any("error", err))
		Responser(w, r, err.Error(), http.StatusConflict)
		return

	// === 5xx: Ошибки сервера ===
	case errors.Is(err, ErrFailedDecodeJSON):
		log.Error("Wallets not found", slog.Any("error", err))
		Responser(w, r, err.Error(), http.StatusInternalServerError)
		return
	case errors.Is(err, storage.ErrInvalidExchangeRate):
		// Курс = 0 или отрицательный — внутренняя ошибка расчёта.
		// Клиент ни при чём → 500.
		log.Error("Invalid exchange rate from provider", slog.Any("error", err))
		Responser(w, r, "Exchange rate temporarily unavailable", http.StatusInternalServerError)
		return

	case errors.Is(err, storage.ErrConversionResultZero):
		// Результат конвертации = 0 из-за округления или курса — техническая проблема.
		// Не должно происходить при корректной логике → 500.
		log.Error("Conversion resulted in zero", slog.Any("error", err))
		Responser(w, r, "Conversion failed due to rate precision", http.StatusInternalServerError)
		return

	default:
		// Любая неожиданная ошибка — логируем как внутреннюю.
		log.Error("Unexpected internal error", slog.Any("error", err))
		Responser(w, r, "Internal server error", http.StatusInternalServerError)
		return
	}
}
func Responser(
	w http.ResponseWriter,
	r *http.Request,
	textOfError string,
	statusCode int,
){
	render.Status(r, statusCode)
	render.JSON(w,r, Response{
		Status: "Error",
		Error: textOfError,
	})
}