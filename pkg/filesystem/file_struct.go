package filesystem

import "net/http"

type FileMock struct {
	FilePath string
	Contents []byte
	Domain   string
}

type Filesystem interface {
	GetMocks() []FileMock
	SetOptionString(key string, value string)
	WriteMockFile(r *http.Response, bodyBytes []byte, originalHost string) error
}
