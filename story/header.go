package story

import (
	"errors"
	"fmt"
	"github.com/cryptopunkscc/lore/id"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
)

type Header struct {
	Version int
	Type    string
	Rel     []string `yaml:",omitempty"`
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
	err = story.Story.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid header: %w", err)
	}

	return story.Story, nil
}

func (header Header) Validate() error {
	if header.Version != 0 {
		return errors.New("invalid version")
	}
	if header.Type == "" {
		return errors.New("missing type")
	}
	for _, edge := range header.Rel {
		if _, err := id.Parse(edge); err != nil {
			return errors.New("invalid edge")
		}
	}
	return nil
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
