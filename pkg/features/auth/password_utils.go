package auth

import "golang.org/x/crypto/bcrypt"

// HashPassword will take a password and hash them using the bcrypt algorithm
// the max length accepted is 72 bytes based on the [bcrypt.GenerateFromPassword] documentation
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// PasswordMatch will take a hash and a password and will compare if the hash is generated with
// the current password
func PasswordMatch(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
