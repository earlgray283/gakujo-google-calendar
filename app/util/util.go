package util

import (
	"errors"
	"os"
	"path/filepath"
)

func DefaultConfigDir() string {
	configPath, _ := os.UserConfigDir()
	configDirPath := filepath.Join(configPath, "gakujo-google-calendar")
	if err := os.MkdirAll(configDirPath, 0755); err != nil {
		panic(err)
	}
	return configDirPath
}

func ExistsFile(path string) bool {
	_, err := os.Stat(path)
	return errors.Is(err, os.ErrNotExist)
}

func DefaultCacheDir() string {
	uCacheDirPath, _ := os.UserCacheDir()
	cacheDirPath := filepath.Join(uCacheDirPath, "gakujo-google-calendar")
	if err := os.MkdirAll(cacheDirPath, 0755); err != nil {
		panic(err)
	}
	return cacheDirPath
}
