package auth

import "golang.org/x/crypto/bcrypt"

type PasswordManager struct{}

func (p PasswordManager) GetHash(password string) (string, error) {
	dat, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}
	return string(dat), nil
}

func (p PasswordManager) CheckHash(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
