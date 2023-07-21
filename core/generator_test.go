package core

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockStatusLoader struct {
	err error
}

func (m *mockStatusLoader) AllStatus() ([]*Status, error) {
	if m.err != nil {
		return nil, m.err
	}
	return []*Status{{LastID: 100}}, nil
}

func (m *mockStatusLoader) Status(name string) (*Status, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &Status{LastID: 100}, nil
}

func (m *mockStatusLoader) SaveStatus(name string, s *Status) error {
	if m.err != nil {
		return m.err
	}
	return nil
}

func TestNewGenerator(t *testing.T) {
	// Positive test case
	sl := &mockStatusLoader{}
	g := NewGenerator("test", sl)
	assert.NotNil(t, g)

	// Negative test case
	sl = &mockStatusLoader{err: errors.New("status load error")}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("NewGenerator did not panic")
		}
	}()
	NewGenerator("test", sl)
}

func TestIDGenerator_NextID(t *testing.T) {
	sl := &mockStatusLoader{}
	g := &IDGenerator{
		cfg:    &Config{LoadNum: 5, LoadThreshold: 2},
		status: &Status{LastID: 100},
		ch:     make(chan struct{}, 1),
		num:    0,
		lastID: 100,
		loader: sl,
	}

	for i := 0; i < 6; i++ {
		id := g.NextID()
		assert.Equal(t, id, int64(i+101))
		if (g.cfg.LoadNum - int32(i+1)) == g.cfg.LoadThreshold {
			select {
			case <-g.ch:
			default:
				t.Errorf("NextID did not send signal to g.ch")
			}
		}
		if i+1 == int(g.cfg.LoadNum) {
			assert.Equal(t, g.num, int32(0))
		}
	}
}

func TestIDGenerator_loadAndUpdateStatus(t *testing.T) {
	sl := &mockStatusLoader{}
	g := &IDGenerator{
		cfg:    &Config{LoadNum: 5, IncrFirstID: 5},
		loader: sl,
	}
	// Positive test case
	err := g.loadAndUpdateStatus()
	if err != nil {
		t.Errorf("loadAndUpdateStatus returned error: %s", err)
	}
	assert.Equal(t, int64(10), g.status.LastID)

	// Negative test case
	sl.err = errors.New("status update error")
	err = g.loadAndUpdateStatus()
	assert.EqualError(t, err, "status update error")
}

func TestIDGenerator_watchLoadThreshold(t *testing.T) {
	sl := &mockStatusLoader{}
	g := &IDGenerator{
		cfg:                &Config{LoadNum: 5, LoadThreshold: 2},
		status:             &Status{LastID: 100},
		watchRetryInterval: 10 * time.Millisecond,
		ch:                 make(chan struct{}),
		loader:             sl,
	}
	// Positive test case
	go g.watchLoadThreshold()
	g.ch <- struct{}{}
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, g.status.LastID, int64(105))

	// Negative test case
	sl.err = errors.New("status update error")
	g.ch <- struct{}{}
	time.Sleep(20 * time.Millisecond)
	sl.err = nil
	time.Sleep(20 * time.Millisecond)
	assert.Equal(t, g.status.LastID, int64(110))
}
