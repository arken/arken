package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

//this isn't so much a test as a way to execute the code I want without executing
//the whole Arken program. Eventually will be more test-like.
func TestParsePoolSize(t *testing.T) {
	GlobalDiskInfo.Init()
	ParsePoolSize(&GlobalDiskInfo)
}

func TestNumLen(t *testing.T) {
	assert.Equal(t, 0, numLen(0))
	assert.Equal(t, 2, numLen(10))
}

func TestGetUnit(t *testing.T) {
	assert.Equal(t, "TB", getUnit(12))
	assert.Equal(t, "GB", getUnit(9))
	assert.Equal(t, "MB", getUnit(6))
	assert.Equal(t, "KB", getUnit(3))
	assert.Equal(t, "B", getUnit(0))
	assert.Equal(t, "", getUnit(-1))
}

