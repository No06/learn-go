package crypto

import "golang.org/x/crypto/bcrypt"

// HashPassword hashes a raw password with bcrypt.
func HashPassword(raw string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// ComparePassword compares a hashed password with the provided input.
func ComparePassword(hashed, candidate string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(candidate))
}
