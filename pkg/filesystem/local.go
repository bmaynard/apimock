package filesystem

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	m "github.com/bmaynard/apimock/pkg/mocks"
	l "github.com/bmaynard/apimock/pkg/utils/logger"
)

type LocalFileOptions struct {
	Path string
}

func NewLoclFileOptions() *LocalFileOptions {
	return &LocalFileOptions{
		Path: "",
	}
}

func (o *LocalFileOptions) SetOptionString(key string, value string) {
	err := setField(o, key, value)

	if err != nil {
		l.Log.Fatal(err)
	}
}

func (o *LocalFileOptions) GetMocks() []FileMock {
	folders, err := o.readDir()

	if err != nil {
		l.Log.Fatal(err)
	}

	var fileMocks []FileMock

	for _, folder := range folders {
		var files []string
		rootPath := fmt.Sprintf("%s/%s", o.Path, folder.Name())
		err := filepath.Walk(rootPath, visit(&files))

		if err != nil {
			l.Log.Fatal(err)
		}

		for _, filePath := range files {
			if filepath.Ext(filePath) != ".json" {
				continue
			}

			file, err := ioutil.ReadFile(filePath)

			if err != nil {
				l.Log.Fatal(err)
			}

			fileMocks = append(fileMocks, FileMock{
				FilePath: filePath,
				Domain:   filepath.Base(filepath.Dir(filePath)),
				Contents: file,
			})
			l.Log.Debugf("Added %s to potential mocks list", filePath)
		}
	}

	return fileMocks
}

func (o *LocalFileOptions) WriteMockFile(r *http.Response, bodyBytes []byte, originalHost string) error {
	if r == nil {
		return nil
	}

	meta := m.Meta{
		StatusCode:    r.StatusCode,
		RequestMethod: r.Request.Method,
		RequestPath:   r.Request.URL.Path,
	}

	mock := m.MockResponse{Meta: meta}

	if err := json.Unmarshal(bodyBytes, &mock.Response); err != nil {
		return err
	}

	file, _ := json.MarshalIndent(mock, "", " ")

	h := md5.New()
	h.Write(file)

	requestPath := r.Request.URL.EscapedPath()
	requestPath = strings.Replace(requestPath, "/", "_", -1)
	folderName := fmt.Sprintf("%s/%s", o.Path, originalHost)
	fileName := fmt.Sprintf("%s/%s_%s.json", folderName, requestPath, hex.EncodeToString(h.Sum(nil)))

	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		os.Mkdir(folderName, 0700)
	}

	err := ioutil.WriteFile(fileName, file, 0644)

	if err == nil {
		l.Log.Infof("Saved mock file to %s", fileName)
	}

	return err
}

func (o *LocalFileOptions) readDir() ([]os.FileInfo, error) {
	f, err := os.Open(o.Path)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })
	return list, nil
}

func visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			l.Log.Fatal(err)
		}
		if info.IsDir() {
			return nil
		}
		*files = append(*files, path)
		return nil
	}
}
