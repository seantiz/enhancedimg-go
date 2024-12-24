package enhancedimg

import (
	"os"
)

// EnhanceImages acts as the public API to use this package as a dependency
// path argument should be the path holding all html templates you want to be processed
func EnhanceImages(path string) error {
	if err := os.MkdirAll("static/processed", 0755); err != nil {
		return err
	}
	return FindAllImageElements(path)
}
