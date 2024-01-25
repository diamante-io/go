package account

import "errors"

type Store interface {
	Add(a Account) error
	Delete(address string) error
	Get(address string) (Account, error)
	Update(a Account) error
	FindWithIdentityAddress(address string) ([]Account, error)
	FindWithIdentityPhoneNumber(phoneNumber string) ([]Account, error)
	FindWithIdentityEmail(email string) ([]Account, error)
	Count() (int, error)
}

var ErrNotFound = errors.New("account not found")
var ErrAlreadyExists = errors.New("account already exists")
