package io

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

// ReadConfigs reads all files with extension ext from dirName and returns a slice of T.
func ReadConfigs[T any](dirName, ext string) ([]T, error) {
	var cfgs []T

	// open dir
	dir := os.DirFS(dirName)

	// walk dir
	err := fs.WalkDir(dir, ".", func(path string, d fs.DirEntry, err error) error {
		var cfg T

		if path == "." {
			return nil
		}

		if err != nil {
			log.Fatalln("couldn't read file:", err)
		}
		if filepath.Ext(path) != ext {
			return nil
		}

		f, err := os.Open("./" + dirName + "/" + path)
		if err != nil {
			return err
		}
		defer func() { _ = f.Close() }()

		if err := json.NewDecoder(f).Decode(&cfg); err != nil {
			return fmt.Errorf("wrong tower format: %w", err)
		}

		cfgs = append(cfgs, cfg)

		return nil
	})
	if err != nil {
		return nil, err
	}

	cfgsCopy := make([]T, len(cfgs))
	copy(cfgsCopy, cfgs)
	return cfgsCopy, nil
}
