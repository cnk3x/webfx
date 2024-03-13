package eddsa

import (
	"crypto"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

var (
	ErrKeyMustBePEMEncoded = errors.New("invalid key: Key must be a PEM encoded PKCS1 or PKCS8 key")
	ErrNotEdPrivateKey     = errors.New("key is not a valid Ed25519 private key")
	ErrNotEdPublicKey      = errors.New("key is not a valid Ed25519 public key")
)

// 生成一对密钥对
func GenerateKey() (publicKey ed25519.PublicKey, privateKey ed25519.PrivateKey, err error) {
	return ed25519.GenerateKey(nil)
}

// PKCS8编码一个私钥
func EncodePrivate(privateKey ed25519.PrivateKey) (pemBytes []byte, err error) {
	keyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyBytes}), nil
}

// PKCS8解码一个私钥
func DecodePrivate(pemBytes []byte) (privateKey ed25519.PrivateKey, err error) {
	var block *pem.Block
	if block, _ = pem.Decode(pemBytes); block == nil {
		err = ErrKeyMustBePEMEncoded
		return
	}

	var key crypto.PrivateKey
	if key, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
		return
	}

	var ok bool
	if privateKey, ok = key.(ed25519.PrivateKey); !ok {
		err = ErrNotEdPrivateKey
	}

	return
}

// 加载一个私钥到文件
func ReadPrivateFile(privateKeyFile string) (privateKey ed25519.PrivateKey, err error) {
	var pemBytes []byte
	if pemBytes, err = os.ReadFile(privateKeyFile); err != nil {
		return
	}
	return DecodePrivate(pemBytes)
}

// 保存一个私钥到文件
func WritePrivateFile(privateKeyFile string, privateKey ed25519.PrivateKey) (err error) {
	var pemBytes []byte
	if pemBytes, err = EncodePrivate(privateKey); err != nil {
		return
	}
	return os.WriteFile(privateKeyFile, pemBytes, 0o644)
}

// PKIX编码一个公钥
func EncodePublic(publicKey ed25519.PublicKey) (pemBytes []byte, err error) {
	keyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: keyBytes}), nil
}

// PKIX解码一个公钥
func DecodePublic(pemBytes []byte) (publicKey ed25519.PublicKey, err error) {
	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(pemBytes); block == nil {
		err = ErrKeyMustBePEMEncoded
		return
	}

	var key crypto.PublicKey
	// Parse the key
	if key, err = x509.ParsePKIXPublicKey(block.Bytes); err != nil {
		return
	}

	var ok bool
	if publicKey, ok = key.(ed25519.PublicKey); !ok {
		err = ErrNotEdPublicKey
	}

	return
}

// 加载一个公钥到文件
func ReadPublicFile(publicKeyFile string) (publicKey ed25519.PublicKey, err error) {
	var pemBytes []byte
	if pemBytes, err = os.ReadFile(publicKeyFile); err != nil {
		return
	}
	return DecodePublic(pemBytes)
}

// 保存一个公钥到文件
func WritePublicFile(publicKeyFile string, publicKey ed25519.PublicKey) (err error) {
	var pemBytes []byte
	if pemBytes, err = EncodePublic(publicKey); err != nil {
		return
	}
	return os.WriteFile(publicKeyFile, pemBytes, 0o644)
}

func GeneratePublicKey(privateKey ed25519.PrivateKey) (publicKey ed25519.PublicKey) {
	return privateKey.Public().(ed25519.PublicKey)
}
