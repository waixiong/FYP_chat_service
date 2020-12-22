package token

import (
	"crypto/aes"
	"crypto/cipher"
	c_rand "crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	m_rand "math/rand"
	"time"
	"unsafe"
)

// GenerateAccessToken generate token
func GenerateAccessToken(user string, refresh string) (string, error) {
	// expire in an hour
	token := &Get_Token{User: user, Exp: time.Now().Add(1 * time.Hour).Unix(), Refresh: refresh}
	rawToken := []byte(token.toString())
	key := []byte("passphrasewhichneedstobe32bytes!")

	// generate a new aes cipher using our 32 byte long key
	c, err := aes.NewCipher(key)
	if err != nil {
		// err in generate token
		fmt.Println(err)
		return "", err
	}
	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err := cipher.NewGCM(c)
	// if any error generating new GCM
	// handle them
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(c_rand.Reader, nonce); err != nil {
		fmt.Println(err)
		return "", err
	}

	// here we encrypt our text using the Seal function
	// Seal encrypts and authenticates plaintext, authenticates the
	// additional data and appends the result to dst, returning the updated
	// slice. The nonce must be NonceSize() bytes long and unique for all
	// time, for a given key.
	encoder := base64.StdEncoding
	// fmt.Println(enc)
	encodeByte := gcm.Seal(nonce, nonce, rawToken, nil)
	// fmt.Println(encByte)
	encodeToken := encoder.EncodeToString(encodeByte)
	return "token_" + encodeToken, nil
}

// GenerateRefreshToken generate token
func GenerateRefreshToken() (string, string) {
	// return "refresh_" + randStringBytesMaskImprSrcUnsafe(22)
	r := randStringBytesMaskImprSrcUnsafe(12)
	return "refresh_" + r, r
}

func RawTokenFromAccessToken(accessToken string) (*Get_Token, error) {
	decoder := base64.StdEncoding
	decByte, err := decoder.DecodeString(accessToken[6:])
	// if our program was unable to read the file
	// print out the reason why it can't
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	key := []byte("passphrasewhichneedstobe32bytes!")
	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(decByte) < nonceSize {
		fmt.Println(err)
		return nil, err
	}

	nonce, ciphertext := decByte[:nonceSize], decByte[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(len(plaintext))
	fmt.Println(string(plaintext))
	fmt.Println("")
	token, err := get_token_FromString(string(plaintext))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return token, nil
}

// ---------------- Random String Genarator ---------------- //
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = m_rand.NewSource(time.Now().UnixNano())

// source : https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func randStringBytesMaskImprSrcUnsafe(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

// Service one time token
func GenerateOneTimeToken(user string) (string, error) {
	// expire in an hour
	token := &Get_Token{User: user, Exp: time.Now().Add(1 * time.Hour).Unix(), Refresh: ""}
	rawToken := []byte(token.toString())
	key := []byte("passphrasewhichneedstobe32bytes!")

	// generate a new aes cipher using our 32 byte long key
	c, err := aes.NewCipher(key)
	if err != nil {
		// err in generate token
		fmt.Println(err)
		return "", err
	}
	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err := cipher.NewGCM(c)
	// if any error generating new GCM
	// handle them
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(c_rand.Reader, nonce); err != nil {
		fmt.Println(err)
		return "", err
	}

	// here we encrypt our text using the Seal function
	// Seal encrypts and authenticates plaintext, authenticates the
	// additional data and appends the result to dst, returning the updated
	// slice. The nonce must be NonceSize() bytes long and unique for all
	// time, for a given key.
	encoder := base64.StdEncoding
	// fmt.Println(enc)
	encodeByte := gcm.Seal(nonce, nonce, rawToken, nil)
	// fmt.Println(encByte)
	encodeToken := encoder.EncodeToString(encodeByte)
	return "token_" + encodeToken, nil
}
