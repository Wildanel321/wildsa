//go:build !linux

package system

func getDiskUsageLinux(path string) (totalBytes uint64, freeBytes uint64) {
	return 100 * 1024 * 1024 * 1024, 60 * 1024 * 1024 * 1024
}
