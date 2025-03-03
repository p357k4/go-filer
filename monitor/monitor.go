package monitor

import (
	"os"
	"path/filepath"
	"time"

	"myapp/config"
	"myapp/logger"
	"myapp/common"
)

func Run(cfg *config.Config) error {
	logger.Info.Println("Starting monitoring component")
	landing := cfg.LandingDir
	incoming := cfg.IncomingDir
	fileSizes := make(map[string]int64)
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		files, err := os.ReadDir(landing)
		if err != nil {
			logger.Error.Println("Error reading landing directory:", err)
			continue
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			path := filepath.Join(landing, file.Name())
			info, err := os.Stat(path)
			if err != nil {
				logger.Error.Println("Error stating file:", err)
				continue
			}
			size := info.Size()
			// Check if size unchanged.
			if prevSize, ok := fileSizes[path]; ok && prevSize == size {
				dest := filepath.Join(incoming, file.Name())
				if err := common.MoveFile(path, dest); err != nil {
					logger.Error.Println("Error moving file:", err)
					continue
				}
				logger.Info.Println("Moved file to incoming:", file.Name())
				delete(fileSizes, path)
			} else {
				fileSizes[path] = size
			}
		}
	}
	// (Never reached)
	// return nil
}
