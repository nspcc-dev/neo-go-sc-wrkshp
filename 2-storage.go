package storage

import (
	"github.com/nspcc-dev/neo-go/pkg/interop/runtime"
	"github.com/nspcc-dev/neo-go/pkg/interop/storage"
)

func Main() interface{} {
	ctx := storage.GetContext()
	itemKey := "test-storage-key"
	itemValue := storage.Get(ctx, itemKey)
	msg := "Value read from storage"

	runtime.Notify(msg)

	if itemValue == nil {
		runtime.Notify("Storage key not yet set. Setting to 1")
		itemValue = 1
	} else {
		runtime.Notify("Storage key already set. Incrementing by 1")
		itemValue = itemValue.(int) + 1
	}

	storage.Put(ctx, itemKey, itemValue)
	runtime.Notify("New value written into storage")
	return itemValue
}
