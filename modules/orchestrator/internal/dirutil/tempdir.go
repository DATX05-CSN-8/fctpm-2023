package dirutil

import (
	"os"
	"strings"

	"github.com/google/uuid"
)

func JoinPath(paths ...string) string {
	return strings.Join(paths, string(os.PathSeparator))
}

func EnsureDirectory(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	return err
}

func CreateTempDirWithId(base string, id string) (string, error) {
	path := JoinPath(base, id)
	err := EnsureDirectory(path)
	if err != nil {
		return "", err
	}
	return path, nil
}
func CreateTempDir(base string) (string, error) {
	id := uuid.NewString()
	return CreateTempDirWithId(base, id)
}

func RemoveTempDir(path string) error {
	return os.RemoveAll(path)
}
