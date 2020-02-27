package files

import "os"

//In this file, is everything required to find and traverse the file directories.

//Exists returns whether a file with a given file exists.
func Exists(fn string) bool {
	if _, err := os.Stat(fn); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false

	} else {
		// Schrodinger: file may or may not exist. See err for details.
		return false

	}
}
