package worker

import "os"

// FileFilter returns true if path is a file
func FileFilter(path string, info os.FileInfo) bool {
	return !info.IsDir()
}

// DirFilter returns true if path is a directory
func DirFilter(path string, info os.FileInfo) bool {
	return info.IsDir()
}
