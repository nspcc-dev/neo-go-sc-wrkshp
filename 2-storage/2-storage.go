package storage

import (
	"github.com/nspcc-dev/neo-go/pkg/interop/runtime"
	"github.com/nspcc-dev/neo-go/pkg/interop/storage"
)

const itemKey = "test-storage-key"

func _deploy(_ interface{}, isUpdate bool) {
	if !isUpdate {
		ctx := storage.GetContext()
		runtime.Notify("info", []byte("Storage key not yet set. Setting to 0"))
		itemValue := 0
		storage.Put(ctx, itemKey, itemValue)
		runtime.Notify("info", []byte("Storage key is initialised"))
	}
}

func Main() interface{} {
	ctx := storage.GetContext()
	itemValue := storage.Get(ctx, itemKey)
	runtime.Notify("info", []byte("Value read from storage"))

	runtime.Notify("info", []byte("Storage key already set. Incrementing by 1"))
	itemValue = itemValue.(int) + 1

	storage.Put(ctx, itemKey, itemValue)
	runtime.Notify("info", []byte("New value written into storage"))
	return itemValue
}
