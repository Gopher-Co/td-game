package io

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
)

func ReadConfigs[T any](dirName, ext string) ([]T, error) {
	var cfg T
	var cfgs []T

	dir := os.DirFS(dirName)

	err := fs.WalkDir(dir, ".", func(path string, d fs.DirEntry, err error) error {
		if path == "." {
			return nil
		}

		if err != nil {
			log.Fatalln("couldn't read tower file:", err)
		}
		if filepath.Ext(path) != ext {
			return nil
		}

		f, err := os.Open("./" + dirName + "/" + path)
		if err != nil {
			return err
		}

		if err := json.NewDecoder(f).Decode(&cfg); err != nil {
			return fmt.Errorf("wrong tower format: %w", err)
		}

		cfgs = append(cfgs, cfg)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return slices.Clip(cfgs), nil
}
