package dirutil

import "os"

func RemoveFileIfExists(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return nil
	}
	return os.Remove(filename)
}

func RemoveDirIfExists(dirpath string) error {
	_, err := os.Stat(dirpath)
	if os.IsNotExist(err) {
		return nil
	}
	return os.RemoveAll(dirpath)
}
