package storage

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *Storage) RegistAUTH(email string, password string) (uuid.UUID, error){
	var ID uuid.UUID
	if len(password) < 8{
		return ID, ErrShortPass
	}
	if email == "" || password == "" {
    return ID, ErrInvalidInput
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil{
		return ID, err
	}
	var exists int64
	s.db.Model(&User{}).Where("email = ?", email).Count(&exists)
	if exists != 0{
		return ID, ErrEmailAlredyExists
	}
	user := User{
		ID: uuid.New(),
		Email: email,
		Hash: string(hash),
		CreatedAt: time.Now(),
	}
	tx := s.db.Begin()
	defer tx.Rollback()
	result := s.db.Create(&user)
	if result.Error != nil{
		return  ID, result.Error
	}
	
	if err = s.CreateEWallet(user.ID); err != nil{
		return ID, err
	}
	ID = user.ID
	tx.Commit()
	return ID, nil
}

func (s *Storage) CreateEWallet(UserID uuid.UUID) (error){
	ewallet := Wallet{
		ID: uuid.New(),
		UserID: UserID,
		Balance: 1000, //Just bonus from my "company"
		Currency: "RUB",
		CreatedAt: time.Now(),
	}
	result := s.db.Create(&ewallet)
	if result.Error != nil{
		return result.Error
	}
	return nil
}