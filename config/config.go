package config

import (
	"encoding/json"
	"os"
)

// Config holds application settings.
type Config struct {
	LandingDir        string `json:"landing_dir"`
	IncomingDir       string `json:"incoming_dir"`
	AcceptedDir       string `json:"accepted_dir"`
	RejectedDir       string `json:"rejected_dir"`
	UploadedDir       string `json:"uploaded_dir"`
	FailedDir         string `json:"failed_dir"`
	ConfirmedDir      string `json:"confirmed_dir"`
	CompletedDir      string `json:"completed_dir"`
	ManifestIncoming  string `json:"manifest_incoming"`
	ManifestCompleted string `json:"manifest_completed"`
	ManifestUploaded  string `json:"manifest_uploaded"`
	ManifestFailed    string `json:"manifest_failed"`
	GCSBucket         string `json:"gcs_bucket"`
	ServiceAccountKey string `json:"service_account_key"`
	LogLevel          string `json:"log_level"`
	MaxUploadAttempts int    `json:"max_upload_attempts"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
