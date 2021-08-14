package utils

import "os"

func CheckFileExists(path string) bool {
	_, err := os.Stat(path)
	return err != nil || os.IsExist(err)
}
