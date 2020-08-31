package main

import (
	"errors"
	"sync"
)

var bridgePrefixPool = []string{
	"192.168.2",
	"192.168.3",
	"192.168.4",
	"192.168.5",
	"192.168.6",
	"192.168.7",
}

var bridgePrefixLock sync.Mutex

func getBridgePrefix() (string, error) {
	bridgePrefixLock.Lock()
	defer bridgePrefixLock.Unlock()

	if len(bridgePrefixPool) == 0 {
		return "", errors.New("bridge prefix pool is empty")
	}

	bridgePrefix := bridgePrefixPool[0]
	bridgePrefixPool = bridgePrefixPool[1:]
	return bridgePrefix, nil
}
