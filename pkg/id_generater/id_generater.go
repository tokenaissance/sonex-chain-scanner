package id_generater

import (
	"fmt"
	"sync"
	"time"
)

const (
	epoch          = int64(1704135845000)
	machineIDBits  = uint(10)
	sequenceBits   = uint(12)
	maxMachineID   = int64(-1 ^ (-1 << machineIDBits))
	maxSequence    = int64(-1 ^ (-1 << sequenceBits))
	timestampShift = machineIDBits + sequenceBits
	machineIDShift = sequenceBits
)

type IdGenerater struct {
	mu        sync.Mutex
	lastTime  int64
	machineID int64
	sequence  int64
}

func NewIdGenerater(machineID int64) (*IdGenerater, error) {
	if machineID < 0 || machineID > maxMachineID {
		return nil, fmt.Errorf("%s", "machine id out of range")
	}

	return &IdGenerater{
		machineID: machineID,
	}, nil
}

func (sf *IdGenerater) Generate() (int64, error) {
	sf.mu.Lock()
	defer sf.mu.Unlock()

	now := time.Now().UnixMilli()
	if now < sf.lastTime {
		return 0, fmt.Errorf("%s", "clock moved backwards")
	}
	if now == sf.lastTime {
		sf.sequence = (sf.sequence + 1) & maxSequence
		if sf.sequence == 0 {
			for now <= sf.lastTime {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		sf.sequence = 0
	}
	sf.lastTime = now

	return ((now - epoch) << timestampShift) | (sf.machineID << machineIDShift) | sf.sequence, nil
}
