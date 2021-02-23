package server

import (
	"errors"
	"os"

	"github.com/gofrs/flock"
	"github.com/tliron/kutil/ard"
	"github.com/tliron/kutil/format"
)

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

func (self *Server) InitializeState() error {
	_, err := self.SetState(initialState)
	return err
}

func (self *Server) GetState() (ard.Map, error) {
	if file, err := os.Open(self.StateFilename); err == nil {
		defer file.Close()
		if state, _, err := ard.ReadYAML(file, false); err == nil {
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

func (self *Server) SetState(state ard.Map) (ard.Map, error) {
	if file, err := os.Create(self.StateFilename); err == nil {
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

func (self *Server) ProvideBridgePrefix() (string, error) {
	// Lock on state file
	lock := flock.New(self.StateFilename)
	if err := lock.Lock(); err == nil {
		defer lock.Unlock()
	} else {
		return "", err
	}

	if state, err := self.GetState(); err == nil {
		if availableBridgePrefixes, ok := state["availableBridgePrefixes"]; ok {
			if availableBridgePrefixes_, ok := availableBridgePrefixes.(ard.List); ok {
				if len(availableBridgePrefixes_) > 0 {
					bridgePrefix := availableBridgePrefixes_[0]
					if bridgePrefix_, ok := bridgePrefix.(string); ok {
						availableBridgePrefixes_ = availableBridgePrefixes_[1:]
						state["availableBridgePrefixes"] = availableBridgePrefixes_
						if _, err := self.SetState(state); err == nil {
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
