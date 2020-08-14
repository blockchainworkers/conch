package cvm

import (
	"github.com/blockchainworkers/conch/cvm/common"
	"github.com/blockchainworkers/conch/cvm/vm"
	"math/big"
)

// NewEVMContext creates a new context for use in the EVM.
func NewEVMContext(from common.Address, blockNum, timeStamp, difficulty int64) vm.Context {
	// If we don't have an explicit author (i.e. not mining), extract from the header
	return vm.Context{
		CanTransfer: CanTransfer,
		Transfer:    Transfer,
		GetHash:     GetHashFn(),
		Origin:      from,
		Coinbase:    common.Address{},
		BlockNumber: new(big.Int).Set(big.NewInt(blockNum)),
		Time:        new(big.Int).Set(big.NewInt(timeStamp)),
		Difficulty:  new(big.Int).Set(big.NewInt(difficulty)),
		GasLimit:    0xfffffffffffffff, //header.GasLimit,
		GasPrice:    new(big.Int).Set(big.NewInt(10)),
	}
}

// GetHashFn returns a GetHashFunc which retrieves header hashes by number 获取块号码对于的块hash
func GetHashFn() func(n uint64) common.Hash {

	return func(n uint64) common.Hash {
		// If there's no hash cache yet, make one
		// if cache == nil {
		// 	cache = map[uint64]common.Hash{
		// 		ref.Number.Uint64() - 1: ref.ParentHash,
		// 	}
		// }
		// // Try to fulfill the request from the cache
		// if hash, ok := cache[n]; ok {
		// 	return hash
		// }
		// // Not cached, iterate the blocks and cache the hashes
		// for header := chain.GetHeader(ref.ParentHash, ref.Number.Uint64()-1); header != nil; header = chain.GetHeader(header.ParentHash, header.Number.Uint64()-1) {
		// 	cache[header.Number.Uint64()-1] = header.ParentHash
		// 	if n == header.Number.Uint64()-1 {
		// 		return header.ParentHash
		// 	}
		// }
		return common.Hash{}
	}
}

// CanTransfer checks wether there are enough funds in the address' account to make a transfer.
// This does not take the necessary gas in to account to make the transfer valid.
func CanTransfer(db vm.StateDB, addr common.Address, amount *big.Int) bool {
	return db.GetBalance(addr).Cmp(amount) >= 0
}

// Transfer subtracts amount from sender and adds amount to recipient using the given Db
func Transfer(db vm.StateDB, sender, recipient common.Address, amount *big.Int) {
	db.SubBalance(sender, amount)
	db.AddBalance(recipient, amount)
}
