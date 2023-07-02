package memory_store

import (
	diskStore "github.com/goodstemy/redis-on-go/store/disk"
)

var _store map[string]map[string]interface{}

func Init() {
	_store = make(map[string]map[string]interface{})
}

func Hset(hash string, key string, value interface{}) {
	diskStore.Write()

	_, ok := _store[hash]

	if !ok {
		_store[hash] = make(map[string]interface{})
	}

	_store[hash][key] = value
}

func HGet(hash string, key string) interface{} {
	diskStore.Read()

	_, ok := _store[hash]

	if !ok {
		return nil
	}

	_, ok = _store[hash][key]

	if !ok {
		return nil
	}

	return _store[hash][key]
}
