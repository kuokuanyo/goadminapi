package file

import (
	"goadminapi/plugins/admin/modules"
	"io"
	"mime/multipart"
	"os"
	"path"
	"sync"
)

var copyBufPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 4096)
	},
}

var uploaderList = map[string]UploaderGenerator{
	"local": GetLocalFileUploader,
}

// UploaderGenerator is a function return an Uploader.
type UploaderGenerator func() Uploader

// UploadFun is a function to process the uploading logic.
type UploadFun func(*multipart.FileHeader, string) (string, error)

// Uploader is a file uploader which contains the method Upload.
type Uploader interface {
	Upload(*multipart.Form) error
}

// GetFileEngine return the Uploader of given name.
func GetFileEngine(name string) Uploader {
	if up, ok := uploaderList[name]; ok {
		return up()
	}
	panic("wrong uploader name")
}

func copyZeroAlloc(w io.Writer, r io.Reader) (int64, error) {
	buf := copyBufPool.Get().([]byte)
	n, err := io.CopyBuffer(w, r, buf)
	copyBufPool.Put(buf)
	return n, err
}

// SaveMultipartFile used in a local Uploader which help to save file in the local path.
func SaveMultipartFile(fh *multipart.FileHeader, path string) error {
	f, err := fh.Open()
	if err != nil {
		return err
	}

	if ff, ok := f.(*os.File); ok {
		// Windows can't rename files that are opened.
		if err := f.Close(); err != nil {
			return err
		}

		// If renaming fails we try the normal copying method.
		// Renaming could fail if the files are on different devices.
		if os.Rename(ff.Name(), path) == nil {
			return nil
		}

		// Reopen f for the code below.
		f, err = fh.Open()
		if err != nil {
			return err
		}
	}

	defer func() {
		if err2 := f.Close(); err2 != nil {
			err = err2
		}
	}()

	ff, err := os.Create(path)
	if err != nil {
		return err
	}

	defer func() {
		if err2 := ff.Close(); err2 != nil {
			err = err2
		}
	}()
	_, err = copyZeroAlloc(ff, f)
	return err
}

// Upload receive the return value of given UploadFun and put them into the form.
func Upload(c UploadFun, form *multipart.Form) error {
	var (
		suffix   string
		filename string
	)

	for k := range form.File {
		for _, fileObj := range form.File[k] {
			suffix = path.Ext(fileObj.Filename)
			filename = modules.Uuid() + suffix

			pathStr, err := c(fileObj, filename)

			if err != nil {
				return err
			}

			form.Value[k] = append(form.Value[k], pathStr)
		}
	}

	return nil
}
