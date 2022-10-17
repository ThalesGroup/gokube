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
	"errors"
	"fmt"
	"gopkg.in/cheggaaa/pb.v2"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"strings"
)

// GetAppDataHome ...
func GetAppDataHome() string {
	return os.Getenv("APPDATA")
}

// GetUserHome ...
func GetUserHome() string {
	userHome, err := user.Current()
	if err != nil {
		fmt.Println("Error: cannot determine user home directory")
		os.Exit(1)
	}
	return userHome.HomeDir
}

// GetBinDir ...
func GetBinDir(executable string) string {
	path, err := exec.LookPath(executable)
	if err != nil {
		fmt.Printf("Error: cannot determine %s directory\n", executable)
		os.Exit(1)
	}
	if errors.Is(err, exec.ErrDot) {
		path, err = os.Getwd()
		if err != nil {
			fmt.Printf("Error: cannot determine %s directory\n", executable)
			os.Exit(1)
		}
	} else {
		path = strings.TrimSuffix(path, string(os.PathSeparator)+"gokube.exe")
	}
	return path
}

// CreateDirs ...
func CreateDirs(dirPath string) error {
	// Check if file already exist
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

// CleanDir ...
func CleanDir(dirPath string) error {
	entries, err := os.ReadDir(dirPath)
	if err == nil {
		for _, e := range entries {
			err = os.RemoveAll(path.Join([]string{dirPath, e.Name()}...))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// DeleteDir ...
func DeleteDir(dirPath string) {
	err := os.RemoveAll(dirPath)
	if err != nil {
		fmt.Printf("Warning: cannot remove directory %s: %s", dirPath, err)
	}
}

func Close(stream io.Closer) {
	if stream != nil {
		err := stream.Close()
		if err != nil {
			fmt.Printf("Warning: cannot close stream: %s\n", err)
		}
	}
}

func CloseFile(file *os.File) {
	if file != nil {
		err := file.Close()
		if err != nil {
			fmt.Printf("Warning: cannot close file %s: %s\n", file.Name(), err)
		}
	}
}

func CloseGZipReader(reader *gzip.Reader) {
	if reader != nil {
		err := reader.Close()
		if err != nil {
			fmt.Printf("Warning: cannot close reader %s: %s\n", reader.Name, err)
		}
	}
}

func CloseZipReader(reader *zip.ReadCloser) {
	if reader != nil {
		err := reader.Close()
		if err != nil {
			fmt.Printf("Warning: cannot close reader: %s\n", err)
		}
	}
}

func ClosePBReader(reader *pb.Reader) {
	if reader != nil {
		err := reader.Close()
		if err != nil {
			fmt.Printf("Warning: cannot close reader: %s\n", err)
		}
	}
}

// Untar takes a destination path and a reader; a tar reader loops over the tarfile
// creating the file structure at 'dst' along the way, and writing any files
func Untar(src string, dst string) error {

	file, err := os.Open(src)
	defer CloseFile(file)
	if err != nil {
		panic(err)
	}
	gzr, err := gzip.NewReader(bufio.NewReader(file))
	defer CloseGZipReader(gzr)
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

		// check the file type
		switch header.Typeflag {
		// if it's a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}
		// if it's a file create it
		case tar.TypeReg:
			if _, err := os.Stat(filepath.Dir(target)); err != nil {
				if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
					return err
				}
			}
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}
			err = f.Close()
			if err != nil {
				return err
			}
		}
	}
}

// Unzip ...
func Unzip(src string, dest string) error {

	var fileNames []string

	r, err := zip.OpenReader(src)
	defer CloseZipReader(r)
	if err != nil {
		return err
	}

	for _, f := range r.File {

		rc, err := f.Open()
		if err != nil {
			return err
		}

		// Store filename/path for returning and using later on
		fileName := filepath.Join(dest, f.Name)
		fileNames = append(fileNames, fileName)

		if f.FileInfo().IsDir() {
			// Make Folder
			err = os.MkdirAll(fileName, 0755)
			if err != nil {
				return err
			}
		} else {
			// Make File
			if err = os.MkdirAll(filepath.Dir(fileName), 0755); err != nil {
				return err
			}
			outFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			_, err = io.Copy(outFile, rc)
			err = outFile.Close()
			if err != nil {
				return err
			}
		}
		err = rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// GetValueFromEnv ...
func GetValueFromEnv(envVar string, defaultValue string) string {
	var value = os.Getenv(envVar)
	if len(value) == 0 {
		value = defaultValue
	}
	return value
}
