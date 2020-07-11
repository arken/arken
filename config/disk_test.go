package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

//this isn't so much a test as a way to execute the code I want without executing
//the whole Arken program. Eventually will be more test-like.
func TestParsePoolSize(t *testing.T) {
	var test DiskInfo
	test.Init()
	ParsePoolSize(&test)
}

func TestDiskInfo_PrettyPoolSize(t *testing.T) {
	var dip DiskInfoProvider = &DiskInfo{}
	dip.SetPoolSizeBytes(30123019213)
	assert.Equal(t, "30.123019213GB", dip.GetPrettyPoolSize())
	dip.SetPoolSizeBytes(1003012000000)
	assert.Equal(t, "1.003012TB", dip.GetPrettyPoolSize())
	dip.SetPoolSizeBytes(10000000000000)
	assert.Equal(t,"10TB", dip.GetPrettyPoolSize())
}

func TestBytesToUnitString(t *testing.T) {
	assert.Equal(t, "10GB", BytesToUnitString(10000000000))
	assert.Equal(t, "100GB", BytesToUnitString(100000000000))
	assert.Equal(t, "1TB", BytesToUnitString(1000000000000))
	assert.Equal(t, "1.231KB", BytesToUnitString(1231))
	assert.Equal(t, "653B", BytesToUnitString(653))
}
