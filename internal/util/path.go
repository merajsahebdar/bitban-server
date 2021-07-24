package util

import (
	"fmt"
	"os"
)

// GetAssetPath Returns path to the requested asset file.
func GetAssetPath(asset string) (string, error) {
	dirs := []string{
		"./configs",
		"/etc/giteam",
	}

	for _, dir := range dirs {
		path := dir + asset
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("cannot find asset: %s", asset)
}
