package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var currencys []string = []string{RUB, USD, CNY, EUR}

const(
	RUB = "RUB"
	USD = "USD"
	CNY = "CNY"
	EUR = "EUR"
)

//Func registration user, I use it to register user in DB
func (s *Storage) RegistAUTH(ctx context.Context, email string, password string) (uuid.UUID, error){
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
	tx := s.db.Begin()
	defer tx.Rollback()
	s.db.WithContext(ctx).Model(&User{}).Where("email = ?", email).Count(&exists)
	if exists != 0{
		return ID, ErrEmailAlredyExists
	}
	user := User{
		ID: uuid.New(),
		Email: email,
		Hash: string(hash),
		CreatedAt: time.Now(),
	}
	result := tx.WithContext(ctx).Create(&user)
	if result.Error != nil{
		return  ID, result.Error
	}
	
	if err = s.CreateEWallet(ctx, user.ID, RUB); err != nil{
		return ID, err
	}
	ID = user.ID
	tx.Commit()
	return ID, nil
}
//TODO: When I make all, rewrite it with config data, to take settings when creating
//For example: bonus of money, I can set it in config.
//It would be better but now this func is enough
func (s *Storage) CreateEWallet(ctx context.Context, UserID uuid.UUID, currency string) (error){
	var bonus int64
	tx := s.db.Begin()
	defer tx.Rollback()
	wallets, err := s.UserWallets(ctx, UserID)
	if err != nil{
		return err
	}
	for _, i := range wallets{
		if i.Currency == currency{
			return ErrCurrencyWalletExist
		}
	}
	switch currency{
	case USD: // a cent
		bonus = 1000
	case RUB: // a penny
		bonus = 100000
	case CNY: // a fyn
		bonus = 1000
	case EUR: // a eurocent
		bonus = 1000
	default:
		return ErrCurencyNotInCurrencies
	}
	ewallet := Wallet{
		ID: uuid.New(),
		UserID: UserID,
		Balance: bonus, //Just bonus from my "company"
		Currency: currency,
		CreatedAt: time.Now(),
	}
	result := tx.WithContext(ctx).Create(&ewallet)
	if result.Error != nil{
		return result.Error
	}
	tx.Commit()
	return nil
}

func (s *Storage) EnterAuth(ctx context.Context, email string, password string) (uuid.UUID, error){
	tx := s.db.Begin()
	defer tx.Rollback()
	var user User
	var ID uuid.UUID
	result := tx.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil{
		return ID, result.Error
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password)); err != nil{
		return ID, ErrPasswordIncorrect
	}
	tx.Commit()
	return user.ID, nil
}

func GetCurrencyes() ([]string){
	return currencys
}

//I think it now unuseful, because many people can cheat. They can just open wallet make transaction to friend
//then they can delete wallet and create wallet again
// I dont delete it, because it can be useful in future
func (s *Storage) DeleteWallet(ctx context.Context, userID uuid.UUID, currency string) error{
	var wallet Wallet
	tx := s.db.Begin()
	defer tx.Rollback()
	result := tx.WithContext(ctx).Where("user_id = ?", userID).Where("currency = ?", currency).First(&wallet)
	if result.Error != nil{
		return result.Error
	}
	result = tx.WithContext(ctx).Delete(&wallet, wallet.ID)
	if result.Error != nil{
		return result.Error
	}
	tx.Commit()
	return nil
}