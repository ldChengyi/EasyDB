package api

import "github.com/ldChengYi/EasyDB/core/storage"

func NewGenericStore[T any](opts storage.Options, indexSetup func(im *storage.IndexManager[T])) *storage.Store[T] {
	store := storage.New[T](opts)
	indexSetup(store.IndexManager)
	return store
}
