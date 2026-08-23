package store

import (
	"encoding/json"
	"fmt"
)

func encode(value any) ([]byte, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("encode entity: %w", err)
	}
	return data, nil
}

func decode(data []byte, target any) error {
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("decode entity: %w", err)
	}
	return nil
}

func cloneBytes(data []byte) []byte {
	copyOfData := make([]byte, len(data))
	copy(copyOfData, data)
	return copyOfData
}
