package quartzflake

import (
	"errors"
	"sync"
	"time"
)

const (
	TimestampBits = 63
	MachineIDBits = 14
	SequenceBits  = 18

	MaxTimestamp = (1 << TimestampBits) - 1
	MaxMachineID = (1 << MachineIDBits) - 1
	MaxSequence  = (1 << SequenceBits) - 1
)

type Quartzflake struct {
	mutex *sync.Mutex

	timestamp int64
	machineid int
	sequence  int
}

type QuartzflakeID struct {
	timestamp int64 // Timestamp in milliseconds since unix epoch
	metadata  int   // Machine ID and Sequence
}

var (
	ErrInvalidTimestamp     = errors.New("timpestamp out of range")
	ErrInvalidMachineIDBits = errors.New("invalid bit length for machine id")
	ErrInvalidSequenceBits  = errors.New("invalid bit length for sequence number")
)

func (mf *Quartzflake) NextID() (QuartzflakeID, error) {
	mf.mutex.Lock()
	defer mf.mutex.Unlock()

	now := time.Now().UnixMilli()
	if mf.timestamp != now {
		mf.timestamp = now
		mf.sequence = 0
	} else {
		mf.sequence++
		if mf.sequence > MaxSequence {
			// Wait for the next millisecond
			for now <= mf.timestamp {
				time.Sleep(time.Millisecond)
				now = time.Now().UnixMilli()
			}
			mf.timestamp = now
			mf.sequence = 0
		}
	}

	return mf.toID()
}

func (mf *Quartzflake) toID() (QuartzflakeID, error) {
	if mf.timestamp < 0 || mf.timestamp > MaxTimestamp {
		return QuartzflakeID{}, ErrInvalidTimestamp
	}

	return QuartzflakeID{
		timestamp: mf.timestamp,
		metadata:  (mf.machineid << SequenceBits) | mf.sequence,
	}, nil
}

func (id QuartzflakeID) Timestamp() int64 {
	return id.timestamp
}

func (id QuartzflakeID) MachineID() int {
	return id.metadata >> SequenceBits
}

func (id QuartzflakeID) Sequence() int {
	return id.metadata & MaxSequence
}
