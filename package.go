package annotate

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

// Package is list of files
type Package []File

var (
	lock = sync.Mutex{}
	p    Package
)

func walkPackage(path string, f os.FileInfo, err error) error {
	if !f.IsDir() {
		if filepath.Ext(path) == ".go" {
			r, err := os.Open(path)
			if err != nil {
				return err
			}
			data, err := ioutil.ReadAll(r)
			if err != nil {
				return err
			}

			f, err := ParseFile(string(data))
			if err != nil {
				return err
			}
			f.FileName = path
			p = append(p, f)
		}
	}

	return nil
}

// ParsePackage is here for loading a single package and parse all files in it
func ParsePackage(path string) (Package, error) {
	lock.Lock()
	defer lock.Unlock()
	p = nil
	err := filepath.Walk(path, walkPackage)
	if err != nil {
		return nil, err
	}

	return p, nil
}
