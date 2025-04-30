package transaction

import (
	"context"
	"errors"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// TransactionManager handles secure transaction operations
type TransactionManager struct {
	client    *ethclient.Client
	gasLimits map[string]*big.Int
	mutex     sync.RWMutex
}

// NewTransactionManager creates a new transaction manager instance
func NewTransactionManager(client *ethclient.Client) *TransactionManager {
	return &TransactionManager{
		client:    client,
		gasLimits: make(map[string]*big.Int),
		mutex:     sync.RWMutex{},
	}
}

// SetGasLimit sets a gas limit for a specific contract address
func (tm *TransactionManager) SetGasLimit(contractAddr string, limit *big.Int) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	tm.gasLimits[contractAddr] = limit
}

// GetGasLimit retrieves the gas limit for a specific contract address
func (tm *TransactionManager) GetGasLimit(contractAddr string) *big.Int {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()
	if limit, exists := tm.gasLimits[contractAddr]; exists {
		return limit
	}
	return big.NewInt(0)
}

// VerifyTransaction performs security checks on a transaction
func (tm *TransactionManager) VerifyTransaction(ctx context.Context, tx *types.Transaction) error {
	// Verify gas limit
	if tx.Gas() == 0 {
		return errors.New("gas limit cannot be zero")
	}

	// Check against contract-specific gas limits
	if tx.To() != nil {
		contractAddr := tx.To().Hex()
		limit := tm.GetGasLimit(contractAddr)
		if limit.Cmp(big.NewInt(0)) > 0 && tx.Gas() > limit.Uint64() {
			return errors.New("transaction exceeds gas limit for contract")
		}
	}

	// Verify gas price is reasonable
	suggestedGasPrice, err := tm.client.SuggestGasPrice(ctx)
	if err != nil {
		return err
	}

	// Allow up to 50% more than suggested gas price
	maxGasPrice := new(big.Int).Mul(suggestedGasPrice, big.NewInt(150))
	maxGasPrice = maxGasPrice.Div(maxGasPrice, big.NewInt(100))

	if tx.GasPrice().Cmp(maxGasPrice) > 0 {
		return errors.New("gas price too high")
	}

	return nil
}

// EstimateGas estimates gas needed for a transaction with safety margin
func (tm *TransactionManager) EstimateGas(ctx context.Context, from common.Address, to *common.Address, data []byte) (uint64, error) {
	// Create a message call transaction
	msg := ethereum.CallMsg{
		From: from,
		To:   to,
		Data: data,
	}

	// Estimate gas
	gas, err := tm.client.EstimateGas(ctx, msg)
	if err != nil {
		return 0, err
	}

	// Add 10% safety margin
	gas = gas + (gas / 10)

	return gas, nil
}

// WaitForTransaction waits for transaction confirmation with timeout
func (tm *TransactionManager) WaitForTransaction(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	retry := 0
	maxRetries := 30 // Adjust based on block time and desired timeout

	for retry < maxRetries {
		receipt, err := tm.client.TransactionReceipt(ctx, txHash)
		if err == nil && receipt != nil {
			// Verify transaction success
			if receipt.Status == types.ReceiptStatusSuccessful {
				return receipt, nil
			}
			return receipt, errors.New("transaction failed")
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(2 * time.Second):
			retry++
			continue
		}
	}

	return nil, errors.New("transaction confirmation timeout")
}
