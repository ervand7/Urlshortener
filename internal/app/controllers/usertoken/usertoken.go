package usertoken

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

type UserToken struct {
	Key       [32]byte
	nonce     []byte
	encrypted []byte
	aesGCM    cipher.AEAD
}

func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (u *UserToken) Encode(src string) (string, error) {
	aesCipher, err := aes.NewCipher(u.Key[:])
	if err != nil {
		return "", err
	}
	u.aesGCM, err = cipher.NewGCM(aesCipher)
	if err != nil {
		return "", err
	}

	u.nonce, err = generateRandom(u.aesGCM.NonceSize())
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return "", err
	}
	u.encrypted = u.aesGCM.Seal(nil, u.nonce, []byte(src), nil) // зашифровываем
	return string(u.encrypted), nil
}
func (u *UserToken) Decode(encrypted string) (decrypted string, err error) {
	result, err := u.aesGCM.Open(nil, u.nonce, []byte(encrypted), nil) // расшифровываем
	if err != nil {
		return "", err
	}
	return string(result), nil
}
