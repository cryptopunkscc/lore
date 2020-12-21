package story

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
)

type Header struct {
	Version int
	Type    string
	Rel     []string
}

const MaxStorySize = 65535

// ParseHeaderFromBytes tries to parse a story header from data
func ParseHeaderFromBytes(data []byte) (*Header, error) {
	// Check data size
	if len(data) > MaxStorySize {
		return nil, ErrDataTooBig
	}

	// Try to parse the header
	story := struct{ Story *Header }{}
	err := yaml.Unmarshal(data, &story)
	if err != nil {
		return nil, fmt.Errorf("parser error: %w", err)
	}
	if story.Story == nil {
		return nil, ErrHeaderMissing
	}

	return story.Story, nil
}

// ParseHeaderFromFile tries to parse a story header from file
func ParseHeaderFromFile(file string) (*Header, error) {
	// Stat the file
	stat, err := os.Stat(file)
	if err != nil {
		return nil, err
	}

	// Check data size
	if stat.Size() > MaxStorySize {
		return nil, ErrDataTooBig
	}

	// Read the data
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// Parse the data
	return ParseHeaderFromBytes(data)
}

// ParseHeader reads all data from the provided reader and parses a story header
func ParseHeader(reader io.Reader) (*Header, error) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return ParseHeaderFromBytes(bytes)
}
