package clean

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"myapp/common"
	"myapp/config"
	"myapp/logger"
)

// Dummy uploadCleaningData simulates an upload step during cleaning.
func uploadCleaningData(filePath, bucket, serviceAccountKey string) error {
	time.Sleep(1 * time.Second)
	logger.Info.Println("Uploaded cleaning data:", filePath)
	return nil
}

func Run(cfg *config.Config) error {
	logger.Info.Println("Starting cleaning component")
	uploaded := cfg.UploadedDir
	completed := cfg.CompletedDir
	manifestUploaded := cfg.ManifestUploaded
	manifestFailed := cfg.ManifestFailed

	// For demonstration, we list manifest files in the manifestUploaded folder.
	manifests, err := os.ReadDir(manifestUploaded)
	if err != nil {
		return err
	}

	for _, manifest := range manifests {
		if manifest.IsDir() {
			continue
		}
		manifestPath := filepath.Join(manifestUploaded, manifest.Name())
		data, err := os.ReadFile(manifestPath)
		if err != nil {
			logger.Error.Println("Error reading manifest file:", err)
			continue
		}
		// Each line in the manifest is a file name in the uploaded folder.
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}
			src := filepath.Join(uploaded, line)
			dest := filepath.Join(completed, line)
			if err := common.MoveFile(src, dest); err != nil {
				logger.Error.Printf("Error moving file %s: %v", line, err)
				attempts := 0
				for attempts < cfg.MaxUploadAttempts {
					attempts++
					time.Sleep(2 * time.Second)
					if err := common.MoveFile(src, dest); err == nil {
						break
					}
					if attempts >= cfg.MaxUploadAttempts {
						failedManifest := filepath.Join(manifestFailed, manifest.Name())
						if err := common.MoveFile(manifestPath, failedManifest); err != nil {
							logger.Error.Println("Error moving manifest to failed:", err)
						}
						logger.Error.Printf("Max attempts reached for cleaning file %s", line)
					}
				}
			} else {
				logger.Info.Println("Cleaned file:", line)
			}
		}
		logger.Info.Println("Completed cleaning for manifest:", manifest.Name())
	}
	return nil
}
