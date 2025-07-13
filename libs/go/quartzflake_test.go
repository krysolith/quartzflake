package quartzflake

import (
	"sync"
	"testing"
	"time"
)

func newTestQuartzflake(machine int) *Quartzflake {
	return &Quartzflake{
		mutex:     &sync.Mutex{},
		machineid: machine,
		sequence:  0,
	}
}

func TestBasicIDGeneration(t *testing.T) {
	mf := newTestQuartzflake(1)
	id, err := mf.NextID()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id.timestamp == 0 {
		t.Error("timestamp should not be zero")
	}
	if id.metadata>>SequenceBits != 1 {
		t.Errorf("machine ID not encoded correctly: got %d", id.metadata>>SequenceBits)
	}
}

func TestSequenceIncrement(t *testing.T) {
	mf := newTestQuartzflake(2)
	mf.timestamp = time.Now().UnixMilli()
	id1, _ := mf.NextID()
	id2, _ := mf.NextID()
	if (id2.metadata & MaxSequence) != ((id1.metadata & MaxSequence) + 1) {
		t.Error("sequence did not increment correctly")
	}
}

func TestSequenceOverflow(t *testing.T) {
	mf := newTestQuartzflake(3)
	mf.timestamp = time.Now().UnixMilli()
	mf.sequence = MaxSequence
	id, _ := mf.NextID()
	if id.metadata&MaxSequence != 0 {
		t.Error("sequence should reset to 0 after overflow")
	}
}

func TestUniqueIDs(t *testing.T) {
	mf := newTestQuartzflake(4)
	ids := make(map[int64]struct{})
	for i := 0; i < 1000; i++ {
		id, _ := mf.NextID()
		key := (id.timestamp << 16) | int64(id.metadata)
		if _, exists := ids[key]; exists {
			t.Fatalf("duplicate ID detected: %v", id)
		}
		ids[key] = struct{}{}
	}
}
