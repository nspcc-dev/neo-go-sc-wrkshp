package domain

import (
	"github.com/nspcc-dev/neo-go/pkg/interop/runtime"
	"github.com/nspcc-dev/neo-go/pkg/interop/storage"
)

// Main is a very useful function.
func Main(operation string, args []interface{}) interface{} {
	ctx := storage.GetContext()

	// Queries the domain owner
	if operation == "query" {
		if len(args) != 1 {
			return false
		}

		domainName := args[0].([]byte)
		message := "QueryDomain: " + args[0].(string)
		runtime.Log(message)

		owner := storage.Get(ctx, domainName)
		if owner != nil {
			runtime.Log(owner.(string))
			return owner
		}

		runtime.Log("Domain is not yet registered")
	}

	// Deletes the domain
	if operation == "delete" {
		if len(args) != 1 {
			return false
		}
		domainName := args[0].([]byte)
		message := "DeleteDomain: " + args[0].(string)
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
		runtime.Notify([]interface{}{"deleted", owner, domainName})
		return true
	}

	// Registers new domain
	if operation == "register" {
		if len(args) != 2 {
			return false
		}
		domainName := args[0].([]byte)
		owner := args[1].([]byte)
		message := "RegisterDomain: " + args[0].(string)
		runtime.Log(message)

		if !runtime.CheckWitness(owner) {
			runtime.Log("Owner argument is not the same as the sender")
			return false
		}

		exists := storage.Get(ctx, domainName)
		if exists != nil {
			runtime.Log("Domain is already registered")
			return false
		}

		storage.Put(ctx, domainName, owner)
		runtime.Notify([]interface{}{"registered", owner, domainName})
		return true
	}

	// Transfers domain from one address to another
	if operation == "transfer" {
		if len(args) != 2 {
			return false
		}
		domainName := args[0].([]byte)
		message := "TransferDomain: " + args[0].(string)
		runtime.Log(message)

		owner := storage.Get(ctx, domainName)
		if owner == nil {
			runtime.Log("Domain is not yet registered")
			return false
		}

		if !runtime.CheckWitness(owner.([]byte)) {
			runtime.Log("Sender is not the owner, cannot transfer")
			return false
		}

		toAddress := args[1].([]byte)
		storage.Put(ctx, domainName, toAddress)
		runtime.Notify([]interface{}{"deleted", owner, domainName})
		runtime.Notify([]interface{}{"registered", toAddress, domainName})
		return true
	}

	return false
}
