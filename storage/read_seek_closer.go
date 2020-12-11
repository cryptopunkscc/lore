package storage

import "io"

type ReadSeekCloser interface {
	io.Reader
	io.Seeker
	io.Closer
}
