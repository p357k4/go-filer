package upload

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"myapp/config"
	"myapp/logger"
	"myapp/common"
)

// Dummy uploadToGCS simulates uploading a file.
func uploadToGCS(filePath, bucket, serviceAccountKey string) (string, error) {
	// In a real implementation, initialize a GCP client using the serviceAccountKey.
	time.Sleep(1 * time.Second)
	targetPath := fmt.Sprintf("gs://%s/%s", bucket, filepath.Base(filePath))
	return targetPath, nil
}

func Run(cfg *config.Config) error {
	logger.Info.Println("Starting uploading component")
	accepted := cfg.AcceptedDir
	confirmed := cfg.ConfirmedDir
	uploaded := cfg.UploadedDir
	failed := cfg.FailedDir

	files, err := os.ReadDir(accepted)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		srcPath := filepath.Join(accepted, file.Name())
		var targetPath string
		attempts := 0
		maxAttempts := cfg.MaxUploadAttempts
		for {
			attempts++
			tp, err := uploadToGCS(srcPath, cfg.GCSBucket, cfg.ServiceAccountKey)
			if err == nil {
				targetPath = tp
				break
			}
			logger.Error.Printf("Upload attempt %d failed for %s: %v", attempts, file.Name(), err)
			if attempts >= maxAttempts {
				// Move file to failed folder.
				dest := filepath.Join(failed, file.Name())
				if err := common.MoveFile(srcPath, dest); err != nil {
					logger.Error.Println("Error moving file to failed:", err)
				}
				logger.Error.Println("Max upload attempts reached for", file.Name())
				break
			}
			time.Sleep(2 * time.Second)
		}
		if targetPath != "" {
			// Create confirmation file.
			confirmationPath := filepath.Join(confirmed, file.Name()+".confirm")
			content := fmt.Sprintf("source: %s\ntarget: %s\n", srcPath, targetPath)
			if err := os.WriteFile(confirmationPath, []byte(content), 0644); err != nil {
				logger.Error.Println("Error writing confirmation file:", err)
			}
			// Move file to uploaded folder.
			dest := filepath.Join(uploaded, file.Name())
			if err := common.MoveFile(srcPath, dest); err != nil {
				logger.Error.Println("Error moving file to uploaded:", err)
			}
			logger.Info.Println("Uploaded file:", file.Name())
		}
	}
	return nil
}
