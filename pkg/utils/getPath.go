package utils

import (
	"path"
	"runtime"
)

// GetCurrentPath relative 表示层级，为1是当前，2向上跳一级，以此类推
func GetCurrentPath(relative int) string {
	_, filename, _, _ := runtime.Caller(relative)
	return path.Dir(filename)
}
