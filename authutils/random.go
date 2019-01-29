package authutils

import "crypto/rand"

// GenerateRandomBytes generates random bytes
func GenerateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// // GenerateRandomString generates a random string
// func GenerateRandomBase64String(len uint32) (string, error) {
// 	randomBytes, err := GenerateRandomBytes(len)
// 	if err != nil {
// 		return "", err
// 	}

// 	``
// }
