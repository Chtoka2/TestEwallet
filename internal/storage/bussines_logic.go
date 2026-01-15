package storage

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (s *Storage) UserWallets(ctx context.Context, userID uuid.UUID) ([]Wallet, error){
	var wallets []Wallet
	result := s.db.WithContext(ctx).Where("user_id = ?", userID).Find(&wallets)
	if errors.Is(result.Error, gorm.ErrRecordNotFound){
		return nil, ErrWalletsNotFound
	}
	if result.Error != nil{
		return nil, result.Error
	}
	return wallets, nil
}

func (s *Storage) Transactions(ctx context.Context, userIDFrom uuid.UUID,
	userIDTo uuid.UUID, currency string, summ int64) (Transaction, error){
	var walletFrom, walletTo Wallet
	tx := s.db.Begin()
	defer tx.Rollback()
	result := tx.WithContext(ctx).Where("currency = ?", currency).Where("user_id = ?", userIDFrom).First(&walletFrom)
	if result.Error != nil{
		return Transaction{}, ErrWalletWithCurrencyNotFound
	}
	result = tx.WithContext(ctx).Where("currency = ?", currency).Where("user_id = ?", userIDTo).First(&walletTo)
	if result.Error != nil{
		return Transaction{},ErrWalletWithCurrencyNotFound
	}
	err := s.changeBalance(ctx, tx, walletFrom.ID, currency, summ, false)
	if err != nil{
		return Transaction{},err
	}
	err = s.changeBalance(ctx, tx, walletTo.ID, currency, summ, true)
	if err != nil{
		return Transaction{},err
	}
	transaction, err := add_transactions(ctx, tx, walletFrom.ID, walletTo.ID, summ)
	if err != nil{
		return Transaction{},err
	}
	tx.Commit()
	return transaction, nil
}

//If op true - add money, else write off money
// changeBalance atomically updates wallet balance.
// If op is true — adds money; if false — deducts (with balance check).
func (s *Storage) changeBalance(ctx context.Context, tx *gorm.DB, walletID uuid.UUID, currency string, amount int64, op bool) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	var result *gorm.DB
	if op {
		// Зачисление: просто добавляем
		result = tx.WithContext(ctx).Exec(
			"UPDATE wallets SET balance = balance + ? WHERE id = ? AND currency = ?",
			amount, walletID, currency,
		)
	} else {
		// Списание: только если баланс >= amount
		result = tx.WithContext(ctx).Exec(
			"UPDATE wallets SET balance = balance - ? WHERE id = ? AND currency = ? AND balance >= ?",
			amount, walletID, currency, amount,
		)
	}

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		// Либо кошелёк не найден, либо недостаточно средств (при списании)
		if !op {
			return ErrInsufficientFunds
		}
		return ErrWalletsNotFound
	}

	return nil
}

//Was maked
//First []Transaction - transactions_from wallet, second []Transaction - transations_to wallet
func (s *Storage) GetTransactions(ctx context.Context, userID uuid.UUID) ([]Transaction, []Transaction, error){
	var transactions_from []Transaction
	var transactions_to []Transaction
	userWallets, err := s.UserWallets(ctx, userID)
	if err != nil{
		return nil, nil, err
	}
	tx := s.db.Begin()
	defer tx.Rollback()
	for _, i := range userWallets{
		var localtransactions []Transaction
		result := tx.WithContext(ctx).Where("wallet_from = ?", i.ID).Find(&localtransactions)
		if result.Error != nil{
			return nil, nil, result.Error
		}
		transactions_from = append(transactions_from, localtransactions...)
		result = tx.WithContext(ctx).Where("wallet_to = ?", i.ID).Find(&localtransactions)
		if result.Error != nil{
			return nil, nil, result.Error
		}
		transactions_to = append(transactions_to, localtransactions...)
	}
	tx.Commit()
	return transactions_from, transactions_to, nil
}

func add_transactions(ctx context.Context, tx *gorm.DB, wal_from uuid.UUID,
	wal_to uuid.UUID, amount int64)(Transaction, error){
		var transaction Transaction
		transaction = Transaction{
			ID: uuid.New(),
			WalletFrom: wal_from,
			WalletTo: wal_to,
			Amount: amount,
			CreatedAt: time.Now(),
		}
		result := tx.WithContext(ctx).Create(&transaction)
		if result.Error != nil{
			return Transaction{}, result.Error
		}
		return transaction, nil
}

func (s *Storage) GetUserIDByEmail(ctx context.Context, email string) (uuid.UUID, error){
	var user User
	var ID uuid.UUID
	result := s.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil{
		return ID, ErrUserNotFound
	}
	return user.ID, nil
}