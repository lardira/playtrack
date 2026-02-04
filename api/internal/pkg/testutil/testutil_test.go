package testutil

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestInit(t *testing.T) {
	initFaker()

	assert.NotEqual(t, nil, faker)
	assert.NotZero(t, seed)

	got := Faker()
	assert.Equal(t, faker, got)

	gotSeed := GetSeed()
	assert.Equal(t, seed, gotSeed)
}
