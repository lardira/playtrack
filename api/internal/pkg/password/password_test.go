package password

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/lardira/playtrack/internal/pkg/testutil"
)

const (
	minPossiblePassLen = 5
)

func TestHash(t *testing.T) {
	password := testutil.Faker().Password(true, true, true, true, false, minPossiblePassLen)

	hash, err := Hash(password)
	assert.NoError(t, err)
	assert.NotZero(t, hash)

	ok := CompareHash(password, hash)
	assert.True(t, ok)
}

func TestHash_EmptyPassword(t *testing.T) {
	password := ""

	hash, err := Hash(password)
	assert.IsError(t, ErrEmptyPass, err)
	assert.Zero(t, hash)

	ok := CompareHash(password, hash)
	assert.False(t, ok)
}
