package store

import "errors"

// ErrUnsupported means the store does not support the invoked method
var ErrUnsupported = errors.New("method unsupported")

// ErrNotFound means the requested file was not found in the store
var ErrNotFound = errors.New("not found")

/// ErrAlreadyExists means the file already exists in the store
var ErrAlreadyExists = errors.New("file already exists")
