package common

import (
	"io"
	"os"
)

// MoveFile tries to rename the file, falling back to copy-and-delete.
func MoveFile(src, dst string) error {
	// Try renaming first.
	if err := os.Rename(src, dst); err == nil {
		return nil
	}
	// Fallback: copy the file.
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	if err := out.Close(); err != nil {
		return err
	}
	return os.Remove(src)
}
