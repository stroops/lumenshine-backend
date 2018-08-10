package helpers

import (
	"fmt"
	"os"
)

//CreateDirIfNotExists creates an directory if it does not exists
func CreateDirIfNotExists(path string) error {
	if !ExistDir(path) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Printf(fmt.Sprintf("Could not create dir <%s>\n", path))
			return err
		}
	}
	return nil
}

//ExistDir checks if an directory exists
func ExistDir(path string) bool {
	if stat, err := os.Stat(path); err == nil && stat.IsDir() {
		return true
	}
	return false
}
