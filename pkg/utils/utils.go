/*
(c) Copyright 2018, Gemalto. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
)

// EnvVar ...
type EnvVar struct {
	Name  string
	Value string
}

// CreateFile ...
func CreateFile(filePath string) {
	var _, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		var file, err = os.Create(filePath)
		if err != nil {
			panic(err)
		}
		defer file.Close()
	}
}

// CreateDir ...
func CreateDir(dirPath string) {
	// Check if file already exist
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.Mkdir(dirPath, 0001)
		if err != nil {
			panic(err)
		}
	}
}

// MoveFile ...
func MoveFile(oldPath string, newPath string) {
	err := os.Rename(oldPath, newPath)
	if err != nil {
		panic(err)
	}
}

// MoveFiles ...
func MoveFiles(oldPath string, newPath string) {
	files, err := filepath.Glob(oldPath)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if err := os.Rename(f, newPath); err != nil {
			panic(err)
		}
	}
}

func CopyFile(src, dst string) error {
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

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

// CleanDir ...
func CleanDir(dirPath string) {
	dir, err := ioutil.ReadDir(dirPath)
	if err == nil {
		for _, e := range dir {
			err := os.RemoveAll(path.Join([]string{dirPath, e.Name()}...))
			if err != nil {
				panic(err)
			}
		}
	}
}

// RemoveDir ...
func RemoveDir(dirPath string) {
	err := os.RemoveAll(dirPath)
	if err != nil {
		panic(err)
	}
}

// RemoveFile ...
func RemoveFile(filePath string) {
	err := os.RemoveAll(filePath)
	if err != nil {
		panic(err)
	}
}

// RemoveFiles ...
func RemoveFiles(filePath string) {
	files, err := filepath.Glob(filePath)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			panic(err)
		}
	}
}

// Untar takes a destination path and a reader; a tar reader loops over the tarfile
// creating the file structure at 'dst' along the way, and writing any files
func Untar(src string, dst string) error {

	file, err := os.Open(src)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	bufio.NewReader(file)

	gzr, err := gzip.NewReader(bufio.NewReader(file))
	defer gzr.Close()
	if err != nil {
		return err
	}

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			defer f.Close()

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}
		}
	}
}

// Unzip ...
func Unzip(src string, dest string) error {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {

		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)
		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {

			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)

		} else {

			// Make File
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}

			_, err = io.Copy(outFile, rc)

			// Close the file without defer to close before next iteration of loop
			outFile.Close()

			if err != nil {
				return err
			}

		}
	}
	return nil
}

// GetUserHome ...
func GetUserHome() string {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	return user.HomeDir
}

// WriteFile ...
func WriteFile(content string, path string) {
	d1 := []byte(content)
	err := ioutil.WriteFile(path, d1, 0644)
	if err != nil {
		panic(err)
	}
}
