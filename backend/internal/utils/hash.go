package utils

import "golang.org/x/crypto/bcrypt"

// HashearPassword genera un hash bcrypt a partir de una password en texto plano.
func HashearPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerificarPassword compara una password en texto plano contra un hash bcrypt.
// Retorna nil si coinciden, o un error si no.
func VerificarPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
