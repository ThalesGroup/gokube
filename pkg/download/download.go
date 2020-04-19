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
	"strings"
	"time"

	"github.com/gemalto/gokube/pkg/gokube"
	"github.com/gemalto/gokube/pkg/utils"
	"gopkg.in/cheggaaa/pb.v2"
)

// FromUrl ...
func FromUrl(name string, tpl string, version string) int64 {

	url := tpl

	if strings.Contains(tpl, "%s") {
		url = fmt.Sprintf(tpl, version)
	}

	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]

	// Templates definition: https://github.com/cheggaaa/pb/blob/v2/preset.go
	var tmpl string
	//tmpl = `{{ yellow "` + name + `: " }}{{counters . }} {{bar . "[" "=" ">" "_" "]" | green }} {{percent . }} {{speed . }}`
	tmpl = `{{ yellow "` + name + `: " }}{{counters . }} {{bar . | green }} {{percent . }} {{speed . }}`

	utils.CreateDirs(gokube.GetTempDir())
	output, err := os.Create(gokube.GetTempDir() + "/" + fileName)
	if err != nil {
		panic(err)
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	filesize := response.ContentLength
	count := int(filesize)
	bar := pb.ProgressBarTemplate(tmpl).Start(count)
	//bar := pb.StartNew(count)
	bar.SetWidth(100)
	defer bar.Finish()

	// create proxy reader
	reader := bar.NewProxyReader(response.Body)
	n, err := io.Copy(output, reader)

	defer reader.Close()
	if err != nil {
		panic(err)
	}

	var fi os.FileInfo
	for fi == nil || int(fi.Size()) < count {
		fi, _ = output.Stat()
		bar.Increment()
		time.Sleep(time.Millisecond)
	}

	tokens = strings.Split(fileName, ".")
	fileType := tokens[len(tokens)-1]
	switch fileType {
	case "zip":
		if err = utils.Unzip(gokube.GetTempDir()+"/"+fileName, gokube.GetTempDir()); err != nil {
			panic(err)
		}
	case "tgz":
		if err = utils.Untar(gokube.GetTempDir()+"/"+fileName, gokube.GetTempDir()); err != nil {
			panic(err)
		}
	case "gz":
		if err = utils.Untar(gokube.GetTempDir()+"/"+fileName, gokube.GetTempDir()); err != nil {
			panic(err)
		}
	}

	return n
}
