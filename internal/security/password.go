package security

import "golang.org/x/crypto/bcrypt"

type BcryptPasswordEncoder struct{}

func NewBcryptPasswordEncoder() PasswordEncoder {
	return &BcryptPasswordEncoder{}
}

func (e *BcryptPasswordEncoder) Encode(rawPassword string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (e *BcryptPasswordEncoder) Verify(hashedPassword, rawPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(rawPassword))
	return err == nil
}
