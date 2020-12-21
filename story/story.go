package story

import (
	"fmt"
	"gopkg.in/yaml.v2"
)

// ParseStory parses the story from bytes into provided obj file using YAML Unmarshal
func ParseStory(data []byte, obj interface{}) error {
	err := yaml.Unmarshal(data, obj)
	if err != nil {
		return fmt.Errorf("error parsing story: %w", err)
	}
	return nil
}
