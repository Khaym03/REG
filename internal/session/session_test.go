package session

import (
	"errors"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsolatedRodSessionCloseRunsCleanupOnce(t *testing.T) {
	var calls atomic.Int32
	s := NewIsolatedRodSession(nil, nil, func() error {
		calls.Add(1)
		return nil
	})

	require.NoError(t, s.Close())
	require.NoError(t, s.Close())
	assert.Equal(t, int32(1), calls.Load())
}

func TestIsolatedRodSessionClosePropagatesCleanupError(t *testing.T) {
	expectedErr := errors.New("cleanup failed")
	s := NewIsolatedRodSession(nil, nil, func() error {
		return expectedErr
	})

	err := s.Close()
	assert.ErrorIs(t, err, expectedErr)
}

func TestRodSessionIDIsStableAndUnique(t *testing.T) {
	s1 := NewRodSession(nil, nil)
	s2 := NewIsolatedRodSession(nil, nil, nil)

	assert.NotEmpty(t, s1.ID())
	assert.Equal(t, s1.ID(), s1.ID(), "ID() must return the same value across calls")
	assert.Equal(t, s2.ID(), s2.ID(), "isolated session ID() must return the same value across calls")
	assert.NotEqual(t, s1.ID(), s2.ID(), "distinct sessions must have distinct IDs")
}

func TestSessionRegistryLookupWithStableRodSessionID(t *testing.T) {
	reg := NewSessionRegistry()
	s := NewRodSession(nil, nil)

	reg.Add(s)

	got, ok := reg.Get(s.ID())
	require.True(t, ok, "session should be retrievable by its ID")
	assert.Same(t, s, got)

	reg.Remove(s.ID())
	_, ok = reg.Get(s.ID())
	assert.False(t, ok, "session should be removed by its ID")
}
