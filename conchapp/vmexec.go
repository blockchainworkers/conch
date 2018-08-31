package conchapp

import (
	"math/big"
)

// Exec transcton update state

// VMActuator vitual machine to exec transaction then update related state
type VMActuator struct {
	appSt *APPState
}

// NewVMActuator return VMActuator inst
func NewVMActuator(state *APPState) *VMActuator {
	return &VMActuator{appSt: state}
}

// ExecuteTx exec tx and update state
func (vm *VMActuator) ExecuteTx(tx *Transaction) error {

	// fee
	fee := tx.FeeCalc()
	vm.appSt.HeadSt.Fee.Add(vm.appSt.HeadSt.Fee, fee)
	value, _ := big.NewInt(0).SetString(tx.Value, 0)

	txrp := &TransactionReceipt{Status: 0, Fee: fee, BlockNum: vm.appSt.HeadSt.CurBlockNum, TxHash: tx.TxID()}

	sender, err := vm.appSt.AccoutSt.LoadAccount(tx.Sender)
	if err != nil {
		txrp.Log = "sender account err"
		vm.appSt.TxRepSt.UpdateTxRep(txrp)
		return err
	}

	receiver, err := vm.appSt.AccoutSt.LoadAccount(tx.Receiver)
	if err != nil {
		txrp.Log = "receiver account load err"
		vm.appSt.TxRepSt.UpdateTxRep(txrp)
		return err
	}

	if receiver.Address == "" {
		receiver.Address = "00000000000000000000"
	}

	// exec success
	txrp.Status = 1
	sender.Amount.Sub(sender.Amount, big.NewInt(0).Add(value, fee))
	receiver.Amount.Add(receiver.Amount, value)
	vm.appSt.AccoutSt.UpdateAccountCache(sender)
	vm.appSt.AccoutSt.UpdateAccountCache(receiver)
	vm.appSt.TxRepSt.UpdateTxRep(txrp)
	return nil
}
