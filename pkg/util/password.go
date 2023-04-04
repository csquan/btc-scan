package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"github.com/btc-scan/pkg/conf"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"unicode"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/nacl/secretbox"
)

func VerifyRule(s string) (eightOrMore, number, upper, special bool) {
	letters := 0
	for _, c := range s {
		letters++
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		case unicode.IsLetter(c) || c == ' ':
		default:
			//return false, false, false, false
		}
	}
	eightOrMore = letters >= 8
	return
}

func HashPassword(plain string) (string, Err) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)

	if err != nil {
		return "", ErrBcryptHash
	}
	return string(hashed), nil
}

func VerifyPassword(hashed string, plain string) Err {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)); err != nil {
		return ErrBcryptComp
	}
	return nil
}

func CreateSha2(key string) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(key))
	return hasher.Sum(nil)
}

func AesGcmEnc(key, data []byte) ([]byte, Err) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, ErrCryptoAesCipher
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, ErrCryptoAesGcm
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, ErrCryptoRand
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func AesGcmDec(key, data []byte) ([]byte, Err) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, ErrDeCryptoAesCipher
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, ErrDeCryptoAesGcm
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, ErrDeCryptoAesDec
	}
	return plaintext, nil
}

func SecretBoxEnc(data []byte, passphrase string) []byte {
	var secretKey [32]byte
	copy(secretKey[:], CreateSha2(passphrase))

	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		panic(err)
	}
	encrypted := secretbox.Seal(nonce[:], data, &nonce, &secretKey)

	return encrypted
}

func SecretBoxDec(encrypted []byte, passphrase string) []byte {
	var secretKey [32]byte
	copy(secretKey[:], CreateSha2(passphrase))

	var decryptNonce [24]byte
	copy(decryptNonce[:], encrypted[:24])
	decrypted, ok := secretbox.Open(nil, encrypted[24:], &decryptNonce, &secretKey)
	if !ok {
		panic("decryption error")
	}

	return decrypted
}

// AesEncrypt encrypts text and given key with AES.
func AesEncrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

// AesDecrypt decrypts text and given key with AES.
func AesDecrypt(key, text []byte) ([]byte, Err) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, ErrCryptoAesCipher
	}
	if len(text) < aes.BlockSize {
		return nil, ErrAesSize
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, ErrMsgDecode
	}
	return data, nil
}

// EncryptSecret encrypts a string with given key into a hex string
func EncryptSecret(key, str string) (string, error) {
	keyHash := sha256.Sum256([]byte(key))
	plaintext := []byte(str)
	ciphertext, err := AesEncrypt(keyHash[:], plaintext)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(ciphertext), nil
}

// DecryptSecret decrypts a previously encrypted hex string
func DecryptSecret(key, cipherHex string) (string, Err) {
	keyHash := sha256.Sum256([]byte(key))
	ciphertext, err := hex.DecodeString(cipherHex)
	if err != nil {
		return "", ErrMsgDecode
	}
	plaintext, er := AesDecrypt(keyHash[:], ciphertext)
	if er != nil {
		return "", er
	}
	return string(plaintext), nil
}

// HashToken return the hashable salt
func HashToken(token, salt string) string {
	tempHash := pbkdf2.Key([]byte(token), []byte(salt), 10000, 50, sha256.New)
	return hex.EncodeToString(tempHash)
}

func GetEncryptionKey() []byte {
	k := CreateSha2(conf.Conf.AesSalt)
	return k[:]
}

// HashSecret 谷歌验证的密码加密存储
func HashSecret(plain string) (string, Err) {
	secretBytes, err := AesGcmEnc(GetEncryptionKey(), []byte(plain))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(secretBytes), nil
}

func DecodeSecret(hashed string) (string, Err) {
	decodedStoredSecret, err := base64.StdEncoding.DecodeString(hashed)
	if err != nil {
		return "", ErrMsgDecode
	}
	secretBytes, er := AesGcmDec(GetEncryptionKey(), decodedStoredSecret)
	if er != nil {
		return "", er
	}
	return string(secretBytes), nil
}

// ValidateTOTP validates the provided passcode.
func ValidateTOTP(passcode, hashedSecret string) (bool, Err) {
	if hashedSecret == "" {
		return false, ErrGaNew
	}
	secretStr, er := DecodeSecret(hashedSecret)
	if er != nil {
		return false, er
	}
	return totp.Validate(passcode, secretStr), nil
}
