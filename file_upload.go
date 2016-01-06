package uploadSrv

import (
	"io"
	"net/http"
	"os"
)

const (
	FILE_PATH = "/var/"
	URL       = "127.0.0.1/var/"
)

type FileUpload struct {
	req      *http.Request
	formName string
	path     string
	url      string
}

func NewFileUpload(req *http.Request, formName string) *FileUpload {
	req.ParseMultipartForm(32 << 20)
	_, handler, err := req.FormFile(formName)
	if err != nil {
		return nil
	}
	return &FileUpload{
		req:      req,
		formName: formName,
		path:     FILE_PATH + handler.Filename,
		url:      URL + handler.Filename,
	}
}

func (fu *FileUpload) GetPath() string {
	return fu.path
}

func (fu *FileUpload) GetUrl() string {
	return fu.url
}
func (fu *FileUpload) fileNameHandle() error {
	_, err := os.Stat(fu.path)
	if err == nil || os.IsExist(err) {
		err = os.Remove(fu.path)
		if err != nil {
			return err
		}
	}
	f, err := os.Create(fu.path)
	defer f.Close()
	return err
}

func (fu *FileUpload) Upload() error {
	file, _, err := fu.req.FormFile(fu.formName)
	if err != nil {
		return err
	}
	defer func() {
		file.Close()
		return
	}()

	err = fu.fileNameHandle()
	if err != nil {
		return err
	}
	f, err := os.OpenFile(fu.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, file)

	return err
}
