package story

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
)

// Parse parses a story into obj from an io.Reader
func Parse(reader io.Reader, obj interface{}) error {
	var buf [MaxStorySize + 1]byte

	n, err := reader.Read(buf[:])
	if err != nil {
		return err
	}

	if n > MaxStorySize {
		return ErrDataTooBig
	}

	return ParseBytes(buf[:n], obj)
}

// ParseBytes parses the story from bytes into provided obj file using YAML Unmarshal
func ParseBytes(data []byte, obj interface{}) error {
	err := yaml.Unmarshal(data, obj)
	if err != nil {
		return fmt.Errorf("error parsing story: %w", err)
	}
	return nil
}
