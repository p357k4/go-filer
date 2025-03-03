package track

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"myapp/config"
	"myapp/logger"
	"myapp/common"
)

// Dummy uploadManifestToGCS simulates manifest upload.
func uploadManifestToGCS(manifestPath, bucket, serviceAccountKey string) error {
	time.Sleep(1 * time.Second)
	logger.Info.Println("Uploaded manifest to GCS:", manifestPath)
	return nil
}

func Run(cfg *config.Config) error {
	logger.Info.Println("Starting tracking component")
	confirmed := cfg.ConfirmedDir
	manifestIncoming := cfg.ManifestIncoming
	manifestCompleted := cfg.ManifestCompleted
	manifestUploaded := cfg.ManifestUploaded

	// Create a new manifest file.
	manifestFile := filepath.Join(manifestIncoming, fmt.Sprintf("manifest_%d.txt", time.Now().Unix()))
	mf, err := os.OpenFile(manifestFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer mf.Close()

	files, err := os.ReadDir(confirmed)
	if err != nil {
		return err
	}
	count := 0
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filePath := filepath.Join(confirmed, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			logger.Error.Println("Error reading confirmation file:", err)
			continue
		}
		// Extract the target file path.
		lines := strings.Split(string(data), "\n")
		var target string
		for _, line := range lines {
			if strings.HasPrefix(line, "target:") {
				target = strings.TrimSpace(strings.TrimPrefix(line, "target:"))
				break
			}
		}
		if target != "" {
			entry := fmt.Sprintf("%s\n", target)
			if _, err := mf.WriteString(entry); err != nil {
				logger.Error.Println("Error writing to manifest file:", err)
			}
			count++
			// Remove the confirmation file.
			os.Remove(filePath)
		}
	}
	mf.Close()

	// Check if the manifest should be completed.
	fi, err := os.Stat(manifestFile)
	if err != nil {
		return err
	}
	if count >= 1000 || time.Since(fi.ModTime()) > time.Hour {
		completedManifest := filepath.Join(manifestCompleted, filepath.Base(manifestFile))
		if err := common.MoveFile(manifestFile, completedManifest); err != nil {
			return err
		}
		if err := uploadManifestToGCS(completedManifest, cfg.GCSBucket, cfg.ServiceAccountKey); err != nil {
			return err
		}
		finalManifest := filepath.Join(manifestUploaded, filepath.Base(manifestFile))
		if err := common.MoveFile(completedManifest, finalManifest); err != nil {
			return err
		}
		logger.Info.Println("Completed manifest processing:", filepath.Base(manifestFile))
	}
	return nil
}
