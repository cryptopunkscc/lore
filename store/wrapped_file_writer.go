package store

// Interface check
var _ Writer = &WrappedWriter{}

type FinalizeFunc func(string, error) error

// WrappedWriter wraps a Writer to add an additional callback on Finalize() call.
type WrappedWriter struct {
	writer     Writer
	onFinalize FinalizeFunc
}

// NewWrappedWriter makes a new WrappedWriter
func NewWrappedWriter(writer Writer, onFinalize FinalizeFunc) *WrappedWriter {
	return &WrappedWriter{writer: writer, onFinalize: onFinalize}
}

// Write pass-through
func (w *WrappedWriter) Write(data []byte) (int, error) {
	return w.writer.Write(data)
}

// Finalize will call Finalize on the underlying writer and pass the results to the provided callback
func (w *WrappedWriter) Finalize() (string, error) {
	id, err := w.Finalize()

	if w.onFinalize != nil {
		err2 := w.onFinalize(id, err)
		if err2 != nil {
			return "", err2
		}
	}

	return id, err
}

// Write pass-through
func (w *WrappedWriter) Discard() error {
	return w.Discard()
}
