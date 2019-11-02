package lp

import "sync"

const uuidLenght = 32

type uuid string

var uuids map[uuid]bool
var uuidLock sync.Mutex

func init() {
	uuids = make(map[uuid]bool)
}

func newUUID() uuid {
	uuidLock.Lock()
	defer uuidLock.Unlock()

	b := make([]byte, uuidLenght)
	var id uuid
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for true {
		for i := range b {
			b[i] = charset[seededRand.Intn(len(charset))]
		}
		id = uuid(b)
		if _, exists := uuids[id]; exists == false {
			uuids[id] = true
			break
		}
	}
	return id
}

func (u uuid) String() string {
	return string(u)
}
