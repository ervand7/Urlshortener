package usertoken

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

type UserToken struct {
	Key    [32]byte
	Nonce  []byte
	AesGCM cipher.AEAD
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
	u.AesGCM, err = cipher.NewGCM(aesCipher)
	if err != nil {
		return "", err
	}

	u.Nonce, err = generateRandom(u.AesGCM.NonceSize())
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return "", err
	}
	result := u.AesGCM.Seal(nil, u.Nonce, []byte(src), nil) // зашифровываем
	return string(result), nil
}
func (u *UserToken) Decode(Encrypted string) (decrypted string, err error) {
	result, err := u.AesGCM.Open(nil, u.Nonce, []byte(Encrypted), nil) // расшифровываем
	if err != nil {
		return "", err
	}
	return string(result), nil
}
