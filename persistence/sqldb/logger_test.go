package sqldb

import (
	"testing"
)

func Test_GorpLogger_Printf_Success(t *testing.T) {
	logger := NewGorpLogger()
	logger.Printf("%v,%v,%v", "test1", "test2", "test3")
}
