package common

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/afero"
)

func LoadData[T any](fs afero.Fs, filename string) (*T, error) {
	if _, err := fs.Stat(filename); os.IsNotExist(err) {
		return nil, nil
	}

	data, err := afero.ReadFile(fs, filename)
	if err != nil {
		return nil, fmt.Errorf("couldn't read data file: %s", err.Error())
	}

	var result T
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("couldn't unmarshal data: %s", err.Error())
	}

	return &result, nil
}

func SaveData[T any](fs afero.Fs, filename string, value *T) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("couldn't marshal data: %s", err.Error())
	}

	err = afero.WriteFile(fs, filename, data, 0644)
	if err != nil {
		return fmt.Errorf("couldn't write data file: %s", err.Error())
	}

	return nil
}
