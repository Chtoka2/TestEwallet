package storage

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (s *Storage) UserWallets(userID uuid.UUID) ([]Wallet, error){
	var wallets []Wallet
	result := s.db.Where("user_id = ?", userID).Find(&wallets)
	if errors.Is(result.Error, gorm.ErrRecordNotFound){
		return []Wallet{}, ErrWalletsNotFound
	}
	if result.Error != nil{
		return []Wallet{}, result.Error 
	}
	return wallets, nil
}

func (s *Storage) Transactions(userIDFrom uuid.UUID, userIDTo uuid.UUID, currency string, summ int64) (error){
	var walletFrom, walletTo Wallet
	result := s.db.Where("currency = ?", currency).Where("user_id = ?", userIDFrom).First(&walletFrom)
	if result.Error != nil{
		return ErrWalletWithCurrencyNotFound
	}
	result = s.db.Where("currency = ?", currency).Where("user_id = ?", userIDTo).First(&walletTo)
	if result.Error != nil{
		return ErrWalletWithCurrencyNotFound
	}
	err := change_balance(s, walletFrom.ID, currency, summ, false)
	if err != nil{
		return err
	}
	err = change_balance(s, walletTo.ID, currency, summ, true)
	return err
}

//If op true - add money, else write off money
func change_balance(s *Storage, walletID uuid.UUID, currency string, summ int64, op bool) (error){
	var wallet Wallet
	result := s.db.Where("id = ?", walletID).Where("currency = ?", currency).First(&wallet)
	if result.Error != nil{
		return result.Error
	}
	if op{
		wallet.Balance += summ
	}else{
		wallet.Balance -= summ
	}
	result = s.db.Model(&wallet).Update("balance", wallet.Balance)
	return result.Error
}