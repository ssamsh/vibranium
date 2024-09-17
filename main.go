package main

import (
	"flag"
	"time"

	"github.com/ssamsh/vibranium/pkg/compress"
	"github.com/ssamsh/vibranium/pkg/discovery"
	"github.com/ssamsh/vibranium/pkg/list"

	log "github.com/sirupsen/logrus"
)

var (
	listVersion  = flag.Int("listVersion", 21, "-listVersion=21")
	inputDir     = flag.String("inputDir", "", "help message for flag n")
	outputDir    = flag.String("outputDir", "", "help message for flag n")
	using12hTime = flag.Bool("using12hTime", false, "-using12hTime")
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})
	startTime := time.Now()

	flag.Parse()

	if inputDir == nil || *inputDir == "" {
		log.Fatalf("Missing argument 'inputDir'")
		return
	}

	if outputDir == nil || *outputDir == "" {
		log.Fatalf("Missing argument 'outputDir'")
		return
	}

	if using12hTime == nil || *using12hTime {
		log.Fatalf("12h time is currently not supported")
		return
	}

	log.Info("Discovering files")

	discoveredFiles, err := discovery.DiscoverFiles(*inputDir)
	if err != nil {
		log.Fatalf("Cannot discover files: %s", err)
	}

	log.Infof("Discovered %d files", len(discoveredFiles))

	if err := compress.Files(discoveredFiles, *outputDir, "NeuroSpace\\RESCLIENT"); err != nil {
		log.Fatalf("Cannot compress files: %s", err)
	}

	if err := list.Make(list.MakeOpts{
		InputDir:     *inputDir,
		OutputDir:    *outputDir,
		Using12hTime: *using12hTime,
		Version:      *listVersion,
		Files:        discoveredFiles,
	}); err != nil {
		log.Fatalf("Cannot make list: %s", err)
	}

	if err := compress.ListFile(*outputDir, "NeuroSpace\\RESCLIENT"); err != nil {
		log.Fatalf("Cannot compress files: %s", err)
	}

	log.Infof("Job is done, took %s", time.Since(startTime).Round(time.Millisecond))
}
