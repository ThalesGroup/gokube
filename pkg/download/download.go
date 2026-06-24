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

package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gemalto/gokube/pkg/utils"
	pb "github.com/cheggaaa/pb/v3"
)

type FileMap struct {
	Src string
	Dst string
}

func fromUrl(url string, name string, dir string, fileName string, bar *pb.ProgressBar) (n int64, retErr error) {
	file, err := os.Create(dir + string(os.PathSeparator) + fileName)
	defer utils.CloseFile(file)
	if err != nil {
		return -1, err
	}

	response, err := http.Get(url)
	defer utils.Close(response.Body)
	if err != nil {
		return -1, err
	}
	if response.StatusCode != 200 {
		return -1, fmt.Errorf("cannot download %s", url)
	}

	count := int(response.ContentLength)
	tmpl := `{{ yellow "` + name + `: " }}{{counters . }} {{bar . | green }} {{percent . }} {{speed . }} {{etime .}}`
	bar.SetTemplateString(tmpl)
	bar.SetTotal(int64(count))
	bar.SetWidth(100)
	dlStart := time.Now()
	defer func() {
		if retErr == nil {
			d := time.Since(dlStart).Round(time.Second)
			bar.SetTemplateString(`{{ green "` + name + `" }} done (` + d.String() + `)`)
		}
		bar.Finish()
	}()

	// create proxy reader
	reader := bar.NewProxyReader(response.Body)
	defer utils.ClosePBReader(reader)
	n, err = io.Copy(file, reader)
	if err != nil {
		return -1, err
	}

	var fi os.FileInfo
	for fi == nil || int(fi.Size()) < count {
		fi, _ = file.Stat()
		bar.Increment()
		time.Sleep(time.Millisecond)
	}

	tokens := strings.Split(fileName, ".")
	fileType := tokens[len(tokens)-1]
	switch fileType {
	case "zip":
		if err = utils.Unzip(file.Name(), dir); err != nil {
			return -1, err
		}
	case "tgz":
		if err = utils.Untar(file.Name(), dir); err != nil {
			return -1, err
		}
	case "gz":
		if err = utils.Untar(file.Name(), dir); err != nil {
			return -1, err
		}
	}
	return n, nil
}

// VersionFile returns the path of the version metadata file for a given binary path.
// All version files are stored in ~/.gokube/metadata/<toolname>.version.
func VersionFile(binaryPath string) string {
	name := strings.TrimSuffix(filepath.Base(binaryPath), ".exe")
	return filepath.Join(utils.GetUserHome(), ".gokube", "metadata", name+".version")
}

// IsCurrentVersion returns true when binaryPath exists and its metadata file records version.
func IsCurrentVersion(binaryPath, version string) bool {
	if _, err := os.Stat(binaryPath); err != nil {
		return false
	}
	data, err := os.ReadFile(VersionFile(binaryPath))
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(data)) == version
}

// WriteVersion writes version to the metadata file for binaryPath.
func WriteVersion(binaryPath, version string) error {
	vf := VersionFile(binaryPath)
	if err := os.MkdirAll(filepath.Dir(vf), 0755); err != nil {
		return err
	}
	return os.WriteFile(vf, []byte(version), 0644)
}

// DeleteAllMetadata removes the entire ~/.gokube/metadata/ directory.
// Call during a full clean to ensure no stale version records survive.
func DeleteAllMetadata() error {
	return os.RemoveAll(filepath.Join(utils.GetUserHome(), ".gokube", "metadata"))
}

// moveFile moves src to dst, falling back to copy+delete when os.Rename
// fails across drive boundaries or filesystem boundaries.
func moveFile(src, dst string) error {
	if err := os.Rename(src, dst); err == nil {
		return nil
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		in.Close()
		return err
	}
	_, err = io.Copy(out, in)
	in.Close() // must close before os.Remove on Windows — open handles block deletion
	if err != nil {
		out.Close()
		os.Remove(dst)
		return err
	}
	if err = out.Close(); err != nil {
		os.Remove(dst)
		return err
	}
	return os.Remove(src)
}

// FromUrl ...
func FromUrl(urlTpl string, version string, name string, fileMaps []*FileMap, dst string, bar *pb.ProgressBar) (int64, error) {

	url := strings.Replace(urlTpl, "%s", version, -1)
	if version[0:1] == "v" {
		name = name + " " + version
	} else {
		name = name + " v" + version
	}
	tokens := strings.Split(url, "/")
	urlFileName := tokens[len(tokens)-1]

	tempDir, err := os.MkdirTemp(os.TempDir(), "*")
	defer utils.DeleteDir(tempDir)

	n, err := fromUrl(url, name, tempDir, urlFileName, bar)
	if err != nil {
		return -1, err
	}

	for _, fileMap := range fileMaps {
		fileDst := dst + string(os.PathSeparator) + fileMap.Dst
		if _, err := os.Stat(filepath.Dir(fileDst)); err != nil {
			if err := os.MkdirAll(filepath.Dir(fileDst), 0755); err != nil {
				return -1, err
			}
		}
		fileSrc := tempDir + string(os.PathSeparator) + fileMap.Src
		err = moveFile(fileSrc, fileDst)
		if err != nil {
			return -1, err
		}
	}

	return n, nil
}
