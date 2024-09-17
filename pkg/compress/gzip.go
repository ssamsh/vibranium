package compress

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sync/atomic"

	"github.com/gammazero/workerpool"
	log "github.com/sirupsen/logrus"
	"github.com/ssamsh/vibranium/pkg/discovery"
)

func ListFile(outputDir string, prefix string) error {
	return processFile(discovery.File{
		Path:         path.Join(outputDir, "list.txt"),
		StrippedPath: "list.txt",
	}, outputDir, prefix)
}

func Files(files []discovery.File, outputDir string, prefix string) error {
	// this will be removed once we have cache
	if err := os.RemoveAll(path.Join(outputDir, prefix)); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("cannot remove old patch directory: %w", err)
		}
	}

	log.Infof("Starting up %d threads", runtime.NumCPU())

	wp := workerpool.New(runtime.NumCPU())

	var lastErr atomic.Value

	for _, file := range files {
		wp.Submit(func() {
			if err := processFile(file, outputDir, prefix); err != nil {
				err = fmt.Errorf("cannot process file %s: %w", file.Path, err)

				log.Errorf(err.Error())

				lastErr.Store(err)
			}
		})
	}

	wp.StopWait()

	return nil
}

func processFile(file discovery.File, outputDir string, prefix string) error {
	log.Infof("Compressing file: %s", file.StrippedPath)

	outputFilePath := path.Join(outputDir, prefix, file.StrippedPath+".gz")

	if err := os.MkdirAll(filepath.Dir(outputFilePath), os.ModePerm); err != nil {
		return fmt.Errorf("cannot make directory: %w", err)
	}

	inputFile, err := os.Open(file.Path)
	if err != nil {
		return fmt.Errorf("cannot open input file: %w", err)
	}

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("cannot create file: %w", err)
	}

	gzipWriter := gzip.NewWriter(outputFile)

	if _, err = io.Copy(gzipWriter, inputFile); err != nil {
		return fmt.Errorf("cannot write gzip data: %w", err)
	}

	defer func() {
		if err = gzipWriter.Close(); err != nil {
			log.Errorf("failed to close gzip writer: %w", err)
		}

		if err = outputFile.Close(); err != nil {
			log.Errorf("failed to close output file: %w", err)
		}

		if err = inputFile.Close(); err != nil {
			log.Errorf("failed to close input file: %w", err)
		}
	}()

	return nil
}
