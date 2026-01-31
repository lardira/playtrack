package envutil

import (
	"os"
	"testing"

	"github.com/alecthomas/assert/v2"
)

const (
	testOSEnvKey   = "TEST"
	testOSEnvValue = "TEST"
)

func TestMustGet(t *testing.T) {
	var got string

	os.Setenv(testOSEnvKey, testOSEnvValue)

	assert.NotPanics(t, func() {
		got = MustGet(testOSEnvKey)
	})

	assert.Equal(t, testOSEnvValue, got)
}

func TestMustGet_NoValue(t *testing.T) {
	os.Setenv(testOSEnvKey, "")

	assert.Panics(t, func() {
		MustGet(testOSEnvKey)
	})
}

func TestGetOrDefault(t *testing.T) {
	os.Setenv(testOSEnvKey, testOSEnvValue)

	got := GetOrDefault(testOSEnvKey, "default")

	assert.Equal(t, testOSEnvValue, got)
}

func TestGetOrDefault_Default(t *testing.T) {
	os.Setenv(testOSEnvKey, "")

	got := GetOrDefault(testOSEnvKey, "default")

	assert.Equal(t, "default", got)
}

func TestLoadEnvs_InvalidPath(t *testing.T) {
	os.Setenv("ENV_PATH", ".invalid")

	err := LoadEnvs()
	assert.Error(t, err)
}
