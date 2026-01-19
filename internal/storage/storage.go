package storage

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//Consts of errs with Storage
var (
	ErrEmailAlredyExists = errors.New("Email alredy exists");
	ErrShortPass = errors.New("Password is shorter than 8");
	ErrInvalidInput = errors.New("Invalid input")
	ErrUserNotFound = errors.New("User not found")
	ErrPasswordIncorrect = errors.New("Password is incorrect")
	ErrWalletsNotFound = errors.New("Wallets is not found")
	ErrWalletWithCurrencyNotFound = errors.New("Currency is not found")
	ErrCurencyNotInCurrencies = errors.New("Our service don't have this currency")
	ErrCurrencyWalletExist = errors.New("Wallet with this currency alredy exist")
	ErrInsufficientFunds = errors.New("Wallet has not enough balance to make transaction")
	ErrInvalidAmount = errors.New("Invalid amount")
	ErrInvalidExchangeRate = errors.New("Invalid exchange rate")
	ErrSameCurrency = errors.New("You cannot convert same curencies")
	ErrConversionResultZero = errors.New("conversion result cannot be zero or smaller")
)

//User model
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Hash      string    `gorm:"not null" json:"hash"`
	CreatedAt time.Time `json:"created_at"`
}

type Wallet struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Balance   int64     `gorm:"not null;default:0" json:"balance"`
	Currency  string    `gorm:"not null;default:'RUB'" json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

type Transaction struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	WalletFrom uuid.UUID `gorm:"type:uuid;index" json:"wallet_from"`
	WalletTo   uuid.UUID `gorm:"type:uuid;index" json:"wallet_to"`
	Amount     int64     `gorm:"not null" json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
}

type Storage struct{
	db *gorm.DB
}

func New(dbUrl string) (*Storage,  error){
	const op = "storage.New"
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil{
		return nil, fmt.Errorf("%s : %w", op, err)
	}
	if err := db.AutoMigrate(&Wallet{}, &Transaction{}, &User{}); err != nil{
		return  nil, fmt.Errorf("%s : %w", op, err)
	}
	return &Storage{db: db}, nil
}