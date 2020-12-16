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

// ParseHeader tries to parse a story header from data
func ParseHeader(data []byte) (*Header, error) {
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
	return ParseHeader(data)
}

// ParseHeaderFromReader reads all data from the provided reader and parses a story header
func ParseHeaderFromReader(reader io.Reader) (*Header, error) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return ParseHeader(bytes)
}

// ParseStory parses the story from bytes into provided obj file using YAML Unmarshal
func ParseStory(data []byte, obj interface{}) error {
	err := yaml.Unmarshal(data, obj)
	if err != nil {
		return fmt.Errorf("error parsing story: %w", err)
	}
	return nil
}
