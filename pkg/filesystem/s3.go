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
	"strings"

	m "github.com/bmaynard/apimock/pkg/mocks"
	l "github.com/bmaynard/apimock/pkg/utils/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3FileOptions struct {
	Bucket string
	Region string
	Prefix string
}

func NewS3FileOptions() *S3FileOptions {
	return &S3FileOptions{
		Bucket: "",
		Region: "",
		Prefix: "",
	}
}

func (o *S3FileOptions) SetOptionString(key string, value string) {
	err := setField(o, key, value)

	if err != nil {
		l.Log.Fatal(err)
	}
}

func (o *S3FileOptions) GetMocks() []FileMock {
	var fileMocks []FileMock

	sess, err := o.getS3Session()

	if err != nil {
		l.Log.Fatal(err)
	}

	svc := s3.New(sess)
	input := &s3.ListObjectsInput{
		Bucket: aws.String(o.Bucket),
		Prefix: aws.String(o.Prefix),
	}

	result, err := svc.ListObjects(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				l.Log.Fatal(s3.ErrCodeNoSuchBucket, aerr.Error())
			default:
				l.Log.Fatal(aerr.Error())
			}
		} else {
			l.Log.Fatal(err.Error())
		}
	}

	downloader := s3manager.NewDownloader(sess)

	for _, object := range result.Contents {
		tmpfile, err := ioutil.TempFile("", "apimock")
		if err != nil {
			l.Log.Fatal(err)
		}

		defer os.Remove(tmpfile.Name())

		_, err = downloader.Download(tmpfile, &s3.GetObjectInput{
			Bucket: aws.String(o.Bucket),
			Key:    aws.String(*object.Key),
		})
		if err != nil {
			l.Log.Fatalf("Failed to download file from S3, %v", err)
		}

		file, err := ioutil.ReadFile(tmpfile.Name())

		if err != nil {
			l.Log.Fatal(err)
		}

		fileMocks = append(fileMocks, FileMock{
			FilePath: tmpfile.Name(),
			Domain:   filepath.Base(filepath.Dir(*object.Key)),
			Contents: file,
		})
		l.Log.Debugf("Added %s to potential mocks list", *object.Key)
	}

	return fileMocks
}

func (o *S3FileOptions) WriteMockFile(r *http.Response, bodyBytes []byte, originalHost string) error {
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

	file, err := json.MarshalIndent(mock, "", " ")

	if err != nil {
		return err
	}

	h := md5.New()
	h.Write(file)

	requestPath := r.Request.URL.EscapedPath()
	requestPath = strings.Replace(requestPath, "/", "_", -1)
	folderName, err := ioutil.TempDir("", "apimock")
	if err != nil {
		return err
	}

	defer os.RemoveAll(folderName)

	filePart := fmt.Sprintf("%s_%s.json", requestPath, hex.EncodeToString(h.Sum(nil)))
	fileName := filepath.Join(folderName, filePart)
	err = ioutil.WriteFile(fileName, file, 0644)

	if err != nil {
		return err
	}

	sess, err := o.getS3Session()

	if err != nil {
		return err
	}

	uploader := s3manager.NewUploader(sess)
	f, err := os.Open(fileName)

	if err != nil {
		return err
	}

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(o.Bucket),
		Key:         aws.String(strings.Join([]string{o.Prefix, originalHost, filePart}, "/")),
		ContentType: aws.String("application/json"),
		Body:        f,
	})

	if err == nil {
		l.Log.Infof("Saved mock file to %s", result.Location)
	}

	return err
}

func (o *S3FileOptions) getS3Session() (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region: aws.String(o.Region)},
	)
}
