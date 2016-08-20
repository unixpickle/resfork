// Package resfork provides Go programs with access to
// resource forks on files (generally on OS X or on
// filesystems written to by OS X).
package resfork

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// Open attempts to open the resource fork associated
// with the file at the given path.
func Open(path string) (*os.File, error) {
	return os.Open(filepath.Join(path, "..namedfork/rsrc"))
}

// Read attempts to fully read the resource fork
// associated with the file at the given path.
func Read(path string) ([]byte, error) {
	f, err := Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}
