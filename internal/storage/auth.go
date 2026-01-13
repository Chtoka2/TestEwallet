package storage

import (
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
	
	if err = s.CreateEWallet(user.ID, RUB); err != nil{
		return ID, err
	}
	ID = user.ID
	tx.Commit()
	return ID, nil
}
//TODO: When I make all, rewrite it with config data, to take settings when creating
//For example: bonus of money, I can set it in config.
//It would be better but now this func is enough
func (s *Storage) CreateEWallet(UserID uuid.UUID, currency string) (error){
	var bonus int64
	wallets, err := s.UserWallets(UserID)
	if err != nil{
		return err
	}
	for _, i := range wallets{
		if i.Currency == currency{
			return ErrCurrencyWalletExist
		}
	}
	switch currency{
	case USD:
		bonus = 1000
	case RUB:
		bonus = 100000
	case CNY:
		bonus = 100
	case EUR:
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
	result := s.db.Create(&ewallet)
	if result.Error != nil{
		return result.Error
	}
	return nil
}

func (s *Storage) EnterAuth(email string, password string) (uuid.UUID, error){
	tx := s.db.Begin()
	defer tx.Rollback()
	var user User
	var ID uuid.UUID
	result := s.db.Where("email = ?", email).First(&user)
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

//I think it now useful, because many people can cheat. They can just open wallet make transaction to friend
//then they can delete wallet and create wallet again
// I dont delete it, because it can be useful in future
func (s *Storage) DeleteWallet(userID uuid.UUID, currency string) error{
	var wallet Wallet
	wallets, err := s.UserWallets(userID)
	if err != nil{
		return err
	}
	for x, i := range wallets{
		if i.Currency == currency{
			break;
		}
		if x +1 == len(wallets){
			return ErrWalletWithCurrencyNotFound
		}
	}
	result := s.db.Where("user_id = ?", userID).Where("currency = ?", currency).First(&wallet)
	if result.Error != nil{
		return result.Error
	}
	result = s.db.Delete(&wallet, wallet.ID)
	if result.Error != nil{
		return result.Error
	}
	return nil
}