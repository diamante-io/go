package http

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"sort"
)

// EqualFileSystems traverses two http.FileSystem instances and returns true
// if they are equal
func EqualFileSystems(fs, otherFS http.FileSystem, currentPath string) bool {
	file, err := fs.Open(currentPath)
	if err != nil {
		return false
	}

	otherFile, err := otherFS.Open(currentPath)
	if err != nil {
		return false
	}

	stat, err := file.Stat()
	if err != nil {
		return false
	}

	otherStat, err := otherFile.Stat()
	if err != nil {
		return false
	}

	if stat.IsDir() != otherStat.IsDir() {
		return false
	}

	if !stat.IsDir() {
		var fileBytes, otherFileBytes []byte
		fileBytes, err = ioutil.ReadAll(file)
		if err != nil {
			return false
		}

		otherFileBytes, err = ioutil.ReadAll(otherFile)
		if err != nil {
			return false
		}

		return bytes.Equal(fileBytes, otherFileBytes)
	}

	files, err := file.Readdir(0)
	if err != nil {
		return false
	}

	otherFiles, err := otherFile.Readdir(0)
	if err != nil {
		return false
	}

	if len(files) != len(otherFiles) {
		return false
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})
	sort.Slice(otherFiles, func(i, j int) bool {
		return otherFiles[i].Name() < otherFiles[j].Name()
	})

	for i := range files {
		if files[i].Name() != otherFiles[i].Name() {
			return false
		}

		nextPath := filepath.Join(currentPath, files[i].Name())

		if !EqualFileSystems(fs, otherFS, nextPath) {
			return false
		}
	}

	return true
}
