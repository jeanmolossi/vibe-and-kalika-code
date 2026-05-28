package installer

import "os"

func SafeFileMode(info os.FileInfo) os.FileMode {
	mode := info.Mode().Perm()
	if mode == 0 {
		return 0o644
	}
	return mode
}
