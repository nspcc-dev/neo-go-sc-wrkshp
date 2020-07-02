package domain

import (
	"github.com/nspcc-dev/neo-go/pkg/interop/runtime"
	"github.com/nspcc-dev/neo-go/pkg/interop/storage"
)

// Main is a very useful function.
func Main(operation string, args []interface{}) interface{} {
	// Queries the domain owner
	if operation == "query" {
		if len(args) != 1 {
			return false
		}

		message := "QueryDomain: " + args[0].(string)
		runtime.Log(message)

		return Query(args[0].([]byte))
	}

	// Deletes the domain
	if operation == "delete" {
		if len(args) != 1 {
			return false
		}

		message := "DeleteDomain: " + args[0].(string)
		runtime.Log(message)

		return Delete(args[0].([]byte))
	}

	// Registers new domain
	if operation == "register" {
		if len(args) != 2 {
			return false
		}

		message := "RegisterDomain: " + args[0].(string)
		runtime.Log(message)

		return Register(args[0].([]byte), args[1].([]byte))
	}

	// Transfers domain from one address to another
	if operation == "transfer" {
		if len(args) != 2 {
			return false
		}

		message := "TransferDomain: " + args[0].(string)
		runtime.Log(message)

		return Transfer(args[0].([]byte), args[1].([]byte))
	}

	return false
}

// Query returns the owner of the domain with the specified name.
func Query(domainName []byte) interface{} {
	ctx := storage.GetContext()
	owner := storage.Get(ctx, domainName)
	if owner != nil {
		runtime.Log(owner.(string))
		return owner
	}

	runtime.Log("Domain is not yet registered")
	return false
}

// Delete deletes domain with the specified name.
func Delete(domainName []byte) bool {
	ctx := storage.GetContext()
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
	runtime.Notify("deleted", owner, domainName)
	return true
}

// Register registers new domain with specified name and owner.
func Register(domainName []byte, owner []byte) bool {
	ctx := storage.GetContext()

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
	ctx := storage.GetContext()
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
	runtime.Notify("deleted", owner, domainName)
	runtime.Notify("registered", toAddress, domainName)
	return true
}
