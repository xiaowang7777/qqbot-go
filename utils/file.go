package utils

import "os"

func CheckFileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func CreateIfNotExists(path string) error {
	if !CheckFileExists(path) {
		file, err := os.Create(path)
		defer file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
