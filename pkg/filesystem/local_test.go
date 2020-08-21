package filesystem

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGetLocalMocks(t *testing.T) {
	fa := NewLocalFileOptions()
	fa.SetOptionString("Path", getMockPath(t))

	if len(fa.GetMocks()) == 0 {
		t.Error("No mocks returned")
	}
}

func TestGetWriteMock(t *testing.T) {
	fa := NewLocalFileOptions()
	fa.SetOptionString("Path", getMockPath(t))

	res, err := http.Get("https://jsonplaceholder.typicode.com/users/1")

	if err != nil {
		t.Error(err)
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		t.Error(err)
	}

	domain := fmt.Sprintf("%s.typicode.com", randString(10))
	fa.WriteMockFile(res, bodyBytes, domain)

	files, _ := ioutil.ReadDir(filepath.Join(getMockPath(t), domain))
	if len(files) != 1 {
		t.Error("Mock files not saved to domain directory")
	}
}

func getMockPath(t *testing.T) string {
	pwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	return filepath.Join(pwd, "..", "..", "tests", "mocks")
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
