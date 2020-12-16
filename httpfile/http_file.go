package httpfile

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Interface checks
var _ io.ReadSeeker = &File{}
var _ io.Closer = &File{}

// File provides a Reader, Seeker and Closer interface over a HTTP url
type File struct {
	url  string
	body io.ReadCloser
	pos  int64
	len  int64
}

// Open opens a HTTP url as File
func Open(url string) (*File, error) {
	f := &File{
		url: url,
	}

	err := f.reopen(0)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// Len returns the file length, -1 if unknown
func (file File) Len() int64 {
	return file.len
}

// Read implements the io.Reader interface
func (file File) Read(p []byte) (n int, err error) {
	n, err = file.body.Read(p)
	if err != nil {
		return
	}

	file.pos = file.pos + int64(n)

	return
}

// Seek implements the io.Seeker interface
func (file *File) Seek(offset int64, whence int) (int64, error) {
	if file.len < 1 {
		return 0, errors.New("unsupported")
	}

	_ = file.body.Close()

	switch whence {
	case io.SeekStart:
		err := file.reopen(offset)
		if err != nil {
			return 0, err
		}
		return offset, nil
	case io.SeekCurrent:
		pos := file.pos + offset
		err := file.reopen(pos)
		if err != nil {
			return 0, err
		}
		return pos, nil
	case io.SeekEnd:
		err := file.reopen(-offset)
		return file.len - offset, err
	}

	return 0, errors.New("invalid whence parameter")
}

// Close implements the io.Closer interface
func (file File) Close() error {
	return file.body.Close()
}

func (file *File) reopen(offset int64) error {
	client := http.Client{}

	var httpRange string
	var pos int64
	if offset >= 0 {
		pos = offset
	} else {
		pos = file.len + offset
	}

	httpRange = fmt.Sprintf("bytes=%d-", pos)

	req, err := http.NewRequest("GET", file.url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Range", httpRange)

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if (res.StatusCode != http.StatusOK) && (res.StatusCode != http.StatusPartialContent) {
		return fmt.Errorf("http code %d", res.StatusCode)
	}

	// Remember content length if not yet known
	if file.len == 0 {
		file.len = res.ContentLength
	}

	file.body = res.Body
	file.pos = pos

	return nil
}
