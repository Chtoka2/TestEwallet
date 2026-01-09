package storage

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//Consts of errs with Storage
const(

)

//User model
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Hash      string    `gorm:"not null"`
	CreatedAt time.Time
}

type Wallet struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
	Balance   int64     `gorm:"not null;default:0"`
	Currency  string    `gorm:"not null;default:'RUB'"`
	CreatedAt time.Time
}

type Transaction struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	WalletFrom uuid.UUID `gorm:"type:uuid;index"`
	WalletTo   uuid.UUID `gorm:"type:uuid;index"`
	Amount     int64     `gorm:"not null"`
	Type       string    `gorm:"not null"`
	Status     string    `gorm:"not null;default:'pending'"`
	CreatedAt  time.Time
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