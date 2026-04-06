package device

import (
	"bytes"
	"errors"
)

var errorPatterns = []string{
	"Command::Base error",
	"Mobile::Interface error",
}

func errorLookup(raw []byte) error {
	for i := range errorPatterns {
		if idx := bytes.Index(raw, []byte(errorPatterns[i])); idx != -1 {
			return errors.New(string(bytes.TrimSpace(bytes.Clone(raw[idx:]))))
		}
	}

	return nil
}
