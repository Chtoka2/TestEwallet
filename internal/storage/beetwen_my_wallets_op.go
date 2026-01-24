package storage

import (
	"context"
	"e-wallet/internal/lib/currency"
	"log/slog"
	"math"

	"github.com/google/uuid"
)

// ConvertCurrency выполняет конвертацию средств между валютными кошельками ОДНОГО пользователя.
// Списание происходит из кошелька currencyFrom, зачисление — в currencyTo по заданному курсу.
//
// Параметры:
//   - ctx: контекст
//   - userID: ID пользователя
//   - fromCurrency, toCurrency: исходная и целевая валюта (например, "RUB", "USD")
//   - amountFrom: сумма к списанию в минимальных единицах (копейках) fromCurrency
//   - rate: курс toCurrency / fromCurrency (например, 0.0111 для RUB→USD при 90 RUB/USD)
func (s *Storage) ConvertCurrency(
	ctx context.Context,
	userID uuid.UUID,
	log *slog.Logger,
	fromCurrency, toCurrency string,
	amountFrom int64,
) error {
	var rate float64
	var err error
	//get rate from central bank of Russia
	rate, err = currency.GetCBRRate(ctx, log, fromCurrency, toCurrency)
	if err != nil{
		log.Error("Can't take rate", slog.String("Error", err.Error()))
		return err
	}
	tx := s.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Находим оба кошелька пользователя
	var walletFrom, walletTo Wallet
	if err = tx.Where("user_id = ? AND currency = ?", userID, fromCurrency).First(&walletFrom).Error; err != nil {
		return ErrWalletsNotFound
	}
	if err = tx.Where("user_id = ? AND currency = ?", userID, toCurrency).First(&walletTo).Error; err != nil {
		return ErrWalletsNotFound
	}

	// Проверяем баланс и списываем
	if walletFrom.Balance < amountFrom {
		return ErrInsufficientFunds
	}
	if err = tx.Exec(
		"UPDATE wallets SET balance = balance - ? WHERE id = ?",
		amountFrom, walletFrom.ID,
	).Error; err != nil {
		return err
	}

	// Рассчитываем сумму зачисления (в копейках/центах)
	amountToFloat := float64(amountFrom) * rate
	comission := amountToFloat * 0.05
	amountTo := int64(math.Round(amountToFloat - comission))

	if amountTo <= 0 {
		return ErrConversionResultZero
	}

	// Зачисляем
	if err = tx.Exec(
		"UPDATE wallets SET balance = balance + ? WHERE id = ?",
		amountTo, walletTo.ID,
	).Error; err != nil {
		return err
	}

	// Логируем транзакцию (опционально, но рекомендуется)
	transaction := Transaction{
		ID:         uuid.New(),
		WalletFrom: walletFrom.ID,
		WalletTo:   walletTo.ID,
		Amount:     amountFrom, // или отдельные поля AmountFrom/AmountTo
	}
	if err := tx.Create(&transaction).Error; err != nil {
		return err
	}

	return tx.Commit().Error
}

