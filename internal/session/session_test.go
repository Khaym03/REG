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
