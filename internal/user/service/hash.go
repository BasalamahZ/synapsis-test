package service

import "golang.org/x/crypto/bcrypt"

// generateHash returns the generated has for the given text
// with the given salt using Bcrypt.
func generateHash(text string, salt string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.MinCost)
	if err != nil {
		return "", nil
	}
	return string(hash), nil
}

// compareHash compares the given text and the given hash
// with the given salt using Bcrypt.
func compareHash(text string, hash string, salt string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(text))
	if err != nil && err != bcrypt.ErrMismatchedHashAndPassword {
		return false, err
	}

	if err != nil {
		return false, nil
	}

	return true, nil
}
