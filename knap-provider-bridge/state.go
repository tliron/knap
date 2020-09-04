package main

import (
	"errors"
	"os"

	"github.com/gofrs/flock"
	"github.com/tliron/kutil/ard"
	"github.com/tliron/kutil/format"
)

const stateFilename = "/tmp/knap-provider-bridge.state"

var initialState ard.Map

func init() {
	availableBridgePrefixes := ard.List{
		"192.168.2",
		"192.168.3",
		"192.168.4",
		"192.168.5",
		"192.168.6",
		"192.168.7",
	}

	initialState = make(ard.Map)
	initialState["availableBridgePrefixes"] = availableBridgePrefixes
}

func GetState() (ard.Map, error) {
	if file, err := os.Open(stateFilename); err == nil {
		defer file.Close()
		if state, err := format.ReadYAML(file); err == nil {
			if map_, ok := state.(ard.Map); ok {
				return map_, nil
			} else {
				return nil, errors.New("state is not a map")
			}
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func SetState(state ard.Map) (ard.Map, error) {
	if file, err := os.Create(stateFilename); err == nil {
		defer file.Close()
		if err := format.WriteYAML(state, file, " ", false); err == nil {
			return state, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func GetBridgePrefix() (string, error) {
	// Lock on state file
	lock := flock.New(stateFilename)
	if err := lock.Lock(); err == nil {
		defer lock.Unlock()
	} else {
		return "", err
	}

	if state, err := GetState(); err == nil {
		if availableBridgePrefixes, ok := state["availableBridgePrefixes"]; ok {
			if availableBridgePrefixes_, ok := availableBridgePrefixes.(ard.List); ok {
				if len(availableBridgePrefixes_) > 0 {
					bridgePrefix := availableBridgePrefixes_[0]
					if bridgePrefix_, ok := bridgePrefix.(string); ok {
						availableBridgePrefixes_ = availableBridgePrefixes_[1:]
						state["availableBridgePrefixes"] = availableBridgePrefixes_
						if _, err := SetState(state); err == nil {
							return bridgePrefix_, nil
						} else {
							return "", err
						}
					}
				}
			}
		}
	} else {
		return "", err
	}

	return "", errors.New("no bridge prefixes available")
}
