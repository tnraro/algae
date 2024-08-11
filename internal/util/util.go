package util

import (
	"os"
	"path"
)

func DataDir(v ...string) string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if path.IsAbs(os.Getenv("DATA_DIR")) {
		return path.Join(os.Getenv("DATA_DIR"), path.Join(v...))
	}
	return path.Join(cwd, os.Getenv("DATA_DIR"), path.Join(v...))
}
