package hasher

import (
    "golang.org/x/crypto/bcrypt"
)

// Hashes password based on MD5 algoritm
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

//Desides if the hashed password matches to common password
func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}