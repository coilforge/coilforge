package appsettings

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	configDirName  = "coilforge"
	configFileName = "appsettings.json"
)

// LoadLocal loads app settings from the platform user config directory.
// If no stored settings file exists, Defaults() are returned.
func LoadLocal() (Values, error) {
	path, err := localSettingsPath()
	if err != nil {
		return Defaults(), err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Defaults(), nil
		}
		return Defaults(), err
	}
	var values Values
	if err := json.Unmarshal(data, &values); err != nil {
		return Defaults(), err
	}
	return values, nil
}

// SaveLocalCurrent stores Current in the platform user config directory.
func SaveLocalCurrent() error {
	path, err := localSettingsPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(Current, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func localSettingsPath() (string, error) {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cfgDir, configDirName, configFileName), nil
}
