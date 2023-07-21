// Package storagecontract contains RPC wrappers for Storage example contract.
package storagecontract

import (
	"github.com/nspcc-dev/neo-go/pkg/core/transaction"
	"github.com/nspcc-dev/neo-go/pkg/smartcontract"
	"github.com/nspcc-dev/neo-go/pkg/util"
)

// Actor is used by Contract to call state-changing methods.
type Actor interface {
	MakeCall(contract util.Uint160, method string, params ...any) (*transaction.Transaction, error)
	MakeRun(script []byte) (*transaction.Transaction, error)
	MakeUnsignedCall(contract util.Uint160, method string, attrs []transaction.Attribute, params ...any) (*transaction.Transaction, error)
	MakeUnsignedRun(script []byte, attrs []transaction.Attribute) (*transaction.Transaction, error)
	SendCall(contract util.Uint160, method string, params ...any) (util.Uint256, uint32, error)
	SendRun(script []byte) (util.Uint256, uint32, error)
}

// Contract implements all contract methods.
type Contract struct {
	actor Actor
	hash util.Uint160
}

// New creates an instance of Contract using provided contract hash and the given Actor.
func New(actor Actor, hash util.Uint160) *Contract {
	return &Contract{actor, hash}
}

func (c *Contract) scriptForDelete(key []byte) ([]byte, error) {
	return smartcontract.CreateCallWithAssertScript(c.hash, "delete", key)
}

// Delete creates a transaction invoking `delete` method of the contract.
// This transaction is signed and immediately sent to the network.
// The values returned are its hash, ValidUntilBlock value and error if any.
func (c *Contract) Delete(key []byte) (util.Uint256, uint32, error) {
	script, err := c.scriptForDelete(key)
	if err != nil {
		return util.Uint256{}, 0, err
	}
	return c.actor.SendRun(script)
}

// DeleteTransaction creates a transaction invoking `delete` method of the contract.
// This transaction is signed, but not sent to the network, instead it's
// returned to the caller.
func (c *Contract) DeleteTransaction(key []byte) (*transaction.Transaction, error) {
	script, err := c.scriptForDelete(key)
	if err != nil {
		return nil, err
	}
	return c.actor.MakeRun(script)
}

// DeleteUnsigned creates a transaction invoking `delete` method of the contract.
// This transaction is not signed, it's simply returned to the caller.
// Any fields of it that do not affect fees can be changed (ValidUntilBlock,
// Nonce), fee values (NetworkFee, SystemFee) can be increased as well.
func (c *Contract) DeleteUnsigned(key []byte) (*transaction.Transaction, error) {
	script, err := c.scriptForDelete(key)
	if err != nil {
		return nil, err
	}
	return c.actor.MakeUnsignedRun(script, nil)
}

// Find creates a transaction invoking `find` method of the contract.
// This transaction is signed and immediately sent to the network.
// The values returned are its hash, ValidUntilBlock value and error if any.
func (c *Contract) Find(value []byte) (util.Uint256, uint32, error) {
	return c.actor.SendCall(c.hash, "find", value)
}

// FindTransaction creates a transaction invoking `find` method of the contract.
// This transaction is signed, but not sent to the network, instead it's
// returned to the caller.
func (c *Contract) FindTransaction(value []byte) (*transaction.Transaction, error) {
	return c.actor.MakeCall(c.hash, "find", value)
}

// FindUnsigned creates a transaction invoking `find` method of the contract.
// This transaction is not signed, it's simply returned to the caller.
// Any fields of it that do not affect fees can be changed (ValidUntilBlock,
// Nonce), fee values (NetworkFee, SystemFee) can be increased as well.
func (c *Contract) FindUnsigned(value []byte) (*transaction.Transaction, error) {
	return c.actor.MakeUnsignedCall(c.hash, "find", nil, value)
}

// FindReturnIter creates a transaction invoking `findReturnIter` method of the contract.
// This transaction is signed and immediately sent to the network.
// The values returned are its hash, ValidUntilBlock value and error if any.
func (c *Contract) FindReturnIter(prefix []byte) (util.Uint256, uint32, error) {
	return c.actor.SendCall(c.hash, "findReturnIter", prefix)
}

// FindReturnIterTransaction creates a transaction invoking `findReturnIter` method of the contract.
// This transaction is signed, but not sent to the network, instead it's
// returned to the caller.
func (c *Contract) FindReturnIterTransaction(prefix []byte) (*transaction.Transaction, error) {
	return c.actor.MakeCall(c.hash, "findReturnIter", prefix)
}

// FindReturnIterUnsigned creates a transaction invoking `findReturnIter` method of the contract.
// This transaction is not signed, it's simply returned to the caller.
// Any fields of it that do not affect fees can be changed (ValidUntilBlock,
// Nonce), fee values (NetworkFee, SystemFee) can be increased as well.
func (c *Contract) FindReturnIterUnsigned(prefix []byte) (*transaction.Transaction, error) {
	return c.actor.MakeUnsignedCall(c.hash, "findReturnIter", nil, prefix)
}

// Get creates a transaction invoking `get` method of the contract.
// This transaction is signed and immediately sent to the network.
// The values returned are its hash, ValidUntilBlock value and error if any.
func (c *Contract) Get(key []byte) (util.Uint256, uint32, error) {
	return c.actor.SendCall(c.hash, "get", key)
}

// GetTransaction creates a transaction invoking `get` method of the contract.
// This transaction is signed, but not sent to the network, instead it's
// returned to the caller.
func (c *Contract) GetTransaction(key []byte) (*transaction.Transaction, error) {
	return c.actor.MakeCall(c.hash, "get", key)
}

// GetUnsigned creates a transaction invoking `get` method of the contract.
// This transaction is not signed, it's simply returned to the caller.
// Any fields of it that do not affect fees can be changed (ValidUntilBlock,
// Nonce), fee values (NetworkFee, SystemFee) can be increased as well.
func (c *Contract) GetUnsigned(key []byte) (*transaction.Transaction, error) {
	return c.actor.MakeUnsignedCall(c.hash, "get", nil, key)
}

// Get_0 creates a transaction invoking `get` method of the contract.
// This transaction is signed and immediately sent to the network.
// The values returned are its hash, ValidUntilBlock value and error if any.
func (c *Contract) Get_0() (util.Uint256, uint32, error) {
	return c.actor.SendCall(c.hash, "get")
}

// Get_0Transaction creates a transaction invoking `get` method of the contract.
// This transaction is signed, but not sent to the network, instead it's
// returned to the caller.
func (c *Contract) Get_0Transaction() (*transaction.Transaction, error) {
	return c.actor.MakeCall(c.hash, "get")
}

// Get_0Unsigned creates a transaction invoking `get` method of the contract.
// This transaction is not signed, it's simply returned to the caller.
// Any fields of it that do not affect fees can be changed (ValidUntilBlock,
// Nonce), fee values (NetworkFee, SystemFee) can be increased as well.
func (c *Contract) Get_0Unsigned() (*transaction.Transaction, error) {
	return c.actor.MakeUnsignedCall(c.hash, "get", nil)
}

// Put creates a transaction invoking `put` method of the contract.
// This transaction is signed and immediately sent to the network.
// The values returned are its hash, ValidUntilBlock value and error if any.
func (c *Contract) Put(key []byte, value []byte) (util.Uint256, uint32, error) {
	return c.actor.SendCall(c.hash, "put", key, value)
}

// PutTransaction creates a transaction invoking `put` method of the contract.
// This transaction is signed, but not sent to the network, instead it's
// returned to the caller.
func (c *Contract) PutTransaction(key []byte, value []byte) (*transaction.Transaction, error) {
	return c.actor.MakeCall(c.hash, "put", key, value)
}

// PutUnsigned creates a transaction invoking `put` method of the contract.
// This transaction is not signed, it's simply returned to the caller.
// Any fields of it that do not affect fees can be changed (ValidUntilBlock,
// Nonce), fee values (NetworkFee, SystemFee) can be increased as well.
func (c *Contract) PutUnsigned(key []byte, value []byte) (*transaction.Transaction, error) {
	return c.actor.MakeUnsignedCall(c.hash, "put", nil, key, value)
}

// Put_1 creates a transaction invoking `put` method of the contract.
// This transaction is signed and immediately sent to the network.
// The values returned are its hash, ValidUntilBlock value and error if any.
func (c *Contract) Put_1(value []byte) (util.Uint256, uint32, error) {
	return c.actor.SendCall(c.hash, "put", value)
}

// Put_1Transaction creates a transaction invoking `put` method of the contract.
// This transaction is signed, but not sent to the network, instead it's
// returned to the caller.
func (c *Contract) Put_1Transaction(value []byte) (*transaction.Transaction, error) {
	return c.actor.MakeCall(c.hash, "put", value)
}

// Put_1Unsigned creates a transaction invoking `put` method of the contract.
// This transaction is not signed, it's simply returned to the caller.
// Any fields of it that do not affect fees can be changed (ValidUntilBlock,
// Nonce), fee values (NetworkFee, SystemFee) can be increased as well.
func (c *Contract) Put_1Unsigned(value []byte) (*transaction.Transaction, error) {
	return c.actor.MakeUnsignedCall(c.hash, "put", nil, value)
}
