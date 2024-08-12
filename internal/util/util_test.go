package util

import (
	"path"
	"testing"
)

func TestDataDir(t *testing.T) {
	dir := SetupDataDir(t)

	AssertEq(t, DataDir(), dir)
	AssertEq(t, DataDir("abc"), path.Join(dir, "abc"))
}
