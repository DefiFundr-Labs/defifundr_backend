package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// WalletManager handles secure wallet operations
type WalletManager struct {
	keystore *keystore.KeyStore
	keysDir  string
}

// NewWalletManager creates a new wallet manager instance
func NewWalletManager(keysDir string) (*WalletManager, error) {
	if err := os.MkdirAll(keysDir, 0700); err != nil {
		return nil, err
	}

	ks := keystore.NewKeyStore(
		keysDir,
		keystore.StandardScryptN,
		keystore.StandardScryptP,
	)

	return &WalletManager{
		keystore: ks,
		keysDir:  keysDir,
	}, nil
}

// CreateWallet generates a new wallet with encrypted private key
func (wm *WalletManager) CreateWallet(password string) (accounts.Account, error) {
	if len(password) < 8 {
		return accounts.Account{}, errors.New("password must be at least 8 characters")
	}

	account, err := wm.keystore.NewAccount(password)
	if err != nil {
		return accounts.Account{}, err
	}

	return account, nil
}

// ImportPrivateKey imports an existing private key
func (wm *WalletManager) ImportPrivateKey(privateKeyHex string, password string) (accounts.Account, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return accounts.Account{}, err
	}

	account, err := wm.keystore.ImportECDSA(privateKey, password)
	if err != nil {
		return accounts.Account{}, err
	}

	return account, nil
}

// SignTransaction signs a transaction with the account's private key
func (wm *WalletManager) SignTransaction(account accounts.Account, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	signedTx, err := wm.keystore.SignTx(account, tx, chainID)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

// EncryptData encrypts sensitive data using AES-GCM
func (wm *WalletManager) EncryptData(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

// DecryptData decrypts AES-GCM encrypted data
func (wm *WalletManager) DecryptData(encryptedData []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, errors.New("encrypted data too short")
	}

	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
