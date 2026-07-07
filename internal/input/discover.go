package input

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/marcinbohm/search-index-preflight/internal/model"
)

var ignoredDirectories = map[string]struct{}{
	".git":         {},
	"node_modules": {},
	"vendor":       {},
	"dist":         {},
	"build":        {},
	".local":       {},
}

func LoadFile(path string, kind model.DocumentKind) (Source, error) {
	info, err := os.Stat(path)
	if err != nil {
		return Source{}, fmt.Errorf("load %q: %w", path, err)
	}
	if info.IsDir() {
		return Source{}, fmt.Errorf("load %q: expected file, got directory", path)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return Source{}, fmt.Errorf("load %q: %w", path, err)
	}

	return Source{
		Path:         filepath.Clean(path),
		RelativePath: relativePath(path),
		Kind:         kind,
		Content:      content,
	}, nil
}

func Discover(path string) ([]Source, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("discover %q: %w", path, err)
	}
	if !info.IsDir() {
		source, err := LoadFile(path, inferKind(path))
		if err != nil {
			return nil, err
		}
		return []Source{source}, nil
	}

	var sources []Source
	err = filepath.WalkDir(path, func(current string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if shouldIgnoreDirectory(entry.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if !entry.Type().IsRegular() {
			return nil
		}
		if !isSupportedInputFile(current) {
			return nil
		}

		source, err := LoadFile(current, inferKind(current))
		if err != nil {
			return err
		}
		sources = append(sources, source)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("discover %q: %w", path, err)
	}

	sort.Slice(sources, func(i, j int) bool {
		return sources[i].RelativePath < sources[j].RelativePath
	})
	return sources, nil
}

func shouldIgnoreDirectory(name string) bool {
	_, ok := ignoredDirectories[name]
	return ok
}

func isSupportedInputFile(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".json", ".jsonl", ".ndjson":
		return true
	default:
		return false
	}
}

func inferKind(path string) model.DocumentKind {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".jsonl", ".ndjson":
		return model.DocumentKindSampleDocs
	default:
		return model.DocumentKindUnknown
	}
}

func relativePath(path string) string {
	rel, err := filepath.Rel(".", path)
	if err != nil {
		return filepath.Clean(path)
	}
	return filepath.Clean(rel)
}
