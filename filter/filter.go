package filter

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"myapp/config"
	"myapp/logger"
	"myapp/common"
)

func Run(cfg *config.Config) error {
	logger.Info.Println("Starting filtering component")
	incoming := cfg.IncomingDir
	accepted := cfg.AcceptedDir
	rejected := cfg.RejectedDir

	files, err := os.ReadDir(incoming)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		srcPath := filepath.Join(incoming, file.Name())
		data, err := ioutil.ReadFile(srcPath)
		if err != nil {
			logger.Error.Println("Error reading file:", err)
			continue
		}
		// Dummy filter: accept if file contains the word "accept".
		if strings.Contains(string(data), "accept") {
			dest := filepath.Join(accepted, file.Name())
			if err := common.MoveFile(srcPath, dest); err != nil {
				logger.Error.Println("Error moving file to accepted:", err)
				continue
			}
			logger.Info.Println("File accepted:", file.Name())
		} else {
			dest := filepath.Join(rejected, file.Name())
			if err := common.MoveFile(srcPath, dest); err != nil {
				logger.Error.Println("Error moving file to rejected:", err)
				continue
			}
			logger.Info.Println("File rejected:", file.Name())
		}
	}
	return nil
}
