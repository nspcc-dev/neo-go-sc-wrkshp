package domain

import (
	"github.com/nspcc-dev/neo-go/pkg/interop/runtime"
	"github.com/nspcc-dev/neo-go/pkg/interop/storage"
)

var ctx storage.Context

func init() {
	ctx = storage.GetContext()
}

// Query returns the owner of the domain with the specified name.
func Query(domainName []byte) interface{} {
	message := "QueryDomain: " + string(domainName)
	runtime.Log(message)

	owner := storage.Get(ctx, domainName)
	if owner != nil {
		runtime.Log("Domain is already registered")
		return owner
	}

	runtime.Log("Domain is not yet registered")
	return false
}

// Delete deletes domain with the specified name.
func Delete(domainName []byte) bool {
	message := "DeleteDomain: " + string(domainName)
	runtime.Log(message)

	owner := storage.Get(ctx, domainName)
	if owner == nil {
		runtime.Log("Domain is not yet registered")
		return false
	}

	if !runtime.CheckWitness(owner.([]byte)) {
		runtime.Log("Sender is not the owner, cannot delete")
		return false
	}

	storage.Delete(ctx, domainName)
	runtime.Notify("deleted", owner.([]byte), domainName)
	return true
}

// Register registers new domain with specified name and owner.
func Register(domainName []byte, owner []byte) bool {
	message := "RegisterDomain: " + string(domainName)
	runtime.Log(message)

	if !runtime.CheckWitness(owner) {
		runtime.Log("Owner argument is not the same as the sender")
		return false
	}

	if storage.Get(ctx, domainName) != nil {
		runtime.Log("Domain is already registered")
		return false
	}

	storage.Put(ctx, domainName, owner)
	runtime.Notify("registered", owner, domainName)
	return true
}

// Transfer transfers domain from owner to the specified address.
func Transfer(domainName []byte, toAddress []byte) bool {
	message := "TransferDomain: " + string(domainName)
	runtime.Log(message)

	owner := storage.Get(ctx, domainName)
	if owner == nil {
		runtime.Log("Domain is not yet registered")
		return false
	}

	if !runtime.CheckWitness(owner.([]byte)) {
		runtime.Log("Not co-signed by an owner, cannot transfer")
		return false
	}

	storage.Put(ctx, domainName, toAddress)
	runtime.Notify("deleted", owner.([]byte), domainName)
	runtime.Notify("registered", toAddress, domainName)
	return true
}
