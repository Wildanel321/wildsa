//go:build linux

package system

import "syscall"

func getDiskUsageLinux(path string) (totalBytes uint64, freeBytes uint64) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return 100 * 1024 * 1024 * 1024, 60 * 1024 * 1024 * 1024
	}
	totalBytes = stat.Blocks * uint64(stat.Bsize)
	freeBytes = stat.Bavail * uint64(stat.Bsize)
	return totalBytes, freeBytes
}
