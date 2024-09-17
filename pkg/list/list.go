package list

import (
	"fmt"
	"os"
	"path"

	"github.com/ssamsh/vibranium/pkg/discovery"

	log "github.com/sirupsen/logrus"
)

type MakeOpts struct {
	InputDir     string
	OutputDir    string
	Using12hTime bool
	Version      int
	Files        []discovery.File
}

func Make(opts MakeOpts) error {
	log.Info("Building list.txt")

	listFile, err := os.OpenFile(path.Join(opts.OutputDir, "list.txt"), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("cannot open list.txt: %s", err)
	}

	defer func() {
		if err := listFile.Close(); err != nil {
			log.Warningf("cannot close list.txt: %s", err)
		}
	}()

	if err = listFile.Truncate(0); err != nil {
		return fmt.Errorf("cannot erase list.txt content (truncate): %s", err)
	}

	listFile.WriteString(fmt.Sprintf("ver:%d\r\n", opts.Version))
	listFile.WriteString(" Z:\\\\RESCLIENT\r\n")

	for _, discoveredFile := range opts.Files {
		ampm := "a"

		hour := discoveredFile.Info.ModTime().Hour()

		if !opts.Using12hTime {
			if hour == 0 {
				ampm = "p"
				hour = 12
			} else if hour == 12 {
				ampm = "p"
			} else if hour > 12 {
				ampm = "p"
				hour -= 12
			}
		}

		fileEntry := fmt.Sprintf("%04d-%02d-%02d  %02d:%02d%s %19d %s\r\n",
			discoveredFile.Info.ModTime().Year(),
			discoveredFile.Info.ModTime().Month(),
			discoveredFile.Info.ModTime().Day(),
			hour,
			discoveredFile.Info.ModTime().Minute(),
			ampm,
			discoveredFile.Info.Size(),
			discoveredFile.StrippedPath,
		)

		if _, err = listFile.WriteString(fileEntry); err != nil {
			return fmt.Errorf("cannot write list.txt file entry: %w", err)
		}
	}

	log.Info("list.txt built successfully")
	return nil
}
