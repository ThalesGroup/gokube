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
	"gopkg.in/cheggaaa/pb.v2"
)

type FileMap struct {
	Src string
	Dst string
}

func fromUrl(url string, name string, dir string, fileName string) (int64, error) {
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
	tmpl := `{{ yellow "` + name + `: " }}{{counters . }} {{bar . | green }} {{percent . }} {{speed . }}`
	bar := pb.ProgressBarTemplate(tmpl).Start(count)
	defer bar.Finish()
	bar.SetWidth(100)

	// create proxy reader
	reader := bar.NewProxyReader(response.Body)
	defer utils.ClosePBReader(reader)
	n, err := io.Copy(file, reader)
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

// FromUrl ...
func FromUrl(urlTpl string, version string, name string, fileMaps []*FileMap, dst string) (int64, error) {

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

	n, err := fromUrl(url, name, tempDir, urlFileName)
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
		err = os.Rename(fileSrc, fileDst)
		if err != nil {
			return -1, err
		}
	}

	return n, nil
}
