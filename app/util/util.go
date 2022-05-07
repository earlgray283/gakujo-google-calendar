package util

import (
	"errors"
	"os"
	"path/filepath"
	"time"
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
	cacheDirPath := filepath.Join(uCacheDirPath, "gakujo-google-calendar", time.Now().Format("2006_01_02_15_04_05"))
	if err := os.MkdirAll(cacheDirPath, 0755); err != nil {
		panic(err)
	}
	return cacheDirPath
}

func DoWithRetry(f func() error, retryCount int, interval time.Duration) error {
	var err error = nil
	for i := 0; i < retryCount; i++ {
		err = f()
		if err == nil {
			break
		}
		time.Sleep(interval)
	}
	return err
}
