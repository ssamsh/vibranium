package discovery

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

type File struct {
	Path         string
	StrippedPath string
	Info         fs.FileInfo
}

func DiscoverFiles(basePath string) ([]File, error) {
	out := make([]File, 0, 1024)

	if err := filepath.WalkDir(basePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk dir %s: %w", path, err)
		}

		// Quick fix for testing

		if strings.Contains(path, "NeuroSpace") {
			return nil
		}

		// We only care about files

		if d.IsDir() {
			return nil
		}

		fileInfo, err := d.Info()
		if err != nil {
			return fmt.Errorf("cannot get file %s info: %w", path, err)
		}

		file := File{
			Path:         path,
			StrippedPath: strings.TrimPrefix(path, basePath+"\\"),
			Info:         fileInfo,
		}

		out = append(out, file)

		return nil
	}); err != nil {
		return nil, fmt.Errorf("error while walking directory %s: %w", basePath, err)
	}

	return out, nil
}
