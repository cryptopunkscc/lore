package storage

import "errors"

var ErrAlreadyAdded = errors.New("file already added")
var ErrFileNotFound = errors.New("file not found")
