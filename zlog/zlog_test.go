package zlog

import (
	"fmt"
	"testing"

	"github.com/DavidGamba/go-getoptions"
	"github.com/stretchr/testify/assert"
)

func TestOpts(t *testing.T) {
	osArgs := []string{"--log-level=debug", "--log-file=/tmp/1.log"}
	opts := getoptions.New()
	zl := NewWithOpts(opts)
	_, err := opts.Parse(osArgs)
	assert.NoError(t, err)
	assert.Equal(t, zl.optsLogLevel, "debug")
	assert.Equal(t, zl.optsLogFile, "/tmp/1.log")
	fmt.Println(zl.GetLogger())

}
